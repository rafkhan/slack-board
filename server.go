package main;

import (
  "io"
  "log"
  "net/http"
  websocket "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader {
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
  CheckOrigin: func(r *http.Request) bool { return true },
}

func readWriteSocket(conn *websocket.Conn) {
  for {
    messageType, r, err := conn.NextReader()
    if err != nil {
      return
    }

    w, err := conn.NextWriter(messageType)
    if err != nil {
      log.Println(err)
    }

    if _, err := io.Copy(w, r); err != nil {
      log.Println(err)
    }

    if err := w.Close(); err != nil {
      log.Println(err)
    }
  }
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
  conn, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println(err);
    return;
  }

  go readWriteSocket(conn);
}

func main() {
  http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:]);
  });

  http.HandleFunc("/foo", websocketHandler);

  log.Println("Starting on 8080");
  log.Fatal(http.ListenAndServe(":8080", nil))
}


