//formidable uplod file
/*
var http = require('http');
var formidable = require('formidable');
var fs = require('fs');

http.createServer(function (req, res) {
  if (req.url == '/fileupload') {
    var form = new formidable.IncomingForm();
    form.parse(req, function (err, fields, files) {
      var oldpath = files.filetoupload.path;
      var newpath = 'H:/git/emami3/node_project/app/fileupload/' + files.filetoupload.name;
      fs.rename(oldpath, newpath, function (err) {
        if (err) throw err;
        res.write('File uploaded and moved!');
        res.end();
      });
 });
  } else {
    res.writeHead(200, {'Content-Type': 'text/html'});
    res.write('<form action="fileupload" method="post" enctype="multipart/form-data">');
    res.write('<input type="file" name="filetoupload"><br>');
    res.write('<input type="submit">');
    res.write('</form>');
    return res.end();
  }
}).listen(8080);

*/

var formidable = require('formidable'),
http = require('http'),
util = require('util');
var fs = require('fs');



http.createServer(function(req, res) {
if (req.url == '/upload' && req.method.toLowerCase() == 'post') {
// parse a file upload 
var form = new formidable.IncomingForm();
form.encoding = 'utf-8';
form.uploadDir = "H:/git/emami3/node_project/app/fileupload/";
form.type  = true;
//file.name = 'aaaa';

form.parse(req, function(err, fields, files) {
  res.writeHead(200, {'content-type': 'text/plain'});
    //****

     
    



       // var oldpath = files.filetoupload.path;
      //  var newpath = 'H:/git/emami3/node_project/app/fileupload/' + files.filetoupload.name;
     



    //**** 
  res.write('received upload:\n\n');
  res.end(util.inspect({fields: fields, files: files}));
});

return;
}

// show a file upload form   ‎74-D4-35-1D-B6-6C
// ‎74:D4:35:1D:B6:6C
// 74:D4:35:1D:B6:6C
res.writeHead(200, {'content-type': 'text/html'});
res.end(
'<form action="/upload" enctype="multipart/form-data" method="post">'+
'<input type="text" name="title"><br>'+
'<input type="file" name="upload32" multiple="multiple"><br>'+
'<input type="submit" value="Upload">'+
'</form>'
);
}).listen(8080);