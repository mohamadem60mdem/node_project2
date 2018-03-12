// Copyright (c) 2017, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package interp

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"mvdan.cc/sh/syntax"
)

func (r *Runner) expandFormat(format string, args []string) (int, string, error) {
	buf := r.strBuilder()
	esc := false
	var fmts []rune
	initialArgs := len(args)

	for _, c := range format {
		switch {
		case esc:
			esc = false
			switch c {
			case 'n':
				buf.WriteRune('\n')
			case 'r':
				buf.WriteRune('\r')
			case 't':
				buf.WriteRune('\t')
			case '\\':
				buf.WriteRune('\\')
			default:
				buf.WriteRune('\\')
				buf.WriteRune(c)
			}

		case len(fmts) > 0:
			switch c {
			case '%':
				buf.WriteByte('%')
				fmts = nil
			case 'c':
				var b byte
				if len(args) > 0 {
					arg := ""
					arg, args = args[0], args[1:]
					if len(arg) > 0 {
						b = arg[0]
					}
				}
				buf.WriteByte(b)
				fmts = nil
			case '+', '-', ' ':
				if len(fmts) > 1 {
					return 0, "", fmt.Errorf("invalid format char: %c", c)
				}
				fmts = append(fmts, c)
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				fmts = append(fmts, c)
			case 's', 'd', 'i', 'u', 'o', 'x':
				arg := ""
				if len(args) > 0 {
					arg, args = args[0], args[1:]
				}
				var farg interface{} = arg
				if c != 's' {
					n, _ := strconv.ParseInt(arg, 0, 0)
					if c == 'i' || c == 'd' {
						farg = int(n)
					} else {
						farg = uint(n)
					}
					if c == 'i' || c == 'u' {
						c = 'd'
					}
				}
				fmts = append(fmts, c)
				fmt.Fprintf(buf, string(fmts), farg)
				fmts = nil
			default:
				return 0, "", fmt.Errorf("invalid format char: %c", c)
			}
		case c == '\\':
			esc = true
		case args != nil && c == '%':
			// if args == nil, we are not doing format
			// arguments
			fmts = []rune{c}
		default:
			buf.WriteRune(c)
		}
	}
	if len(fmts) > 0 {
		return 0, "", fmt.Errorf("missing format char")
	}
	return initialArgs - len(args), buf.String(), nil
}

func (r *Runner) fieldJoin(parts []fieldPart) string {
	switch len(parts) {
	case 0:
		return ""
	case 1: // short-cut without a string copy
		return parts[0].val
	}
	buf := r.strBuilder()
	for _, part := range parts {
		buf.WriteString(part.val)
	}
	return buf.String()
}

func (r *Runner) escapedGlobStr(val string) string {
	if !anyPatternChar(val) { // short-cut without a string copy
		return val
	}
	buf := r.strBuilder()
	for _, r := range val {
		if patternRune(r) {
			buf.WriteByte('\\')
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

func (r *Runner) escapedGlobField(parts []fieldPart) (escaped string, glob bool) {
	buf := r.strBuilder()
	for _, part := range parts {
		for _, r := range part.val {
			if patternRune(r) {
				if part.quote > quoteNone {
					buf.WriteByte('\\')
				} else {
					glob = true
				}
			}
			buf.WriteRune(r)
		}
	}
	if glob { // only copy the string if it will be used
		escaped = buf.String()
	}
	return escaped, glob
}

// TODO: consider making brace a special syntax Node

type brace struct {
	elems []*braceWord
}

// braceWord is like syntax.Word, but with braceWordPart.
type braceWord struct {
	parts []braceWordPart
}

// braceWordPart contains either syntax.WordPart or brace.
type braceWordPart interface{}

var (
	litLeftBrace  = &syntax.Lit{Value: "{"}
	litComma      = &syntax.Lit{Value: ","}
	litRightBrace = &syntax.Lit{Value: "}"}
)

func (r *Runner) splitBraces(word *syntax.Word) (*braceWord, bool) {
	any := false
	top := &r.braceAlloc
	*top = braceWord{parts: r.bracePartsAlloc[:0]}
	acc := top
	var cur *brace
	open := []*brace{}

	pop := func() *brace {
		old := cur
		open = open[:len(open)-1]
		if len(open) == 0 {
			cur = nil
			acc = top
		} else {
			cur = open[len(open)-1]
			acc = cur.elems[len(cur.elems)-1]
		}
		return old
	}

	for _, wp := range word.Parts {
		lit, ok := wp.(*syntax.Lit)
		if !ok {
			acc.parts = append(acc.parts, wp)
			continue
		}
		last := 0
		for j, r := range lit.Value {
			addlit := func() {
				if last == j {
					return // empty lit
				}
				l2 := *lit
				l2.Value = l2.Value[last:j]
				acc.parts = append(acc.parts, &l2)
			}
			switch r {
			case '{':
				addlit()
				acc = &braceWord{}
				cur = &brace{elems: []*braceWord{acc}}
				open = append(open, cur)
			case ',':
				if cur == nil {
					continue
				}
				addlit()
				acc = &braceWord{}
				cur.elems = append(cur.elems, acc)
			case '}':
				if cur == nil {
					continue
				}
				any = true
				addlit()
				ended := pop()
				if len(ended.elems) > 1 {
					acc.parts = append(acc.parts, ended)
					break
				}
				// return {x} to a non-brace
				acc.parts = append(acc.parts, litLeftBrace)
				acc.parts = append(acc.parts, ended.elems[0].parts...)
				acc.parts = append(acc.parts, litRightBrace)
			default:
				continue
			}
			last = j + 1
		}
		if last == 0 {
			acc.parts = append(acc.parts, lit)
		} else {
			left := *lit
			left.Value = left.Value[last:]
			acc.parts = append(acc.parts, &left)
		}
	}
	// open braces that were never closed fall back to non-braces
	for acc != top {
		ended := pop()
		acc.parts = append(acc.parts, litLeftBrace)
		for i, elem := range ended.elems {
			if i > 0 {
				acc.parts = append(acc.parts, litComma)
			}
			acc.parts = append(acc.parts, elem.parts...)
		}
	}
	return top, any
}

func expandRec(bw *braceWord) []*syntax.Word {
	var all []*syntax.Word
	var left []syntax.WordPart
	for i, wp := range bw.parts {
		br, ok := wp.(*brace)
		if !ok {
			left = append(left, wp.(syntax.WordPart))
			continue
		}
		for _, elem := range br.elems {
			next := *bw
			next.parts = next.parts[i+1:]
			next.parts = append(elem.parts, next.parts...)
			exp := expandRec(&next)
			for _, w := range exp {
				w.Parts = append(left, w.Parts...)
			}
			all = append(all, exp...)
		}
		return all
	}
	return []*syntax.Word{{Parts: left}}
}

func (r *Runner) expandBraces(word *syntax.Word) []*syntax.Word {
	// TODO: be a no-op when not in bash mode
	topBrace, any := r.splitBraces(word)
	if !any { // short-cut without further work
		r.oneWord[0] = word
		return r.oneWord[:]
	}
	return expandRec(topBrace)
}

func (r *Runner) Fields(words ...*syntax.Word) []string {
	fields := make([]string, 0, len(words))
	baseDir := r.escapedGlobStr(r.Dir)
	for _, word := range words {
		for _, expWord := range r.expandBraces(word) {
			for _, field := range r.wordFields(expWord.Parts) {
				path, glob := r.escapedGlobField(field)
				var matches []string
				abs := filepath.IsAbs(path)
				if glob && !r.shellOpts[optNoGlob] {
					if !abs {
						path = filepath.Join(baseDir, path)
					}
					matches, _ = filepath.Glob(path)
				}
				if len(matches) == 0 {
					fields = append(fields, r.fieldJoin(field))
					continue
				}
				for _, match := range matches {
					if !abs {
						match, _ = filepath.Rel(r.Dir, match)
					}
					fields = append(fields, match)
				}
			}
		}
	}
	return fields
}

func (r *Runner) loneWord(word *syntax.Word) string {
	if word == nil {
		return ""
	}
	field := r.wordField(word.Parts, quoteDouble)
	return r.fieldJoin(field)
}

func (r *Runner) lonePattern(word *syntax.Word) string {
	field := r.wordField(word.Parts, quoteSingle)
	buf := r.strBuilder()
	for _, part := range field {
		for _, r := range part.val {
			if part.quote > quoteNone && patternRune(r) {
				buf.WriteByte('\\')
			}
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (r *Runner) expandAssigns(as *syntax.Assign) []*syntax.Assign {
	// Convert "declare $x" into "declare value".
	// Don't use syntax.Parser here, as we only want the basic
	// splitting by '='.
	if as.Name != nil {
		return []*syntax.Assign{as} // nothing to do
	}
	var asgns []*syntax.Assign
	for _, field := range r.Fields(as.Value) {
		as := &syntax.Assign{}
		parts := strings.SplitN(field, "=", 2)
		as.Name = &syntax.Lit{Value: parts[0]}
		if len(parts) == 1 {
			as.Naked = true
		} else {
			as.Value = &syntax.Word{Parts: []syntax.WordPart{
				&syntax.Lit{Value: parts[1]},
			}}
		}
		asgns = append(asgns, as)
	}
	return asgns
}

type fieldPart struct {
	val   string
	quote quoteLevel
}

type quoteLevel uint

const (
	quoteNone quoteLevel = iota
	quoteDouble
	quoteSingle
)

func (r *Runner) wordField(wps []syntax.WordPart, ql quoteLevel) []fieldPart {
	var field []fieldPart
	for i, wp := range wps {
		switch x := wp.(type) {
		case *syntax.Lit:
			s := x.Value
			if i == 0 {
				s = r.expandUser(s)
			}
			if ql == quoteDouble && strings.Contains(s, "\\") {
				buf := r.strBuilder()
				for i := 0; i < len(s); i++ {
					b := s[i]
					if b == '\\' && i+1 < len(s) {
						switch s[i+1] {
						case '\n': // remove \\\n
							i++
							continue
						case '\\', '$', '`': // special chars
							continue
						}
					}
					buf.WriteByte(b)
				}
				s = buf.String()
			}
			field = append(field, fieldPart{val: s})
		case *syntax.SglQuoted:
			fp := fieldPart{quote: quoteSingle, val: x.Value}
			if x.Dollar {
				_, fp.val, _ = r.expandFormat(fp.val, nil)
			}
			field = append(field, fp)
		case *syntax.DblQuoted:
			field = append(field, r.wordField(x.Parts, quoteDouble)...)
		case *syntax.ParamExp:
			field = append(field, fieldPart{val: r.paramExp(x)})
		case *syntax.CmdSubst:
			field = append(field, fieldPart{val: r.cmdSubst(x)})
		case *syntax.ArithmExp:
			field = append(field, fieldPart{
				val: strconv.Itoa(r.arithm(x.X)),
			})
		default:
			panic(fmt.Sprintf("unhandled word part: %T", x))
		}
	}
	return field
}

func (r *Runner) cmdSubst(cs *syntax.CmdSubst) string {
	r2 := r.sub()
	buf := r.strBuilder()
	r2.Stdout = buf
	r2.stmts(cs.StmtList)
	r.setErr(r2.err)
	return strings.TrimRight(buf.String(), "\n")
}

func (r *Runner) wordFields(wps []syntax.WordPart) [][]fieldPart {
	fields := r.fieldsAlloc[:0]
	curField := r.fieldAlloc[:0]
	allowEmpty := false
	flush := func() {
		if len(curField) == 0 {
			return
		}
		fields = append(fields, curField)
		curField = nil
	}
	splitAdd := func(val string) {
		for i, field := range strings.FieldsFunc(val, r.ifsRune) {
			if i > 0 {
				flush()
			}
			curField = append(curField, fieldPart{val: field})
		}
	}
	for i, wp := range wps {
		switch x := wp.(type) {
		case *syntax.Lit:
			s := x.Value
			if i == 0 {
				s = r.expandUser(s)
			}
			if strings.Contains(s, "\\") {
				buf := r.strBuilder()
				for i := 0; i < len(s); i++ {
					b := s[i]
					if b == '\\' {
						i++
						b = s[i]
					}
					buf.WriteByte(b)
				}
				s = buf.String()
			}
			curField = append(curField, fieldPart{val: s})
		case *syntax.SglQuoted:
			allowEmpty = true
			fp := fieldPart{quote: quoteSingle, val: x.Value}
			if x.Dollar {
				_, fp.val, _ = r.expandFormat(fp.val, nil)
			}
			curField = append(curField, fp)
		case *syntax.DblQuoted:
			allowEmpty = true
			if len(x.Parts) == 1 {
				pe, _ := x.Parts[0].(*syntax.ParamExp)
				if elems := r.quotedElems(pe); elems != nil {
					for i, elem := range elems {
						if i > 0 {
							flush()
						}
						curField = append(curField, fieldPart{
							quote: quoteDouble,
							val:   elem,
						})
					}
					continue
				}
			}
			for _, part := range r.wordField(x.Parts, quoteDouble) {
				part.quote = quoteDouble
				curField = append(curField, part)
			}
		case *syntax.ParamExp:
			splitAdd(r.paramExp(x))
		case *syntax.CmdSubst:
			splitAdd(r.cmdSubst(x))
		case *syntax.ArithmExp:
			curField = append(curField, fieldPart{
				val: strconv.Itoa(r.arithm(x.X)),
			})
		default:
			panic(fmt.Sprintf("unhandled word part: %T", x))
		}
	}
	flush()
	if allowEmpty && len(fields) == 0 {
		fields = append(fields, curField)
	}
	return fields
}

func (r *Runner) expandUser(field string) string {
	if len(field) == 0 || field[0] != '~' {
		return field
	}
	name := field[1:]
	rest := ""
	if i := strings.Index(name, "/"); i >= 0 {
		rest = name[i:]
		name = name[:i]
	}
	if name == "" {
		return r.getVar("HOME") + rest
	}
	u, err := user.Lookup(name)
	if err != nil {
		return field
	}
	return u.HomeDir + rest
}
