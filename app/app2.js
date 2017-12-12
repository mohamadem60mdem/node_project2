var MongoClient = require('mongodb').MongoClient;
var url = "mongodb://localhost:27017/mydb";

//********************************** myobj
    var myobj = {
        name: "Company Inc--004",
        address: "Highway 38"
        };
    var myobj1 = {
    name: "Company Inc--004",
    address: "Highway 38"
    };


        var myobj2 = [
            { name: 'John', address: 'Highway 71'},
            { name: 'Peter', address: 'Lowstreet 4'},
            { name: 'Amy', address: 'Apple st 652'},
            { name: 'Hannah', address: 'Mountain 21'},
            { name: 'Michael', address: 'Valley 345'},
            { name: 'Sandy', address: 'Ocean blvd 2'},
            { name: 'Betty', address: 'Green Grass 1'},
            { name: 'Richard', address: 'Sky st 331'},
            { name: 'Susan', address: 'One way 98'},
            { name: 'Vicky', address: 'Yellow Garden 2'},
            { name: 'Ben', address: 'Park Lane 38'},
            { name: 'William', address: 'Central st 954'},
            { name: 'Chuck', address: 'Main Road 989'},
            { name: 'Viola', address: 'Sideway 1633'}
        ];


      var myobj3= [
        { _id: 154, name: 'Chocolate Heaven'},
        { _id: 155, name: 'Tasty Lemon'},
        { _id: 156, name: 'Vanilla Dream'}
      ];


      var myobj4= {  name: 'Chocolate ', address: 'uyguyg'} ;
      //{ _id: 2, name: '44 Heaven2', address: 'Sideway 1633'} ;Chocolate 44 Valley 345




//********************************** myobj







MongoClient.connect(url, function(err, database) {
    if (err) throw err;


    
    //*** create myAwesomeDB database
    const myAwesomeDB = database.db('mydb2');
  

    //***  createCollection
        myAwesomeDB.createCollection("Collection2", function(err, res) {
        if (err) throw err;
        console.log("Collection and DB created!");
        });
    //*** create createCollection


    //*** collection insertOne insert insert insert
        myAwesomeDB.collection("Collection2").insertOne(myobj4, function(err, res) {
            if (err) throw err;
            console.log("1 document inserted");
        });
    //*** create insertOne




    //***insertMany */
    /*
        myAwesomeDB.collection("Collection2").insertMany(myobj3, function(err, res) {
            if (err) throw err;
            console.log("Number of documents inserted: " + res.insertedCount);
        });
    */

    //***insertMany */

    /*
    //***findOne 
    myAwesomeDB.collection("Collection2").findOne({}, function(err, result) {
        if (err) throw err;
        console.log(result.name);
        
      });
    //***findOne 
    */


     //******find  مدل اول 

     /*
     myAwesomeDB.collection("Collection2").find({}).toArray(function(err, result) {
        if (err) throw err;
        console.log(result);
      });
      */   

     //******find  مدل اول 



    //******find مدل دوم با پارامتر
     
     // { _id: true, name: false, address: false } Chocolate 44
     myAwesomeDB.collection("Collection2").find({name:'Chocolate 44' }).toArray(function(err, result) {
        if (err) throw err;
        console.log(result);
      });
      
        
     //******find مدل دوم با پارامتر


     //******updateOne

/*
     var myquery = { address: "Valley 345" };
     var newvalues = { name: "Mickey", address: "Canyon 123" };
     myAwesomeDB.collection("Collection2").updateOne(myquery, newvalues, function(err, res) {
       if (err) throw err;
       console.log("1 document updated");
    });

*/

    //******updateOne


    //*** database.close(); */
    database.close();

});
 