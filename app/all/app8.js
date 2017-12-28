//Split the Query String
// راهکار دریافت درخواست http require
/*
http://localhost:8080/?year=2017&month=July

q.year
q.month

 var txt = q.year + " " + q.month;



خاصیت action : این خاصیت از نوع آدرس ( URL ) بوده و تعیین کننده آدرس فرم یا صفحه ای است که قرار است اطلاعات فرم به آن ارسال شود . اطلاعات ارسال شده در صفحه مقصد مورد استفاده و پردازش قرار می گیرد .
برای مثال اگر مقدار آن را برابر با submit.php قرار دهیم ، اطلاعات فرم پس از submit شدن به صفحه تعیین شده ارسال شده و مرورگر نیز به همان صفحه هدایت می شود . 
تعیین مقدار این خاصیت اجباری است و اگر مقدار آن را خالی رها کنیم ، برنامه دچار نقص می شود . در این حالت فرم به یک آدرس پیش فرض مثل form.html می رود .
خاصیت method : این خاصیت روش ارسال اطلاعات فرم  به صفحه را مقصد تعیین کرده و می تواند یکی از دو مقدار GET یا POST را داشته باشد . 
این دو حالت با هم تفاوت عملکرد نداشته و فقط در نحوه ارسال اطلاعات از روش های متفاوتی استفاده می کنند . در جدول زیر به تشزیح نحوه استفاده از هر 2 متد پرداخته ایم :
متد GET : در این روش اطلاعات بصورت ساده و کد نشده منتقل میشوند. این روش دقیقا مشابه اینه که اطلاعات رو بصورت Query به URL اضافه کرده باشیم و وقتی فرم رو Submit می کنیم، این اطلاعات خودش به URL اضافه میشه و قابل دیدن میشه. باید توجه داشت که اطلاعات حساسی مثل Password نباید در معرض دید قرار بگیره پس نباید برای فرمی که اطلاعات مهمی داره از GET استفاده کنیم تا اطلاعات فرم توی Address Bar قابل رویت نشه. به اضافه اینکه IE توی حجم اطلاعات GET محدودیت داره. به این صورت که وقتی اطلاعات فرم بصورت Query به URL اضافه میشه، طول این URL حداکثر میتونه 2083 کاراکتر باشه. در روش GET، چون اطلاعات فرم کد نمیشوند و ساده منتقل میشوند ، حجم کمتری دارند . اطلاعات فرم در این حالت، توسط متد GET_$ در صفحه مقصد قابل دریافت است . همچنین در روش GET از Upload خبری نیست، یعنی با GET نمیشه آپلود کرد.
پس از submit فرم در این روش ، اطلاعات فرم به صورتی که در کد زیر نمایش داده شده به آدرس صفحه اضافه شده و به صفحه مقصد منتقل می شوند : 
Syntax	http://www.developerstudio.ir/submit.php? fname = Ali & age = 26

*/
var http = require('http');
var url = require('url');
var fs = require('fs');

http.createServer(function (req, res) {
  res.writeHead(200, {'Content-Type': 'text/html'});
  /*
  var q = url.parse(req.url, true).query;
  var txt = q.year + " " + q.month;
  res.write(txt);
  */



  //************************** 
  var q2 = url.parse(req.url, true);
 // var filename = "." + q2.pathname;
  //res.write(q2.query.year);

  //console.log(q2.host); //returns 'localhost:8080'

  var qdata = q2.query; //returns an object: { year: 2017, month: 'february' }
  
  var txt1 =  "<br> qdata.year is  "+ qdata.year + "<br> qdata.month is " + qdata.month + " ";;
  var txt2 =  q2.year + "-------- " + q2.month;
  var txt3 =  "<br> q2.host is  "+ "<br> " + q2.host + " ";
  var txt4 =  "<br> q2.pathname is  "+ "<br> " + q2.pathname + " ";
  var txt5 =  "<br> q2.search is  "+ "<br> " + q2.search + " ";


  res.write(txt1 +txt2 +txt3 +txt4 +txt5  );
  //**************************




 res.end( );

}).listen(8080);