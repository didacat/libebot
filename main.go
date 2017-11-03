package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"math/rand"
	"strconv"  
	"strings"
	"time"  
	"github.com/line/line-bot-sdk-go/linebot"
)



var bot *linebot.Client
var port =""
var addr =""
var isGameStart bool = false
var UserNameSlice []string	//玩家名稱
var UserIDSlice []string	//玩家ID
var UserAnsMap = make(map[string]string) //玩家ID跟骰子數值的MAP表
var UserCanSpeakSlice []bool //玩家是否能說話
var UserDiceCount []int //玩家的骰子數量
var WhoRound int = 0	//輪到誰的INDEX
var m_groupID =""

var test = 6

type DiceValue struct {	
	Values []string
}
 
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
				//如果訊息來自 有發言權的 使用者
				if (groupID == "" && UserIDSlice[WhoRound] == userID){
					log.Print(userID + "講話啦~~")
					log.Print("WhoRound == " + strconv.Itoa(WhoRound))
					NextUserRound := 0
					UserAnser := ""
					if(WhoRound + 1 >= len(UserIDSlice)){
						NextUserRound = 0
					}else{
						NextUserRound += 1
					}
					
					if(message.Text =="1"){
						UserAnser = "單"
						NewValue := ""
						for _, value := range UserIDSlice {
							
							NewValue = strings.TrimRight(UserAnsMap[value],"1")
							NewValue = strings.TrimRight(UserAnsMap[value],"3")
							NewValue = strings.TrimRight(UserAnsMap[value],"5")
							// for _, DiceValue := range UserAnsMap[value] {
							// 	if(strings.EqualFold(DiceValue,"2") || strings.EqualFold(DiceValue,"4") || strings.EqualFold(DiceValue,"6")){
							// 		NewValue += DiceValue
							// 	}
							// }
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])
							
						}
					}else if(message.Text =="2"){
						UserAnser = "雙"
						test = test - 2
					}else if(message.Text =="3"){
						UserAnser = "大"
						test = test - 3
					}else if(message.Text =="4"){
						UserAnser = "小"
						test = test - 2
					}else if(message.Text =="5"){
						UserAnser = "紅"
						test = test - 3
					}else if(message.Text =="6"){
						UserAnser = "黑"
						test = test - 2
					}


					bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound] + "選擇把  " + UserAnser +"拿掉\n換" + UserNameSlice[NextUserRound] + "的回合囉")).Do()
					if(WhoRound + 1 >= len(UserIDSlice)){
						WhoRound = 0
					}else{
						WhoRound += 1
					}
					
					//刪除 MAP 表
					// for _, value := range UserIDSlice {
					// 	delete(UserAnsMap,value)
					// }
					
					for _, value := range UserIDSlice {
						rand.Seed(time.Now().UnixNano())  						
						SliceValue := make([]int, test)
						NumerString := ""
						for _, element := range SliceValue {
							element = rand.Intn(6)  
							element = element + 1 
							NumerString = NumerString + strconv.Itoa(element)  
						}
						// UserAnsMap[NumerString] = value
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
					bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage("請決定你要喊的骰子\n 1)單 \n 2)雙 \n 3)大\n 4)小 \n 5)紅 \n 6)黑" )).Do()
					//發給群組 現在是誰的回合
					// bot.PushMessage(groupID, linebot.NewTextMessage("現在是 " +  UserNameSlice[WhoRound] + "的回合")).Do()

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
							UserDiceCount = append(UserDiceCount, 6)	//初始都給玩家 6顆骰子
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
							UserAnsMap[value] = NumerString
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
						bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage("請決定你要喊的骰子\n 1)單 \n 2)雙 \n 3)大\n 4)小 \n 5)紅 \n 6)黑" )).Do()
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
