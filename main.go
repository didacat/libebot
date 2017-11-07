package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var port = ""
var addr = ""
var isGameStart bool = false     //猜骰子
var isBlowGameStart bool = false //吹牛
var isDice bool = false
var isGuess bool = false
var isBlow bool = false
var UserNameSlice []string                        //玩家名稱
var UserIDSlice []string                          //玩家ID
var UserAnsMap = make(map[string]string)          //玩家ID跟骰子數值的MAP表
var UserCanSpeakSlice []bool                      //玩家是否能說話
var UserDiceCount []int                           //玩家的骰子數量
var WhoRound int = 0                              //輪到誰的INDEX
var AllDiceValueAndCount = make(map[int]int)      //所有玩家骰子數值跟數量的MAP表 一代表任何數
var AllDiceValueAndCountNoOne = make(map[int]int) //所有玩家骰子數值跟數量的MAP表 一被喊掉後 不代表任何數
var isUseOne bool = false                         //點數1是某被喊掉了
var NeedDiceCount = 0                             //最少要喊的骰子數量
var NeedDiceValue = 0                             //最少要喊的骰子數值
var NextUserRound = 0
var PreUserRound = 0
var m_groupID = ""

var diceCount = 10
var blowdiceCount = 6

// type DiceValue struct {
// 	Values []string
// }

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

				//如果訊息來自 有發言權的 使用者 並且在玩猜骰子
				if groupID == "" && UserIDSlice[WhoRound] == userID && isGuess {
					log.Print(userID + "講話啦~~" + message.Text)
					log.Print("WhoRound == " + strconv.Itoa(WhoRound))
					//如果是上一輪玩家剩一顆骰子 補發照片給她
					log.Print("UserIDSlice.length == " + strconv.Itoa(len(UserIDSlice)))
					log.Print("UserAnsMap.length == " + strconv.Itoa(len(UserAnsMap)))
					if UserIDSlice[WhoRound] == userID && len(UserAnsMap[UserIDSlice[WhoRound]]) == 1 {
						bot.PushMessage(
							UserIDSlice[WhoRound],
							linebot.NewImageMessage(
								"https://jenny-web.herokuapp.com/dice/merge/"+UserAnsMap[UserIDSlice[WhoRound]]+"/0/564531635164",
								"https://jenny-web.herokuapp.com/dice/merge/"+UserAnsMap[UserIDSlice[WhoRound]]+"/0/564531635164",
							),
						).Do()
					}

					UserAnser := ""
					if WhoRound+1 >= len(UserIDSlice) {
						NextUserRound = 0
					} else {
						NextUserRound += 1
					}

					if message.Text == "1" {
						UserAnser = "單"
						for _, value := range UserIDSlice {
							NewValue := ""
							for _, DiceValue := range UserAnsMap[value] {
								if strings.EqualFold(string(DiceValue), "2") || strings.EqualFold(string(DiceValue), "4") || strings.EqualFold(string(DiceValue), "6") {
									// log.Print(string(DiceValue))
									NewValue += string(DiceValue)
								}
							}
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])

						}
					} else if message.Text == "2" {
						UserAnser = "雙"
						for _, value := range UserIDSlice {
							NewValue := ""
							for _, DiceValue := range UserAnsMap[value] {
								if strings.EqualFold(string(DiceValue), "1") || strings.EqualFold(string(DiceValue), "3") || strings.EqualFold(string(DiceValue), "5") {
									// log.Print(string(DiceValue))
									NewValue += string(DiceValue)
								}
							}
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])

						}
					} else if message.Text == "3" {
						UserAnser = "大"
						for _, value := range UserIDSlice {
							NewValue := ""
							for _, DiceValue := range UserAnsMap[value] {
								if strings.EqualFold(string(DiceValue), "1") || strings.EqualFold(string(DiceValue), "2") || strings.EqualFold(string(DiceValue), "3") {
									// log.Print(string(DiceValue))
									NewValue += string(DiceValue)
								}
							}
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])

						}
					} else if message.Text == "4" {
						UserAnser = "小"
						for _, value := range UserIDSlice {
							NewValue := ""
							for _, DiceValue := range UserAnsMap[value] {
								if strings.EqualFold(string(DiceValue), "4") || strings.EqualFold(string(DiceValue), "5") || strings.EqualFold(string(DiceValue), "6") {
									// log.Print(string(DiceValue))
									NewValue += string(DiceValue)
								}
							}
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])

						}
					} else if message.Text == "5" {
						UserAnser = "紅"
						for _, value := range UserIDSlice {
							NewValue := ""
							for _, DiceValue := range UserAnsMap[value] {
								if strings.EqualFold(string(DiceValue), "2") || strings.EqualFold(string(DiceValue), "3") || strings.EqualFold(string(DiceValue), "5") || strings.EqualFold(string(DiceValue), "6") {
									// log.Print(string(DiceValue))
									NewValue += string(DiceValue)
								}
							}
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])

						}
					} else if message.Text == "6" {
						UserAnser = "黑"
						for _, value := range UserIDSlice {
							NewValue := ""
							for _, DiceValue := range UserAnsMap[value] {
								if strings.EqualFold(string(DiceValue), "1") || strings.EqualFold(string(DiceValue), "4") {
									// log.Print(string(DiceValue))
									NewValue += string(DiceValue)
								}
							}
							log.Print(NewValue)
							UserAnsMap[value] = NewValue
							log.Print(UserAnsMap[value])

						}
					} else {
						bot.PushMessage(userID, linebot.NewTextMessage("請輸入 1 ~ 6 , 其中一數字")).Do()
						return
					}

					bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+"選擇把  "+UserAnser+"拿掉\n換"+UserNameSlice[NextUserRound]+"的回合囉")).Do()
					Msg := ""
					for i, value := range UserIDSlice {
						Msg = Msg + UserNameSlice[i] + "剩下 " + strconv.Itoa(len(UserAnsMap[value])) + "顆骰子囉~\n"
					}
					bot.PushMessage(m_groupID, linebot.NewTextMessage(Msg)).Do()
					PreUserRound = WhoRound
					if WhoRound+1 >= len(UserIDSlice) {
						WhoRound = 0
					} else {
						WhoRound += 1
					}
					//

					Tie := false //是否平手
					for i, value := range UserIDSlice {
						rand.Seed(time.Now().UnixNano())
						SliceValue := make([]int, len(UserAnsMap[value]))
						NumerString := ""
						for _, element := range SliceValue {
							element = rand.Intn(6)
							element = element + 1
							NumerString = NumerString + strconv.Itoa(element)
						}
						UserAnsMap[value] = NumerString
						log.Print("NumerString = " + NumerString)
						log.Print(UserAnsMap)
						log.Print(value)
						//如果玩家還有骰子的話 發送新的骰子圖片給玩家
						if len(UserAnsMap[value]) > 0 {
							//如果玩家只剩下一顆骰子的話
							if len(UserAnsMap[value]) == 1 {
								//若玩家有發言權 先問玩家要選甚麼類型 再發送圖片
								if UserIDSlice[WhoRound] == value {

								} else {
									//若玩家沒有發言權 直接送圖片
									bot.PushMessage(value, linebot.NewTextMessage("==此局為新的一局牌面==")).Do()
									bot.PushMessage(
										value,
										linebot.NewImageMessage(
											"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
											"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
										),
									).Do()
								}

							} else {
								SomeBodyOut := false
								for _, value := range UserIDSlice {
									if len(UserAnsMap[value]) == 0 {
										SomeBodyOut = true
									}
								}
								//如果沒有人出局的話 就繼續發送圖片
								if SomeBodyOut == false {
									bot.PushMessage(value, linebot.NewTextMessage("==此局為新的一局牌面==")).Do()
									bot.PushMessage(
										value,
										linebot.NewImageMessage(
											"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
											"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
										),
									).Do()
								}

							}
						} else {
							//判斷其他玩家是否也沒有骰子了
							iCount := 0
							for _, value := range UserIDSlice {
								if len(UserAnsMap[value]) == 0 {
									iCount++
								}
							}
							if iCount == 2 {
								Tie = true
							}
							if Tie {
								for _, value := range UserIDSlice {
									bot.PushMessage(
										value,
										linebot.NewImageMessage(
											"https://gss2.bdstatic.com/9fo3dSag_xI4khGkpoWK1HF6hhy/baike/w%3D268%3Bg%3D0/sign=1ccf6bdd9a52982205333ec5eff11cf6/d000baa1cd11728b1dfcdb38c2fcc3cec2fd2cd2.jpg",
											"https://gss2.bdstatic.com/9fo3dSag_xI4khGkpoWK1HF6hhy/baike/w%3D268%3Bg%3D0/sign=1ccf6bdd9a52982205333ec5eff11cf6/d000baa1cd11728b1dfcdb38c2fcc3cec2fd2cd2.jpg",
										),
									).Do()
								}
							} else {
								//骰子沒了 對失敗者發送失敗照片 贏家發送贏照
								for _, value := range UserIDSlice {
									if len(UserAnsMap[value]) > 0 {
										bot.PushMessage(
											value,
											linebot.NewImageMessage(
												"https://www.jiuwa.net/tuku/20170627/wndhdTKu.jpg",
												"https://www.jiuwa.net/tuku/20170627/wndhdTKu.jpg",
											),
										).Do()
									} else {
										bot.PushMessage(
											value,
											linebot.NewImageMessage(
												"https://i.ytimg.com/vi/V3fEhrP_9xc/maxresdefault.jpg",
												"https://i.ytimg.com/vi/V3fEhrP_9xc/maxresdefault.jpg",
											),
										).Do()
									}
								}
							}

							if Tie {
								//發送給群組 告知大家平手 結束遊戲
								bot.PushMessage(m_groupID, linebot.NewTextMessage("此局平手~~請大蝦重新來過!!")).Do()
								//清空所有資料 重新開始一局
								log.Print("END DiceGame")
								break
							} else {
								//發送給群組 告知有人輸了 結束遊戲
								bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[i]+"被清光光了~")).Do()
								//清空所有資料 重新開始一局
								log.Print("END DiceGame")
							}

							isGameStart = false
							WhoRound = 0
							PreUserRound = 0
							NextUserRound = 0
						}
					}
					SomeBodyOut := false
					for _, value := range UserIDSlice {
						if len(UserAnsMap[value]) == 0 {
							SomeBodyOut = true
						}
					}
					//如果沒有人出局的話 就繼續發送訊息給下一位
					if SomeBodyOut == false {
						//讓下一回合有骰子的玩家 可以回答
						if len(UserAnsMap[UserIDSlice[WhoRound]]) > 1 {
							bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage("請決定你要喊的骰子\n 1)單 \n 2)雙 \n 3)大\n 4)小 \n 5)紅 \n 6)黑")).Do()
						} else if len(UserAnsMap[UserIDSlice[WhoRound]]) == 1 {
							bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage("你只剩下一顆骰子\n已經幫您重新搖骰了\n但先不讓你看\n請先決定你要喊的骰子\n 1)單 \n 2)雙 \n 3)大\n 4)小 \n 5)紅 \n 6)黑\n盲骰後會讓你看骰子的~")).Do()
						} else {
							bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage("你已經輸了~~")).Do()
						}
					}

				} else if groupID == "" && UserIDSlice[WhoRound] == userID && isBlow { //如果訊息來自 有發言權的 使用者 並且在玩吹牛
					log.Print(userID + "講話啦~~" + message.Text)
					log.Print("WhoRound == " + strconv.Itoa(WhoRound))
					UserSpeakDiceCount := 0
					UserSpeakDiceValue := 0
					//拆解玩家回答的字串 16/3  == 16個3  或者是 抓
					for i, value := range strings.Split(message.Text, "/") {
						if i == 0 {
							UserSpeakDiceCount, err = strconv.Atoi(value)
						} else if i == 1 {
							UserSpeakDiceValue, err = strconv.Atoi(value)
							if UserSpeakDiceValue == 1 {
								isUseOne = true
							}
						}
					}
					log.Print("UserSpeakDiceCount == " + strconv.Itoa(UserSpeakDiceCount))
					log.Print("UserSpeakDiceValue == " + strconv.Itoa(UserSpeakDiceValue))
					log.Print("NeedDiceCount == " + strconv.Itoa(NeedDiceCount))
					log.Print("NeedDiceValue == " + strconv.Itoa(NeedDiceValue))

					isFirstRound := false
					//第一個人喊的時候
					if NeedDiceValue == 0 {
						isFirstRound = true
					}
					//判斷 是否有符合規則
					isBigger := false
					if UserSpeakDiceCount > NeedDiceCount && UserSpeakDiceValue > 0 && UserSpeakDiceValue < 7 && NeedDiceValue != 0 {
						NeedDiceCount = UserSpeakDiceCount
						NeedDiceValue = UserSpeakDiceValue
						isBigger = true
					} else if UserSpeakDiceCount == NeedDiceCount && UserSpeakDiceValue > NeedDiceValue && UserSpeakDiceValue > 0 && UserSpeakDiceValue < 7 && NeedDiceValue != 0 {
						NeedDiceCount = UserSpeakDiceCount
						NeedDiceValue = UserSpeakDiceValue
						isBigger = true
					}
					SomeBodyOut := false //是否抓到
					//如果符合規則 則換下一個玩家作答
					if isBigger {
						bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+" 吹了   "+strconv.Itoa(UserSpeakDiceCount)+" 個 "+strconv.Itoa(UserSpeakDiceValue)+"\n 現在換"+UserNameSlice[NextUserRound]+"的回合囉")).Do()
					} else if message.Text == "抓" {
						bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+" 選擇抓爆 "+UserNameSlice[PreUserRound])).Do()
						//如果是第一位玩家 喊的時候 不能喊抓
						if isFirstRound {
							bot.PushMessage(userID, linebot.NewTextMessage("請輸入 x/x 這種格式\n並且要大於"+strconv.Itoa(NeedDiceCount)+"個"+strconv.Itoa(NeedDiceValue)+"的牌面")).Do()
						} else {
							//判斷牌面是否有喊中
							//有喊過1的情況
							if isUseOne {
								if NeedDiceCount > AllDiceValueAndCountNoOne[NeedDiceValue] {
									SomeBodyOut = true
									bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+" 抓到了!!! "+UserNameSlice[PreUserRound]+"在吹牛")).Do()
								} else {
									bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+" 沒抓到... ")).Do()
								}

								TotalMsg := ""
								for i, _ := range UserNameSlice {
									TotalMsg = UserNameSlice[i] + " = " + UserAnsMap[UserNameSlice[i]]
								}
								bot.PushMessage(m_groupID, linebot.NewTextMessage("所有玩家的牌面是 : \n "+TotalMsg+"\n 總共有"+strconv.Itoa(AllDiceValueAndCount[UserSpeakDiceValue])+"個"+strconv.Itoa(UserSpeakDiceValue))).Do()
							} else { //沒有喊過1的情況
								if NeedDiceCount > AllDiceValueAndCount[NeedDiceValue] {
									SomeBodyOut = true
									bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+" 抓到了!!! "+UserNameSlice[PreUserRound]+"在吹牛")).Do()
								} else {
									bot.PushMessage(m_groupID, linebot.NewTextMessage(UserNameSlice[WhoRound]+" 沒抓到... ")).Do()
								}

								TotalMsg := ""
								for i, _ := range UserNameSlice {
									TotalMsg = UserNameSlice[i] + " = " + UserAnsMap[UserNameSlice[i]]
								}
								bot.PushMessage(m_groupID, linebot.NewTextMessage("所有玩家的牌面是 : \n "+TotalMsg+"\n 總共有"+strconv.Itoa(AllDiceValueAndCount[UserSpeakDiceValue])+"個"+strconv.Itoa(UserSpeakDiceValue))).Do()
							}
						}
					} else {
						bot.PushMessage(userID, linebot.NewTextMessage("請輸入 x/x 這種格式\n並且要大於上一家喊的牌面 \n或是輸入 抓")).Do()
						return
					}

					PreUserRound = WhoRound
					if WhoRound+1 >= len(UserIDSlice) {
						WhoRound = 0
					} else {
						WhoRound += 1
					}
					//

					// SomeBodyOut := false //是否抓到

					//如果沒有人出局的話 就繼續發送訊息給下一位
					if SomeBodyOut == false {
						//讓下一回合的玩家 開始吹
						if len(UserAnsMap[UserIDSlice[WhoRound]]) > 1 {
							bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage(" 請決定你要喊的骰子點數及數量 格式為 x/x (e.g.: 3/3) 或是輸入 抓")).Do()
						}
					}
				} else { //訊息來自 群組
					if message.Text == "/dice" && isDice == false {
						log.Print("Start DiceGame")
						bot.PushMessage(groupID, linebot.NewTextMessage("Start DiceGame!")).Do()
						isDice = true
					} else if message.Text == "/gg" && isDice == true {
						log.Print("Stop DiceGame")
						UserNameSlice = UserNameSlice[:0]
						UserIDSlice = UserIDSlice[:0]
						bot.PushMessage(groupID, linebot.NewTextMessage("Stop DiceGame!")).Do()
						isDice = false
						isGuess = false
						isBlow = false
						isGameStart = false
						isBlowGameStart = false
						isUseOne = false
						NeedDiceValue = 1
						NeedDiceCount = 0
						WhoRound = 0
						PreUserRound = 0
						NextUserRound = 0
						for i, _ := range AllDiceValueAndCountNoOne {
							AllDiceValueAndCountNoOne[i+1] = 0
						}
						for i, _ := range AllDiceValueAndCount {
							AllDiceValueAndCount[i+1] = 0
						}
					}

					if len(message.Text) > 6 && !isGameStart && !isBlowGameStart {
						if message.Text[0:6] == "/dice " && isDice && !isBlow {
							isGuess = true
							//先判斷名稱是否重複 是的話幫玩家修改名字 或是 禁止玩家取同樣名稱
							UserName := message.Text[6:len(message.Text)]
							if len(UserIDSlice) >= 2 {
								UserDuplicate := false //玩家是否重複
								//判斷玩家ID有無重複 有的話變成修改名稱
								for i, value := range UserIDSlice {
									//判斷玩家是否相同名字
									if UserNameSlice[i] == UserName {
										bot.PushMessage(groupID, linebot.NewTextMessage(UserName+" 玩家名稱已存在 ")).Do()
										return
									}

									if userID == value {
										bot.PushMessage(groupID, linebot.NewTextMessage(UserNameSlice[i]+" 玩家名稱已修改成 "+UserName)).Do()
										UserNameSlice[i] = UserName
										UserDuplicate = true
										// return
									}
								}
								if !UserDuplicate {
									bot.PushMessage(groupID, linebot.NewTextMessage("人數已滿 , 請等下一輪")).Do()
								}
								// break
							} else {

								UserDuplicate := false //玩家是否重複
								//判斷玩家ID有無重複 有的話變成修改名稱
								for i, value := range UserIDSlice {
									if UserNameSlice[i] == UserName {
										bot.PushMessage(groupID, linebot.NewTextMessage(UserName+" 玩家名稱已存在 ")).Do()
										return
									}
									if userID == value {
										bot.PushMessage(groupID, linebot.NewTextMessage(UserNameSlice[i]+" 玩家名稱已修改成 "+UserName)).Do()
										UserNameSlice[i] = UserName
										UserDuplicate = true
										// return
									}
								}

								//log.Print("user input" message.Text)
								// res, err := bot.GetProfile(userID).Do()
								// if err != nil {
								// 	log.Print(err)
								// }
								if !UserDuplicate {
									UserNameSlice = append(UserNameSlice, UserName)      //紀錄玩家輸入的名稱 之後推撥會顯示玩家名稱的回合
									UserIDSlice = append(UserIDSlice, userID)            //紀錄玩家的ID 便於後續發送圖片
									UserDiceCount = append(UserDiceCount, diceCount)     //初始都給玩家 6顆骰子
									UserCanSpeakSlice = append(UserCanSpeakSlice, false) //玩家有無回答的權限
									TotalUser := ""
									for _, value := range UserNameSlice {
										TotalUser += value + ","
									}

									bot.PushMessage(groupID, linebot.NewTextMessage(UserName+" 已加入遊戲\n"+"目前玩家有 : "+TotalUser)).Do()
									// log.Print(res.DisplayName)
									log.Print(UserName)
								}
							}
						} else if message.Text[0:6] == "/blow " && isDice && !isGuess {
							isBlow = true
							//先判斷名稱是否重複 是的話幫玩家修改名字 或是 禁止玩家取同樣名稱
							UserName := message.Text[6:len(message.Text)]
							UserDuplicate := false //玩家是否重複
							//判斷玩家ID有無重複 有的話變成修改名稱
							for i, value := range UserIDSlice {
								if UserNameSlice[i] == UserName {
									bot.PushMessage(groupID, linebot.NewTextMessage(UserName+" 玩家名稱已存在 ")).Do()
									return
								}
								if userID == value {
									bot.PushMessage(groupID, linebot.NewTextMessage(UserNameSlice[i]+" 玩家名稱已修改成 "+UserName)).Do()
									UserNameSlice[i] = UserName
									UserDuplicate = true

								}
							}
							if !UserDuplicate {
								UserNameSlice = append(UserNameSlice, UserName)      //紀錄玩家輸入的名稱 之後推撥會顯示玩家名稱的回合
								UserIDSlice = append(UserIDSlice, userID)            //紀錄玩家的ID 便於後續發送圖片
								UserDiceCount = append(UserDiceCount, blowdiceCount) //初始都給玩家 6顆骰子
								UserCanSpeakSlice = append(UserCanSpeakSlice, false) //玩家有無回答的權限
								TotalUser := ""
								for _, value := range UserNameSlice {
									TotalUser += value + ","
								}

								bot.PushMessage(groupID, linebot.NewTextMessage(UserName+" 已加入遊戲\n"+"目前玩家有 : "+TotalUser)).Do()
								// log.Print(res.DisplayName)
								log.Print(UserName)
							}
						}
					}

					if message.Text == "/blowstart" && !isGameStart && !isBlowGameStart {
						isBlowGameStart = true
						log.Print("blowstart Receive")
						m_groupID = groupID
						NeedDiceCount = len(UserIDSlice) + 1 //最少要 (玩家數量+1) 當作起始值
						for _, value := range UserIDSlice {
							rand.Seed(time.Now().UnixNano())
							SliceValue := make([]int, blowdiceCount)
							// arrValue := [...]int{1,2,3,4,5,6}
							NumerString := ""
							for _, element := range SliceValue {
								element = rand.Intn(6)
								element = element + 1
								NumerString = NumerString + strconv.Itoa(element)
								if element == 1 {
									AllDiceValueAndCountNoOne[element] = AllDiceValueAndCountNoOne[element] + 1
									AllDiceValueAndCount[1] = AllDiceValueAndCount[1] + 1
									AllDiceValueAndCount[2] = AllDiceValueAndCount[2] + 1
									AllDiceValueAndCount[3] = AllDiceValueAndCount[3] + 1
									AllDiceValueAndCount[4] = AllDiceValueAndCount[4] + 1
									AllDiceValueAndCount[5] = AllDiceValueAndCount[5] + 1
									AllDiceValueAndCount[6] = AllDiceValueAndCount[6] + 1
								} else if element == 2 {
									AllDiceValueAndCount[element] = AllDiceValueAndCount[element] + 1
									AllDiceValueAndCountNoOne[element] = AllDiceValueAndCountNoOne[element] + 1
								} else if element == 3 {
									AllDiceValueAndCount[element] = AllDiceValueAndCount[element] + 1
									AllDiceValueAndCountNoOne[element] = AllDiceValueAndCountNoOne[element] + 1
								} else if element == 4 {
									AllDiceValueAndCount[element] = AllDiceValueAndCount[element] + 1
									AllDiceValueAndCountNoOne[element] = AllDiceValueAndCountNoOne[element] + 1
								} else if element == 5 {
									AllDiceValueAndCount[element] = AllDiceValueAndCount[element] + 1
									AllDiceValueAndCountNoOne[element] = AllDiceValueAndCountNoOne[element] + 1
								} else if element == 6 {
									AllDiceValueAndCount[element] = AllDiceValueAndCount[element] + 1
									AllDiceValueAndCountNoOne[element] = AllDiceValueAndCountNoOne[element] + 1
								}
							}
							UserAnsMap[value] = NumerString
							log.Print("NumerString = " + NumerString)
							log.Print(UserAnsMap)
							log.Print(AllDiceValueAndCount)
							log.Print(AllDiceValueAndCountNoOne)
							log.Print(value)
							//發送給玩家圖片
							// bot.PushMessage(value, linebot.NewTextMessage("==此局為新的一局牌面==")).Do()
							bot.PushMessage(
								value,
								linebot.NewImageMessage(
									"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
									"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
								),
							).Do()
						}
						rand.Seed(time.Now().UnixNano())
						WhoRound = rand.Intn(len(UserIDSlice))
						//讓第一位玩家 可以回答
						bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage(" 請決定你要喊的骰子點數及數量 格式為 x/x (e.g.: 3/3)")).Do()
						//發給群組 現在是誰的回合
						bot.PushMessage(groupID, linebot.NewTextMessage("現在是 "+UserNameSlice[WhoRound]+"的回合")).Do()
					}

					if message.Text == "/dicestart" && !isGameStart && !isBlowGameStart {
						isGameStart = true
						log.Print("dicestart Receive")
						m_groupID = groupID
						for _, value := range UserIDSlice {
							rand.Seed(time.Now().UnixNano())
							SliceValue := make([]int, diceCount)
							// arrValue := [...]int{1,2,3,4,5,6}
							NumerString := ""
							for _, element := range SliceValue {
								element = rand.Intn(6)
								element = element + 1
								NumerString = NumerString + strconv.Itoa(element)
							}
							UserAnsMap[value] = NumerString
							log.Print("NumerString = " + NumerString)
							log.Print(UserAnsMap)
							log.Print(value)
							//發送給玩家圖片
							bot.PushMessage(value, linebot.NewTextMessage("==此局為新的一局牌面==")).Do()
							bot.PushMessage(
								value,
								linebot.NewImageMessage(
									"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
									"https://jenny-web.herokuapp.com/dice/merge/"+NumerString+"/0/564531635164",
								),
							).Do()
						}
						rand.Seed(time.Now().UnixNano())
						WhoRound = rand.Intn(len(UserIDSlice))
						//讓第一位玩家 可以回答
						bot.PushMessage(UserIDSlice[WhoRound], linebot.NewTextMessage("請決定你要喊的骰子\n 1)單 \n 2)雙 \n 3)大\n 4)小 \n 5)紅 \n 6)黑")).Do()
						//發給群組 現在是誰的回合
						bot.PushMessage(groupID, linebot.NewTextMessage("現在是 "+UserNameSlice[WhoRound]+"的回合")).Do()

					}
				}

				// log.Print(event.ReplyToken)
				log.Print(message.Text)
			}
		}
	}
}
