$(document).ready(function(){
	var conn;
	var message = document.getElementById("msg");
	
		//takes html from hidden div and builds message from it.
	function newmessage(msg){
		$('#messageblock #username').text(msg.user + " >> ");
		$('#messageblock #message').text(msg.text);
		
		$('#log').append($('#messageblock').html());
		
		$('#messageblock #username').text('');
		$('#messageblock #message').text('');
	}
		/*
			when textarea is submitted, we check to see if we have 
			either lost connection or no message has been typed.  
			If one of these fail, we exit and do nothing.  If both
			succeed, we send message using websocket;
		*/
	$('#form').on('submit', function(){
		
        if (!message.value || !conn) {
            return false;
        }
        conn.send(message.value);
        message.value = "";
        return false;
	});
		/*
			initialize connection with websocket and allow server 
			to send messages to client through websocket.
		*/
	if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function () {
            var ping = {user: "server", text: "Connection Closed."};
       		newmessage(ping);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            var ping = JSON.parse(messages);
            newmessage(ping);
        };
    } else {
		var ping = {user: "server", text: "Your browser does not support WebSockets."};
        newmessage(ping);
    }
});