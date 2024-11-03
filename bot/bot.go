package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/maaw77/telbotnsn/brds"
)

const botUrlAPI string = "https://api.telegram.org/"

// An User represents a Telegram user or bot.
type User struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"Language_code"`
}

// An ObjectMessage represents a message.
type ObjectMessage struct {
	MessageId int    `json:"message_id "`
	From      User   `json:"from"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

// A ResultIncomingUpdate is a field of incoming updates.
type ResultIncomingUpdate struct {
	UpdateId int           `json:"update_id"`
	Message  ObjectMessage `json:"message"`
}

// An IncomingUpdate represents an incoming updates.
type IncomingUpdate struct {
	Ok     bool                   `json:"ok"`
	Result []ResultIncomingUpdate `json:"result"`
}

// A MessageToBot is a message sent to the user.
type MessageToBot struct {
	ChatId    int    `json:"chat_id"` // User or chat ID.
	Text      string `json:"text"`    // The text of the message.
	ParseMode string `json:"parse_mode,omitempty"`
}

// A ParamGetUpdates presents the parameters of the getUpdates method.
type ParamGetUpdates struct {
	Offset          int      `json:"offset"`
	Timeout         int      `json:"timeout"`
	Allowed_updates []string `json:"allowed_updates"`
}

// A botUrls holds the URLs value for bot methods.
type botUrls struct {
	baseUrl        string
	getMeUrl       string
	sendMessageURL string
	getUpdatesURL  string
}

// A Bot is a client of a telegram bot.
type Bot struct {
	botToken string
	urls     botUrls
	// lastUpdate ResultIncomingUpdate
}

// SetURLs sets the URLs value for bot methods.
func (b *Bot) SetURLs(urlAPI string) {
	var err error
	b.urls.baseUrl, err = url.JoinPath(urlAPI, "bot"+b.botToken)
	if err != nil {
		log.Fatal(err)
	}
	b.urls.getMeUrl, err = url.JoinPath(b.urls.baseUrl, "getMe")
	if err != nil {
		log.Fatal(err)
	}
	b.urls.sendMessageURL, err = url.JoinPath(b.urls.baseUrl, "sendMessage")
	if err != nil {
		log.Fatal(err)
	}
	b.urls.getUpdatesURL, err = url.JoinPath(b.urls.baseUrl, "getUpdates")
	if err != nil {
		log.Fatal(err)
	}
}

// CheckAuth checks your bot's authentication token.
func (b *Bot) CheckAuth() error {
	resp, err := http.Get(b.urls.getMeUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("authorization bot error")
	}
	return nil
}

// SendMessage sends text message.
func (b *Bot) SendMessage(msg *MessageToBot) error {

	bJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// log.Println(string(bJSON))
	resp, err := http.Post(b.urls.sendMessageURL, "application/json", bytes.NewReader(bJSON))
	if err != nil {
		return err
	}
	// tx, _ := io.ReadAll(resp.Body)
	// log.Println(string(tx))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return nil
}

// getUpdates returns an incoming update from the bot
func (b *Bot) GetUpdates(prmtrs ParamGetUpdates) (IncomingUpdate, error) {
	var result IncomingUpdate
	bJSON, err := json.Marshal(prmtrs)
	if err != nil {
		return result, err
	}
	resp, err := http.Post(b.urls.getUpdatesURL, "application/json", bytes.NewReader(bJSON))
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

// sengdMessages sends messages from the channel to users.
func sendMessages(bot Bot, mQ <-chan MessageToBot) {
	poolMaxSize := make(chan time.Time, 5)
	for msg := range mQ {
		poolMaxSize <- <-time.After(10 * time.Millisecond)
		go func(b Bot, m MessageToBot) {

			if err := bot.SendMessage(&m); err != nil {
				log.Println(err)
			}
			<-poolMaxSize
		}(bot, msg)

	}
}

func poller(bot Bot, mQ chan<- MessageToBot, regUsers *brds.RegesteredUsers) {
	var lastIDUpdate int

	client, ctx := brds.InitClient()

	for {
		results, err := bot.GetUpdates(ParamGetUpdates{Offset: lastIDUpdate + 1, Timeout: 3, Allowed_updates: []string{"message"}})
		if err != nil {
			log.Println(err)
		} else if results.Ok {
			for _, res := range results.Result {
				if res.UpdateId > lastIDUpdate {
					lastIDUpdate = res.UpdateId
				}
				log.Printf("%#v", res)
				switch res.Message.Text {
				case "/start":
					log.Println("Command start")

					if err := brds.UpdateRegUsers(client, ctx, regUsers); err != nil {
						log.Println(err)
					}

					regUsers.RWD.Lock()
					log.Println(regUsers.Users)
					temp, ok := regUsers.Users[res.Message.From.Username]
					log.Println("Temp", temp, res.Message.From.Username)
					if ok {
						regUsers.Users[res.Message.From.Username] = brds.User{
							Id:           res.Message.From.Id,
							IsBot:        res.Message.From.IsBot,
							FirstName:    res.Message.From.FirstName,
							LastName:     res.Message.From.LastName,
							Username:     res.Message.From.Username,
							LanguageCode: res.Message.From.LanguageCode,
						}
						brds.SaveRegUsers(client, ctx, regUsers)
					} else {
						mQ <- MessageToBot{
							ChatId: res.Message.From.Id,
							// Text:      res.Message.Text,
							Text:      "<i>You are not registered.</i>",
							ParseMode: "HTML",
						}
					}
					regUsers.RWD.Unlock()
				default:
					mQ <- MessageToBot{
						ChatId: res.Message.From.Id,
						// Text:      res.Message.Text,
						Text:      "<i>Unknow command.</i>",
						ParseMode: "HTML",
					}
				}

			}
		}
	}

}

// botRun launches the bot.
func Run(botToken string, mQ chan MessageToBot, rgdUsers *brds.RegesteredUsers) {
	bot := Bot{botToken: botToken}
	bot.SetURLs(botUrlAPI)

	if err := bot.CheckAuth(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("your bot has been authenticated.")
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		sendMessages(bot, mQ)
	}()
	// ParamGetUpdates{Offset: -1, Allowed_updates: []string{"message"}}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		poller(bot, mQ, rgdUsers)
	}()
	waitGroup.Wait()

}

// var urls botUrls
// urls.setURls(botToken, botUrlAPI)

// client := &http.Client{}

// if err := botCheckAuth(client, &urls); err != nil {
// 	log.Fatal(err)
// }

// // var lastUpdate ResultIncomingUpdate
// upd := botGetUpdates(client, &urls, &ParamGetUpdates{Offset: -1, Allowed_updates: []string{"message"}})
// var lastUpdate ResultIncomingUpdate
// if len(upd.Result) != 0 {
// 	lastUpdate = upd.Result[len(upd.Result)-1]
// }
// for {
// 	lastUpdate = Poller(client, &urls, &ParamGetUpdates{Offset: -1, Allowed_updates: []string{"message"}}, users, &lastUpdate)
// 	log.Println(users)

// 	// for _, val := range users {
// 	// 	msg := MessageToBot{
// 	// 		ChatId: val,
// 	// 		Text:   "Hello!",
// 	// 	}
// 	// 	if err := botSendMessage(client, &urls, &msg); err != nil {
// 	// 		log.Println(err)
// 	// 	}
// 	// }
// }
// func botGetUpdates(client *http.Client, urls *botUrls, param *ParamGetUpdates) *IncomingUpdate {
// 	b, err := json.Marshal(param)
// 	if err != nil {
// 		// return err
// 	}
// 	// fmt.Println(string(b))
// 	resp, err := client.Post(urls.getUpdatesURL, "application/json", bytes.NewReader(b))
// 	if err != nil {
// 		// return err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		// return errors.New(resp.Status)
// 	}

// 	defer resp.Body.Close()

// 	var upd IncomingUpdate
// 	dec := json.NewDecoder(resp.Body)
// 	if err := dec.Decode(&upd); err != nil {
// 		log.Println(err)

// 	}

// 	return &upd

// }

// func Poller(client *http.Client, urls *botUrls, param *ParamGetUpdates, users map[string]int, lastUpdate *ResultIncomingUpdate) ResultIncomingUpdate {
// 	upd := botGetUpdates(client, urls, param)
// 	// log.Printf("%s\n", pretty.Sprint(&upd))

// 	indexResult := len(upd.Result) - 1
// 	if indexResult > 0 {

// 		if lastUpdate.Update_id != upd.Result[indexResult].Update_id {
// 			userName := upd.Result[indexResult].Message.From.Username
// 			cmd := upd.Result[indexResult].Message.Text

// 			if _, ok := users[userName]; cmd == "/start" && ok {
// 				users[userName] = upd.Result[indexResult].Message.From.Id

// 				msg := MessageToBot{
// 					ChatId: users[userName],
// 					Text:   "You have been authenticated!\n Welcome!",
// 				}
// 				if err := botSendMessage(client, urls, &msg); err != nil {
// 					log.Println(err)
// 				}

// 			} else if cmd == "/start" {
// 				msg := MessageToBot{
// 					ChatId: upd.Result[indexResult].Message.From.Id,
// 					Text:   "You are not authenticated!",
// 				}
// 				if err := botSendMessage(client, urls, &msg); err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		}
// 		return upd.Result[indexResult]
// 	}
// 	return ResultIncomingUpdate{}
// }
