package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bytes"
	"io"
	// "path"
	// "errors"
	"strings"
	"io/ioutil"

	"image"
	"image/draw"
    "image/gif"
    "image/jpeg"
	"image/png"
	// "github.com/nfnt/resize"

	"github.com/line/line-bot-sdk-go/linebot"
)



var bot *linebot.Client
var port =""
var addr =""
var isVote bool = false
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
	log.Println("port = " + port)
	log.Println("Adde = " + addr)
	events, err := bot.ParseRequest(r)

	getImg("https://raw.githubusercontent.com/didacat/linebot/master/images/1.png")
	getImg("https://raw.githubusercontent.com/didacat/linebot/master/images/2.png")
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
				
				if(len(message.Text) > 6){
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

					src, err := GetImageObj("1.png")
					if err != nil {
						log.Println( "err:", err)
						log.Printf("Error1")
					}
					srcB := src.Bounds().Max

					src1, err := GetImageObj("2.png")
					if err != nil {
						log.Println( "err:", err)
						log.Printf("Error2")
					}
					src1B := src.Bounds().Max

					newWidth := srcB.X + src1B.X 

					newHeight := srcB.Y
					if src1B.Y > newHeight {
						newHeight = src1B.Y
					}
					
					des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
					srcWidth := srcB.X
					draw.Draw(des, des.Bounds(), src, src.Bounds().Min, draw.Over)                      //首先将一个图片信息存入jpg
					draw.Draw(des, image.Rect(srcWidth, 0, newWidth, src1B.Y), src1, image.ZP, draw.Over) //将另外一张图片信息存入jpg

					fSave, err := os.Create("7.png")
					if err != nil {
						log.Println( "err:", err)
						log.Printf("Error3")
					}
				
					defer fSave.Close()
				
					var opt jpeg.Options
					opt.Quality = 80
				
					// newImage := resize.Resize(1024, 0, des, resize.Lanczos3)
				
					err = jpeg.Encode(fSave, des, &opt) // put quality to 80%
					if err != nil {
						log.Println( "err:", err)
						log.Printf("Error4")
					}


					log.Print("Pic Receive")
					bot.PushMessage(
						groupID, 
						linebot.NewImageMessage(
							"https://didacat123.herokuapp.com/callback/7.png",
							"https://didacat123.herokuapp.com/callback/7.png"	,	
							)		,	
					).Do(); 
				}
				log.Print(event.ReplyToken)
				log.Print(message.Text)
			}
		}
	}
}


func GetImageObj(filePath string) (img image.Image, err error) {
    f1Src, err := os.Open(filePath)

    if err != nil {
        return nil, err
    }
    defer f1Src.Close()

    buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
    _, err = f1Src.Read(buff)

    if err != nil {
        return nil, err
    }

    filetype := http.DetectContentType(buff)

    fmt.Println(filetype)

    fSrc, err := os.Open(filePath)
    defer fSrc.Close()

    switch filetype {
    case "image/jpeg", "image/jpg":
        img, err = jpeg.Decode(fSrc)
        if err != nil {
            fmt.Println("jpeg error")
            return nil, err
        }

    case "image/gif":
        img, err = gif.Decode(fSrc)
        if err != nil {
            return nil, err
        }

    case "image/png":
        img, err = png.Decode(fSrc)
        if err != nil {
            return nil, err
        }
    default:
        return nil, err
    }
    return img, nil
}

func getImg(url string) (n int64, err error) {
    path := strings.Split(url, "/")
    var name string
    if len(path) > 1 {
        name = path[len(path)-1]
    }
    fmt.Println(name)
    out, err := os.Create(name)
    defer out.Close()
    resp, err := http.Get(url)
    defer resp.Body.Close()
    pix, err := ioutil.ReadAll(resp.Body)
    n, err = io.Copy(out, bytes.NewReader(pix))
    return

}
