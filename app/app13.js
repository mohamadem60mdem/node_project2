//   Update Files with fs.appendFile()
/*
یعنی به ادامه فایل اضافه میشود  قسمت قبل پاک نمیشود

Update Files
The File System module has methods for updating files:

fs.appendFile()
fs.writeFile()

*/
var fs = require('fs');

fs.appendFile('mynewfile1.txt', '\n This is my text.', function (err) {
  if (err) throw err;
  console.log('Updated!');
});