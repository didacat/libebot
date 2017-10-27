// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)



var bot *linebot.Client
var isVote bool = false
func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
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
				if message.Text == "/vote" && isVote == false {
					log.Print("Start Vote")
					bot.PushMessage(groupID, linebot.NewTextMessage("Start Vote!")).Do()
					isVote = true
				}else if message.Text == "/stopvote" && isVote == true {
					log.Print("Stop Vote")
					bot.PushMessage(groupID, linebot.NewTextMessage("Stop Vote!")).Do()
					isVote = false
				}
				
				if message.Text[0:6] == "/vote " && isVote == true {

				//log.Print("user input" message.Text)
					res, err := bot.GetProfile(userID).Do();
					if err != nil {
						bot.PushMessage(groupID, linebot.NewTextMessage(message.Text)).Do()
					}
					log.Print(res.DisplayName)
					// log.Print(res.PicutureURL)
					// log.Print(res.StatusMessage)

					
				}

				if message.Text == "/pic" {
					bot.PushMessage(
						groupID, 
						linebot.NewImagemapMessage(
						"https://github.com/didacat/linebot/image/",
						"Imagemap alt text",
						linebot.ImagemapBaseSize{1040, 1040},
						linebot.NewURIImagemapAction("https://store.line.me/family/manga/en", linebot.ImagemapArea{0, 0, 520, 520}),
						linebot.NewURIImagemapAction("https://store.line.me/family/music/en", linebot.ImagemapArea{520, 0, 520, 520}),
						linebot.NewURIImagemapAction("https://store.line.me/family/play/en", linebot.ImagemapArea{0, 520, 520, 520}),
						linebot.NewMessageImagemapAction("URANAI!", linebot.ImagemapArea{520, 520, 520, 520}),
						),
					).Do(); 
				}
				log.Print(event.ReplyToken)
				log.Print(message.Text)
			}
		}
	}
}
