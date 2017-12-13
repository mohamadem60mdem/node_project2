//   Update Files with fs.writeFile()
/*
  قسمت قبل پاک میشود

Update Files
 for Replaced files:

fs.appendFile()
fs.writeFile()

*/

var fs = require('fs');

fs.writeFile('mynewfile1.txt', '\nThis is my text1', function (err) {
  if (err) throw err;
  console.log('Replaced!');
});


