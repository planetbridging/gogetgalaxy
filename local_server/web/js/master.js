$( document ).ready(function() {

	console.log("Welcome to ggg");
  //LoadFile("/web/menu.html");
  //refresh_countries()

  $( "button" ).click(function() {
    BtnPress($(this).attr('id'));
  });

  $("#LST_PC").html("");

});


function BtnPress(id){
  console.log("BTN: " + id);

  if(id == "BtnFindAll"){
    LoadFile("/findall");
  }

}


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