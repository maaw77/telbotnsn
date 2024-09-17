package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const botUrlAPI string = "https://api.telegram.org/"

type botUrls struct {
	baseUrl        string
	getMeUrl       string
	sendMessageURL string
	getUpdatesURL  string
}

func (bu *botUrls) setURls(botToken string, urlAPI string) {
	var err error
	bu.baseUrl, err = url.JoinPath(urlAPI, "bot"+botToken)
	if err != nil {
		log.Fatal(err)
	}
	bu.getMeUrl, err = url.JoinPath(bu.baseUrl, "getMe")
	if err != nil {
		log.Fatal(err)
	}
	bu.sendMessageURL, err = url.JoinPath(bu.baseUrl, "sendMessage")
	if err != nil {
		log.Fatal(err)
	}
	bu.getUpdatesURL, err = url.JoinPath(bu.baseUrl, "getUpdates")
	if err != nil {
		log.Fatal(err)
	}
}

// This object represents a Telegram user or bot.
type User struct {
	Id            int
	Is_bot        bool
	First_name    string
	Last_name     string
	Username      string
	Language_code string
}

// This object represents a message.
type ObjectMessage struct {
	Message_id int
	From       User
	Date       int
	Text       string
}

// This object  is a field of incoming updates.
type ResultIncomingUpdate struct {
	Update_id int
	Message   ObjectMessage
}

// This object represents incoming updates.
type IncomingUpdate struct {
	Ok     bool
	Result []ResultIncomingUpdate
}

// This object is a message sent to the user.
type MessageToBot struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

// This object presents the parameters of the getUpdates method.
type ParamGetUpdates struct {
	Offset          int      `json:"offset"`
	Allowed_updates []string `json:"allowed_updates"`
}

// botTestAuth checks your bot's authentication token.
func botCheckAuth(client *http.Client, urls *botUrls) error {
	resp, err := client.Get(urls.getMeUrl)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("authorization bot error")
	}

	defer resp.Body.Close()
	return nil
}

func botSendMessage(client *http.Client, urls *botUrls, msg *MessageToBot) error {

	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := client.Post(urls.sendMessageURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()
	return nil
}

func botGetUpdates(client *http.Client, urls *botUrls, param *ParamGetUpdates) *IncomingUpdate {
	b, err := json.Marshal(param)
	if err != nil {
		// return err
	}
	// fmt.Println(string(b))
	resp, err := client.Post(urls.getUpdatesURL, "application/json", bytes.NewReader(b))
	if err != nil {
		// return err
	}

	if resp.StatusCode != http.StatusOK {
		// return errors.New(resp.Status)
	}

	defer resp.Body.Close()

	var upd IncomingUpdate
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&upd); err != nil {
		log.Println(err)

	}
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	// return err
	// }

	// var upd IncomingUpdate
	// err = json.Unmarshal(body, &upd)
	// if err != nil {
	// 	log.Println(err)
	// }

	return &upd

}

func botPoller(client *http.Client, urls *botUrls, param *ParamGetUpdates, users map[string]int, lastUpdate *ResultIncomingUpdate) ResultIncomingUpdate {
	upd := botGetUpdates(client, urls, param)
	// log.Printf("%s\n", pretty.Sprint(&upd))

	indexResult := len(upd.Result) - 1
	if indexResult >= 0 {

		if lastUpdate.Update_id != upd.Result[indexResult].Update_id {
			userName := upd.Result[indexResult].Message.From.Username
			cmd := upd.Result[indexResult].Message.Text

			if _, ok := users[userName]; cmd == "/start" && ok {
				users[userName] = upd.Result[len(upd.Result)-1].Message.From.Id

				msg := MessageToBot{
					ChatId: users[userName],
					Text:   "You have been authenticated!\n Welcome!",
				}
				if err := botSendMessage(client, urls, &msg); err != nil {
					log.Println(err)
				}

			} else if cmd == "/start" {
				msg := MessageToBot{
					ChatId: upd.Result[len(upd.Result)-1].Message.From.Id,
					Text:   "You are not authenticated!",
				}
				if err := botSendMessage(client, urls, &msg); err != nil {
					log.Println(err)
				}
			}
		}
		return upd.Result[len(upd.Result)-1]
	}
	return ResultIncomingUpdate{}
}

// botRun launches the bot.
func botRun(botToken string, users map[string]int) {
	var urls botUrls
	urls.setURls(botToken, botUrlAPI)

	client := &http.Client{}

	if err := botCheckAuth(client, &urls); err != nil {
		log.Fatal(err)
	}

	// var lastUpdate ResultIncomingUpdate
	upd := botGetUpdates(client, &urls, &ParamGetUpdates{Offset: -1, Allowed_updates: []string{"message"}})
	var lastUpdate ResultIncomingUpdate
	if len(upd.Result) != 0 {
		lastUpdate = upd.Result[len(upd.Result)-1]
	}
	for {
		lastUpdate = botPoller(client, &urls, &ParamGetUpdates{Offset: -1, Allowed_updates: []string{"message"}}, users, &lastUpdate)
		log.Println(users)

		// for _, val := range users {
		// 	msg := MessageToBot{
		// 		ChatId: val,
		// 		Text:   "Hello!",
		// 	}
		// 	if err := botSendMessage(client, &urls, &msg); err != nil {
		// 		log.Println(err)
		// 	}
		// }
	}

}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile) //?????

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	users := map[string]int{"maaw_77": 0}
	botRun(os.Getenv("BOT_TOKEN"), users)

}
