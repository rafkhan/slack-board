/*global window, WebSocket, document*/

window.onload = function() {
  'use strict';

  var h1 = document.getElementById('message');

  var host = window.location.host;

  var sock = new WebSocket('ws://' + host + '/websocket');
  sock.onopen = function() {
    setTimeout(function() {
      sock.send('asd');
    }, 500);
  };

  sock.onmessage = function(e) {
    h1.innerHTML = e.data;
  };
};
