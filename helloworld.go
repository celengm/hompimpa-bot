package main

import (
    "io"
    "net/http"
    "fmt"
    "os"
    "log"
    "strings"
    "regexp"
    "strconv"
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
    fmt.Println(event)
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
    case linebot.EventTypePostback:
      postbackData := event.Postback.Data
      if strings.Contains(postbackData, "numberOfPlayers"){
        nPlayers, _ := strconv.Atoi(regexp.MustCompile("[0-9]+").FindString(event.Postback.Data))
        fmt.Println(nPlayers)
        if event.Source.Type == linebot.EventSourceTypeGroup {
          userChoiceMap[event.Source.GroupID] = make(map[string]string, nPlayers)
        } else if event.Source.Type == linebot.EventSourceTypeRoom {
          userChoiceMap[event.Source.RoomID] = make(map[string]string, nPlayers)
        }
        fmt.Println(event.Source.UserID)
        template := linebot.NewConfirmTemplate(
                              "Mau pilih apa?",
                              linebot.NewPostbackTemplateAction("Putih", "Putih", ""),
                              linebot.NewPostbackTemplateAction("Hitam", "Hitam", ""),
                              )
        if _, err := bot.ReplyMessage(
                              event.ReplyToken,
                              linebot.NewTemplateMessage("Hompimpa", template),
                              ).Do(); err != nil {
                                 log.Print(err)
                              }
      } else if (strings.Contains(postbackData, "Putih") || strings.Contains(postbackData, "Hitam")){
        if event.Source.Type == linebot.EventSourceTypeGroup {
          displayName := getUserProfile(event.Source.UserID)
          if _, ok := userChoiceMap[event.Source.GroupID][displayName]; !ok {
            userChoiceMap[event.Source.GroupID][displayName] = postbackData
            if _, err = bot.ReplyMessage(
                                  event.ReplyToken,
                                  linebot.NewTextMessage(displayName + " sudah milih"),
                                  ).Do(); err != nil {
                                      log.Print(err)
                                  }
          }
        } else if event.Source.Type == linebot.EventSourceTypeRoom {
          displayName := getUserProfile(event.Source.UserID)
          if _, ok := userChoiceMap[event.Source.RoomID][displayName]; !ok {
            userChoiceMap[event.Source.RoomID][displayName] = postbackData
            if _, err = bot.ReplyMessage(
                                  event.ReplyToken,
                                  linebot.NewTextMessage(displayName + " sudah milih"),
                                  ).Do(); err != nil {
                                      log.Print(err)
                                  }
          }
        } else if event.Source.Type == linebot.EventSourceTypeUser {
          if _, err = bot.ReplyMessage(
                                event.ReplyToken,
                                linebot.NewTextMessage("Kamu gak bisa main hompimpa sendiri"),
                                ).Do(); err != nil {
                                    log.Print(err)
                                }
        }
      }
    case linebot.EventTypeMessage:
      switch message := event.Message.(type) {
      case *linebot.TextMessage:
        if (strings.Contains(message.Text, "@bot") && strings.Contains(message.Text, "hompimpa")) {
          /*
          template := linebot.NewButtonsTemplate(
			                          "", "", "Berapa orang yg main?",
			                          linebot.NewPostbackTemplateAction("3", "numberOfPlayers=3", "3"),
			                          linebot.NewPostbackTemplateAction("4", "numberOfPlayers=4", "4"),
                                linebot.NewPostbackTemplateAction("5", "numberOfPlayers=5", "5"),
                                linebot.NewPostbackTemplateAction("6", "numberOfPlayers=6", "6"),
		                            )
                                */
          userChoiceMap[event.Source.GroupID] = make(map[string]string)
          template := linebot.NewConfirmTemplate(
                                "Mau pilih apa?",
                                linebot.NewPostbackTemplateAction("Putih", "Putih", ""),
                                linebot.NewPostbackTemplateAction("Hitam", "Hitam", ""),
                                )
		      if _, err := bot.ReplyMessage(
			                          event.ReplyToken,
			                          linebot.NewTemplateMessage("Hompimpa", template),
		                            ).Do(); err != nil {
			                             log.Print(err)
		                            }
        } else if (strings.Contains(message.Text, "@bot") && strings.Contains(message.Text, "selesai")) {
          //Get fewest choice
          if _, err = bot.ReplyMessage(
                                event.ReplyToken,
                                linebot.NewTextMessage(showUserChoice(event.Source.GroupID)),
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

func showUserChoice (group_id string) (string) {
  var blackText string = "Hitam: \n"
  var whiteText string = "Putih: \n"
  for k, _ := range userChoiceMap[group_id] {
    if userChoiceMap[group_id][k] == "Putih" {
      whiteText = whiteText + k + "\n"
    } else if userChoiceMap[group_id][k] == "Hitam" {
      blackText = blackText + k + "\n"
    }
  }
  return ("Users' choice: \n" + whiteText + "\n" + blackText)
}

func getUserProfile (user_id string) (string) {
  bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
  res, err := bot.GetProfile(user_id).Do();
  if err != nil {
    fmt.Println(err)
  }
  return res.DisplayName
}

func getFewestChoice (group_id string) (string) {
  var whiteChoice int = 0
  var blackChoice int = 0
  for k, _ := range userChoiceMap[group_id] {
    if userChoiceMap[group_id][k] == "Putih" {
      whiteChoice++
    } else if userChoiceMap[group_id][k] == "Hitam" {
      blackChoice++
    }
  }
  if whiteChoice < blackChoice {
    if whiteChoice == 1 {
      for k, _ := range userChoiceMap[group_id] {
        if userChoiceMap[group_id][k] == "Putih" {
          return "Yang menang adalah user: " + k
        }
      }
    } else if whiteChoice != 1 {
      return "Gak ada yang menang nih, ulang lagi ya"
    }
  } else if whiteChoice > blackChoice {
    if blackChoice == 1 {
      for k, _ := range userChoiceMap[group_id] {
        if userChoiceMap[group_id][k] == "Hitam" {
          return "Yang menang adalah user: " + k
        }
      }
    } else if blackChoice != 1 {
      return "Gak ada yang menang nih, ulang lagi ya"
    }
  } else if whiteChoice == blackChoice {
    return "Gak ada yang menang nih, ulang lagi ya"
  }
  return "Error"
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
