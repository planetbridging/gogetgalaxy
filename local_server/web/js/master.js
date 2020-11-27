$( document ).ready(function() {

	console.log("Welcome to ggg");
  //LoadFile("/web/menu.html");
  refresh_countries()
});


function LoadFile(file){
  $.get(file, function(data){
      $(this).children("div:first").html(data);
  });
}

function refresh_countries(){
  var time = 1;

  var interval = setInterval(function() { 
    if (time <= 3) { 

      $( "#LstCountires" ).load("/countries");

        time++;
    }
    else { 
        clearInterval(interval);
    }
  }, 5000);
}