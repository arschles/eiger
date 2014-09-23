function main() {
  var ws = new WebSocket("ws://127.0.0.1:4492/stream/");

  ws.onopen = function(){
    console.log("Socket has been opened!");
  };

  ws.onmessage = function(message) {
    listener(message.data)
  };

  function listener(data) {
    var messageObj = data;
    console.log("received data from websocket: ", messageObj);
    $("#event-list").append("<div>"+messageObj+"</div>");
  }

}

$(document).ready(function() {
  main();
});
