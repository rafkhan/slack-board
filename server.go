package main;

import (
  "io"
  "log"
  "net/url"
  "net/http"
  "io/ioutil"
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

func getBody(r *http.Request) string {
  body, err := ioutil.ReadAll(r.Body);

  if err != nil {
    return "";
  }

  return string(body);
}

func IsSlackbot(v url.Values) bool {
  return v["user_id"][0] == "USLACKBOT";
}

func HasTrigger(text string) bool {
  return text[0] == '~';
}

func slackHandler(w http.ResponseWriter, r *http.Request) {
  body := getBody(r);
  vals, err := url.ParseQuery(body);

  if err != nil || IsSlackbot(vals) {
    return;
  }

  text := vals["text"][0];
  if !HasTrigger(text) {
    return;
  }

  resp := text[1:];

  w.Write([]byte(resp));
}

func main() {
  http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:]);
  });

  http.HandleFunc("/websocket", websocketHandler);
  http.HandleFunc("/slackmessage", slackHandler);

  log.Println("Starting on 8080");
  log.Fatal(http.ListenAndServe(":8080", nil))
}


