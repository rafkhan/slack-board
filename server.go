package main;

import (
  "io"
  "log"
  "bytes"
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

func readWriteSocket(ch chan []byte, conn *websocket.Conn) {
  for {
    // wait on channel
    data := <-ch;
    r := bytes.NewReader(data);

    w, err := conn.NextWriter(1); //not sure why 1
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

func websocketHandler(ch chan []byte) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {

    log.Println("Connected");

    conn, err := upgrader.Upgrade(w, r, nil);

    if err != nil {
      log.Println(err);
      return;
    }

    go readWriteSocket(ch, conn);
  };
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

func slackHandler(ch chan []byte) func(w http.ResponseWriter, r *http.Request){
  return func(w http.ResponseWriter, r *http.Request) {
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

    log.Println(body);
    log.Println(resp);

    ch <- []byte(resp);
  };
}

func main() {
  http.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, r.URL.Path[1:]);
  });

  ch := make(chan []byte);

  http.HandleFunc("/websocket", websocketHandler(ch));
  http.HandleFunc("/slackmessage", slackHandler(ch));

  log.Println("Starting on 8080");
  log.Fatal(http.ListenAndServe(":8080", nil))
}


