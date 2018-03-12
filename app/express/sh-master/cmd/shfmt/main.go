// Copyright (c) 2016, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main // import "mvdan.cc/sh/cmd/shfmt"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"mvdan.cc/sh/fileutil"
	"mvdan.cc/sh/syntax"
)

var (
	showVersion = flag.Bool("version", false, "")

	list   = flag.Bool("l", false, "")
	write  = flag.Bool("w", false, "")
	simple = flag.Bool("s", false, "")
	find   = flag.Bool("f", false, "")

	langStr = flag.String("ln", "", "")
	posix   = flag.Bool("p", false, "")

	indent      = flag.Uint("i", 0, "")
	binNext     = flag.Bool("bn", false, "")
	caseIndent  = flag.Bool("ci", false, "")
	keepPadding = flag.Bool("kp", false, "")

	toJSON = flag.Bool("exp.tojson", false, "")

	parser            *syntax.Parser
	printer           *syntax.Printer
	readBuf, writeBuf bytes.Buffer

	copyBuf = make([]byte, 32*1024)

	out io.Writer = os.Stdout

	version = "v2.1.0"
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `usage: shfmt [flags] [path ...]

If no arguments are given, standard input will be used. If a given path
is a directory, it will be recursively searched for shell files - both
by filename extension and by shebang.

  -version  show version and exit

  -l        list files whose formatting differs from shfmt's
  -w        write result to file instead of stdout
  -s        simplify the code
  -f        recursively find all shell files and print the paths

Parser options:

  -ln str   language variant to parse (bash/posix/mksh, default "bash")
  -p        shorthand for -ln=posix

Printer options:

  -i uint   indent: 0 for tabs (default), >0 for number of spaces
  -bn       binary ops like && and | may start a line
  -ci       switch cases will be indented
  -kp       keep column alignment paddings

  -exp.tojson  print AST to stdout as a typed JSON
`)
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}
	if *posix && *langStr != "" {
		fmt.Fprintf(os.Stderr, "-p and -ln=lang cannot coexist\n")
		os.Exit(1)
	}
	lang := syntax.LangBash
	switch *langStr {
	case "bash", "":
	case "posix":
		lang = syntax.LangPOSIX
	case "mksh":
		lang = syntax.LangMirBSDKorn
	default:
		fmt.Fprintf(os.Stderr, "unknown shell language: %s\n", *langStr)
		os.Exit(1)
	}
	if *posix {
		lang = syntax.LangPOSIX
	}
	parser = syntax.NewParser(syntax.KeepComments, syntax.Variant(lang))
	printer = syntax.NewPrinter(func(p *syntax.Printer) {
		syntax.Indent(*indent)(p)
		if *binNext {
			syntax.BinaryNextLine(p)
		}
		if *caseIndent {
			syntax.SwitchCaseIndent(p)
		}
		if *keepPadding {
			syntax.KeepPadding(p)
		}
	})
	if flag.NArg() == 0 {
		if err := formatStdin(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if *toJSON {
		fmt.Fprintln(os.Stderr, "-tojson can only be used with stdin/out")
		os.Exit(1)
	}
	anyErr := false
	for _, path := range flag.Args() {
		walk(path, func(err error) {
			anyErr = true
			fmt.Fprintln(os.Stderr, err)
		})
	}
	if anyErr {
		os.Exit(1)
	}
}

func formatStdin() error {
	if *write || *list {
		return fmt.Errorf("-w and -l can only be used on files")
	}
	prog, err := parser.Parse(os.Stdin, "")
	if err != nil {
		return err
	}
	if *simple {
		syntax.Simplify(prog)
	}
	if *toJSON {
		return writeJSON(out, prog, true)
	}
	return printer.Print(out, prog)
}

var vcsDir = regexp.MustCompile(`^\.(git|svn|hg)$`)

func walk(path string, onError func(error)) {
	info, err := os.Stat(path)
	if err != nil {
		onError(err)
		return
	}
	if !info.IsDir() {
		if err := formatPath(path, false); err != nil {
			onError(err)
		}
		return
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && vcsDir.MatchString(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			onError(err)
			return nil
		}
		conf := fileutil.CouldBeScript(info)
		if conf == fileutil.ConfNotScript {
			return nil
		}
		err = formatPath(path, conf == fileutil.ConfIfShebang)
		if err != nil && !os.IsNotExist(err) {
			onError(err)
		}
		return nil
	})
}

func formatPath(path string, checkShebang bool) error {
	openMode := os.O_RDONLY
	if *write {
		openMode = os.O_RDWR
	}
	f, err := os.OpenFile(path, openMode, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	readBuf.Reset()
	if checkShebang {
		n, err := f.Read(copyBuf[:32])
		if err != nil {
			return err
		}
		if !fileutil.HasShebang(copyBuf[:n]) {
			return nil
		}
		readBuf.Write(copyBuf[:n])
	}
	if *find {
		fmt.Fprintln(out, path)
		return nil
	}
	if _, err := io.CopyBuffer(&readBuf, f, copyBuf); err != nil {
		return err
	}
	src := readBuf.Bytes()
	prog, err := parser.Parse(&readBuf, path)
	if err != nil {
		return err
	}
	if *simple {
		syntax.Simplify(prog)
	}
	writeBuf.Reset()
	printer.Print(&writeBuf, prog)
	res := writeBuf.Bytes()
	if !bytes.Equal(src, res) {
		if *list {
			if _, err := fmt.Fprintln(out, path); err != nil {
				return err
			}
		}
		if *write {
			if err := f.Truncate(0); err != nil {
				return err
			}
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				return err
			}
			if _, err := f.Write(res); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
	if !*list && !*write {
		if _, err := out.Write(res); err != nil {
			return err
		}
	}
	return nil
}
