const http=require("http");
const urls=["technotip.com",
    "technotip.org",
    "www.capturecaptionapp.com"];

for(var i =0 ;i<urls.length; i++)
{

  ping(urls[i]);

}

function ping (url){
  var start =new Date();

  http.get({host:url},function (res) {
      console.log("url : " +url);


      console.log(" response time : "+(new Date()-start)+"ms");

  })

}