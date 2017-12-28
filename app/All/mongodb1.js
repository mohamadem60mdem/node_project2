//fix version 3.x 
//Node.js MongoDB Create Database and createCollection
 
var MongoClient = require('mongodb').MongoClient;
var url = "mongodb://localhost:27017/mydb";




MongoClient.connect(url, function(err, database) {
  if (err) throw err;
  
  const myAwesomeDB = database.db('mydbXXX1');
  //  myAwesomeDB.collection('theCollectionIwantToAccess'){ 

    myAwesomeDB.createCollection("CollectionXXX1", function(err, res) {
    if (err) throw err;
    console.log("Collection and DB created!");
    database.close();
  });
});