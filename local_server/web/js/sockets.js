var input = document.getElementById("input");
var output = document.getElementById("output");
var socket = new WebSocket("ws://localhost:4848/echo");

socket.onopen = function () {
    output.innerHTML += "Status: Connected\n";
};

socket.onmessage = function (e) {
    //output.innerHTML += "Server: " + e.data + "\n";
    console.log(e.data);
    var updatepc =  e.data.startsWith("pc:");
    if(updatepc){
        update_lstpc(e.data);
    }
};



function send() {
    socket.send(input.value);
    input.value = "";
  }

  setInterval(second_update, 1000);

 function second_update( )
 {
    socket.send("refresh");
    var numItems = $('.computers').length;
    if(numItems == 0){
      console.log("reload computers");
      socket.send("pagereload");
    }
 }


 function update_lstpc(pc){
     var split1 = pc.split("pc:");
     var split2 = split1[1].split(",");
    if($("#Card" + split2[0] + "PC").length == 0) {
        //it doesn't exist
        append_pc(split1[1]);
    }
 }

 function append_pc(data){
   var ip = "";
   var ssh = "";
   if(data.includes(",")){
   // console.log(data);
    var datasplit = data.split(",");
    ip = datasplit[0];
    ssh = datasplit[1];
   }else{
    ip = data;
   }


   ip = ip.replace(" ","");

   //var numItems = $('.Card' + ip + "PC").length;

  var classcount = document.getElementsByClassName('Card' + ip + 'PC');
   if(classcount.length == 0){
      $("#LST_PC").append(template_pc_card(ip,ssh));
   }

 }

 function template_pc_card(ip, ssh){
   var bg = "";
   if(ssh == "ready"){
    bg = "bg-success";
   }else{
    bg = "bg-secondary";
   }
	var pc_card = `<div class="col-lg-4 col-md-6 col-sm-12 computers"><div class="Card`+ip+`PC card text-white `+bg+` mb-3" style="width: 300px;">
<div class="card-header">`+ip+`</div>
<div class="card-body">
  <h5 class="card-title"></h5>
  <p class="card-text"></p>
</div>
<div class="card-footer text-white">
  <button id="Btn`+ip+`Ip" type="button" class="btn btn-dark">SSH</button>
  <button type="button" class="btn btn-dark"  data-toggle="modal" data-target="#cliModalFullscreen">CLI</button>
</div>
</div></div>`;
	return pc_card;
}