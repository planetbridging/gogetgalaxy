$( document ).ready(function() {

	console.log("Welcome to bewear");
  LoadFile("/web/menu.html");
});


function LoadFile(file){
  $.get(file, function(data){
      $(this).children("div:first").html(data);
  });
}