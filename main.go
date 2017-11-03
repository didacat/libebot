package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"math/rand"
	"strconv"  
	"time"  
	"github.com/line/line-bot-sdk-go/linebot"
)



var bot *linebot.Client
var port =""
var addr =""
var isGameStart bool = false
var UserNameSlice []string
var UserIDSlice []string
var UserAnsMap = make(map[string]string)
var UserCanSpeakSlice []bool
var WhoRound int = 0
var m_groupID =""
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
				//如果訊息來自 使用者
				if (groupID == ""){
					log.Print(userID + "講話啦~~")
					bot.PushMessage(m_groupID, linebot.NewTextMessage(userID + "講話啦~~  " + message.Text)).Do()
				}else{ //訊息來自 群組
					if message.Text == "/dice" && isGameStart == false {
						log.Print("Start DiceGame")
						bot.PushMessage(groupID, linebot.NewTextMessage("Start DiceGame!")).Do()
						isGameStart = true
					}else if message.Text == "/dicestop" && isGameStart == true {
						log.Print("Stop DiceGame")
						UserNameSlice = UserNameSlice[:0]
						UserIDSlice = UserIDSlice[:0]
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
							UserName := message.Text[6:len(message.Text)]
							UserNameSlice= append(UserNameSlice, UserName)	//紀錄玩家輸入的名稱 之後推撥會顯示玩家名稱的回合
							UserIDSlice = append(UserIDSlice, userID)	//紀錄玩家的ID 便於後續發送圖片
							UserCanSpeakSlice = append(UserCanSpeakSlice, false) //玩家有無回答的權限
							TotalUser := ""
							for _, value := range UserNameSlice {
								TotalUser += value + ","
							}
	
							bot.PushMessage(groupID, linebot.NewTextMessage(UserName + " 已加入遊戲\n" + "目前玩家有 : " + TotalUser )).Do()
							log.Print(res.DisplayName)
							log.Print(UserName)					
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
	
					if message.Text == "/dicestart" {
						log.Print("dicestart Receive")
						m_groupID = groupID
						for _, value := range UserIDSlice {
							rand.Seed(time.Now().UnixNano())  
							arrValue := [...]int{1,2,3,4,5,6}
							NumerString := ""
							for _, element := range arrValue {
								element = rand.Intn(6)  
								element = element + 1 
								NumerString = NumerString + strconv.Itoa(element)  
							}
							UserAnsMap[NumerString] = value
							log.Print("NumerString = " + NumerString)
							log.Print(UserAnsMap)
							log.Print(value)
							//發送給玩家圖片
							bot.PushMessage(
								value, 
								linebot.NewImageMessage(
									"https://jenny-web.herokuapp.com/dice/merge/"+ NumerString +"/0/564531635164",
									"https://jenny-web.herokuapp.com/dice/merge/"+ NumerString +"/0/564531635164",
									)		,	
							).Do(); 
						}
						//讓第一位玩家 可以回答
						bot.PushMessage(UserIDSlice[0], linebot.NewTextMessage("請決定你要喊的骰子\n 1)單雙 \n 2)大小 \n 3)紅黑" )).Do()
						//發給群組 現在是誰的回合
						bot.PushMessage(groupID, linebot.NewTextMessage("現在是 " +  UserNameSlice[0] + "的回合")).Do()
						
					}
				}
				
				// log.Print(event.ReplyToken)
				log.Print(message.Text)
			}
		}
	}
}
