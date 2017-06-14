package main

import (
    "io"
    "net/http"
    "fmt"
    "os"
    "log"
    "github.com/line/line-bot-sdk-go/linebot"
)

func hello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if len(name) != 0 {
      io.WriteString(w, "Hello " + name)
    }
}

func callback(w http.ResponseWriter, req *http.Request) {
  bot, err := linebot.New(
		"6db72166ed2b37fbfd0a4a00f7bd01ac",
		"HeJ3LGMwxYtnGWdvtbG2pliVYpHxKSFyhJ5aoFOLORxhuMDqffoMGNhhlurDMepGX0IVHuM2sO67jtoK693UbdhrxroaNW/8ar74gB0aK5lZy4/P5kQ1vMHYWKNLWcluR/6XdNIo6QRM8n2jfwMOVwdB04t89/1O/w1cDnyilFU=",
	)
  events, err := bot.ParseRequest(req)
  if err != nil {
    if err == linebot.ErrInvalidSignature {
      w.WriteHeader(400)
    } else {
      w.WriteHeader(500)
    }
    return
  }
  for _, event := range events {
    if event.Type == linebot.EventTypeMessage {
      switch message := event.Message.(type) {
      case *linebot.TextMessage:
        fmt.Println(message.Text)
        if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
          log.Print(err)
        }
      }
    }
  }
}

func main() {

    port := os.Getenv("PORT")
    if port == "" {
  		log.Fatal("$PORT must be set")
  	}


    mux := http.NewServeMux()
    mux.HandleFunc("/", hello)
    mux.HandleFunc("/callback", callback)
    // http.ListenAndServe(":8000", mux)
    http.ListenAndServe(":"+port, mux)
}
