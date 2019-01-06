$(document).ready(function() {

  var socket = new WebSocket("ws://spyfall.carsonseese.com/api")

  $("#startselector").change(function() {
    if ($("input[name='startselection']:checked").val() === "create") {
      $("#create-game").show();
      $("#join-game").hide();
    } else {
      $("#create-game").hide();
      $("#join-game").show();
    }
  });

  $("#create-button").click(function(){
    console.log("Create Game")
  });

  $("#join-button").click(function(){
    console.log("Join Game")
  });
});