/*global window, WebSocket, console*/

window.onload = function() {
  'use strict';

  var sock = new WebSocket('ws://127.0.0.1:8080/foo');
  sock.onopen = function() {
    sock.send('asd');
  };

  sock.onmessage = function(e) {
    console.log(e.data);
  };
};
