package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)



var bot *linebot.Client
var port =""
var addr =""
var isGameStart bool = false
func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port = os.Getenv("PORT")
	addr = fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
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
			
			userID := event.Source.UserID
			groupID := event.Source.GroupID
			RoomID := event.Source.RoomID
			log.Print(userID)
			log.Print(groupID)
			log.Print(RoomID)
			
			
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				/*if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK2!")).Do(); err != nil {
					log.Print(err)
				}*/
				if message.Text == "/dice" && isGameStart == false {
					log.Print("Start DiceGame")
					bot.PushMessage(groupID, linebot.NewTextMessage("Start DiceGame!")).Do()
					isGameStart = true
				}else if message.Text == "/stopdice" && isGameStart == true {
					log.Print("Stop DiceGame")
					bot.PushMessage(groupID, linebot.NewTextMessage("Stop DiceGame!")).Do()
					isGameStart = false
				}
				
				if(len(message.Text) > 6){
					if message.Text[0:6] == "/dice " && isGameStart == true {						
					//log.Print("user input" message.Text)
						res, err := bot.GetProfile(userID).Do();
						if err != nil {
							log.Print(err)
						}
						bot.PushMessage(groupID, linebot.NewTextMessage(message.Text)).Do()
						log.Print(res.DisplayName)
						// log.Print(res.PicutureURL)
						// log.Print(res.StatusMessage)						
					}
				}

					if message.Text == "/pic" {
					log.Print("Pic Receive")
					bot.PushMessage(
						groupID, 
						linebot.NewImageMessage(
						"https://raw.githubusercontent.com/didacat/linebot/master/images/1.png",
						"https://raw.githubusercontent.com/didacat/linebot/master/images/1.png"	,	
						)		,		
					).Do(); 
				}

				if message.Text == "/picall" {
					log.Print("Pic Receive")
					bot.PushMessage(
						groupID, 
						linebot.NewImageMessage(
							"https://raw.githubusercontent.com/didacat/linebot/master/images/2.png",
							"https://raw.githubusercontent.com/didacat/linebot/master/images/2.png",
							)		,	
					).Do(); 
				}
				log.Print(event.ReplyToken)
				log.Print(message.Text)
			}
		}
	}
}
