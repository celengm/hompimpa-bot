package main

import (
    "io"
    "net/http"
    "fmt"
    "os"
    "log"
    "strings"
    "github.com/line/line-bot-sdk-go/linebot"
)

var userChoiceMap map[string]map[string]string

func hello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if len(name) != 0 {
      io.WriteString(w, "Hello " + name)
    }
}

func callback(w http.ResponseWriter, req *http.Request) {
  bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
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
    fmt.Println(event.Type)
    switch event_type := event.Type; event_type {
    case linebot.EventTypeJoin:
      if event.Source.Type == linebot.EventSourceTypeGroup {
        if _, err = bot.PushMessage(
                      event.Source.GroupID,
                      linebot.NewTextMessage("Halo, sapa aku dong dengan kirim \"Hai, @bot\""),
                      ).Do(); err != nil {
                          log.Print(err)
                      }
      } else if event.Source.Type == linebot.EventSourceTypeRoom {
        if _, err = bot.PushMessage(
                      event.Source.RoomID,
                      linebot.NewTextMessage("Halo, sapa aku dong dengan kirim \"Hai, @bot\""),
                      ).Do(); err != nil {
                          log.Print(err)
                      }
      }
    // case linebot.EventTypePostback:
    //   if event.Source.Type == linebot.EventSourceTypeGroup {
    //     runHompimpaGame(event.Source.GroupID, event.Source.UserId, event.Postback.Data)
    //   } else if event.Source.Type == linebot.EventSourceTypeRoom {
    //     runHompimpaGame(event.Source.RoomID, event.Source.UserId, event.Postback.Data)
    //   } else if event.Source.Type == linebot.EventSourceTypeUser {
    //     if _, err = bot.ReplyMessage(
    //                           event.ReplyToken,
    //                           linebot.NewTextMessage("Kamu cuma bisa main hompimpa di group atau room"),
    //                           ).Do(); err != nil {
    //                               log.Print(err)
    //                           }
    //   }
    case linebot.EventTypeMessage:
      switch message := event.Message.(type) {
      case *linebot.TextMessage:
        if (strings.Contains(message.Text, "@bot") && strings.Contains(message.Text, "hompimpa")) {
          template := linebot.NewButtonsTemplate(
			                          "", "", "Mau pilih apa?",
			                          linebot.NewPostbackTemplateAction("3", "numberOfPlayers=3", "3"),
			                          linebot.NewPostbackTemplateAction("4", "numberOfPlayers=4", "4"),
                                linebot.NewPostbackTemplateAction("5", "numberOfPlayers=5", "5"),
                                linebot.NewPostbackTemplateAction("6", "numberOfPlayers=6", "6"),
		                            )
		      if _, err := bot.ReplyMessage(
			                          event.ReplyToken,
			                          linebot.NewTemplateMessage("Hompimpa", template),
		                            ).Do(); err != nil {
			                             log.Print(err)
		                            }
        } else {
          if _, err = bot.ReplyMessage(
                                event.ReplyToken,
                                linebot.NewTextMessage("Wah, ak gak ngerti kamu mau apa, aku cuma bisa kasih game hompimpa"),
                                ).Do(); err != nil {
                                    log.Print(err)
                                }
        }
      }
    }
  }
}

func runHompimpaGame(group_id, user_id, choice string) {

}

func main() {
    userChoiceMap = make(map[string]map[string]string)

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
