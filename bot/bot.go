package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

const botUrlAPI string = "https://api.telegram.org/"

// An User represents a Telegram user or bot.
type User struct {
	Id            int
	Is_bot        bool
	First_name    string
	Last_name     string
	Username      string
	Language_code string
}

// An ObjectMessage represents a message.
type ObjectMessage struct {
	Message_id int
	From       User
	Date       int
	Text       string
}

// A ResultIncomingUpdate is a field of incoming updates.
type ResultIncomingUpdate struct {
	Update_id int
	Message   ObjectMessage
}

// An IncomingUpdate represents an incoming updates.
type IncomingUpdate struct {
	Ok     bool
	Result []ResultIncomingUpdate
}

// A MessageToBot is a message sent to the user.
type MessageToBot struct {
	ChatId int    `json:"chat_id"` // User or chat ID.
	Text   string `json:"text"`    // The text of the message.
}

// A ParamGetUpdates presents the parameters of the getUpdates method.
type ParamGetUpdates struct {
	Offset          int      `json:"offset"`
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
	client   *http.Client
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
	resp, err := b.client.Get(b.urls.getMeUrl)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("authorization bot error")
	}

	defer resp.Body.Close()
	return nil
}

// SendMessage sends text messages.
func (b *Bot) SendMessage(msg *MessageToBot) error {

	bJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := b.client.Post(b.urls.sendMessageURL, "application/json", bytes.NewReader(bJSON))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()
	return nil
}

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

// botRun launches the bot.
func Run(botToken string, users map[string]int) {
	bot := Bot{
		client:   &http.Client{},
		botToken: botToken,
	}
	bot.SetURLs(botUrlAPI)

	if err := bot.CheckAuth(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("your bot has been authenticated.")
	}

	msg := MessageToBot{80901973, "Hello!"}
	if err := bot.SendMessage(&msg); err != nil {
		log.Println(err)
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

}
