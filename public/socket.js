/*global window, WebSocket, document*/

window.onload = function() {
  'use strict';

  var h1 = document.getElementById('message');

  var sock = new WebSocket('ws://127.0.0.1:8080/websocket');
  sock.onopen = function() {
    setTimeout(function() {
      sock.send('asd');
    }, 500);
  };

  sock.onmessage = function(e) {
    h1.innerHTML = e.data;
  };
};
