package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/maaw77/telbotnsn/brds"
	"github.com/maaw77/telbotnsn/msgmngr"
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

// // A MessageToBot is a message sent to the user.
// type MessageToBot struct {
// 	ChatId    int    `json:"chat_id"` // User or chat ID.
// 	Text      string `json:"text"`    // The text of the message.
// 	ParseMode string `json:"parse_mode,omitempty"`
// }

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
func (b *Bot) SendMessage(msg *msgmngr.MessageToBot) error {

	bJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// log.Println(string(bJSON))
	resp, err := http.Post(b.urls.sendMessageURL, "application/json", bytes.NewReader(bJSON))
	if err != nil {
		return err
	}
	tx, _ := io.ReadAll(resp.Body)
	log.Println(string(tx))
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
func sendMessages(bot Bot, mQ <-chan msgmngr.MessageToBot) {
	poolMaxSize := make(chan time.Time, 5)
	for msg := range mQ {
		poolMaxSize <- <-time.After(10 * time.Millisecond)
		go func(b Bot, m msgmngr.MessageToBot) {

			if err := bot.SendMessage(&m); err != nil {
				log.Println(err)
			}
			<-poolMaxSize
		}(bot, msg)

	}
}

func poller(bot Bot, comandToMM chan<- msgmngr.CommandFromBot, regUsers *brds.RegesteredUsers) {
	var lastIDUpdate int

	client, ctx := brds.InitClient()

	for {
		results, err := bot.GetUpdates(ParamGetUpdates{Offset: lastIDUpdate + 1, Timeout: 3, Allowed_updates: []string{"message"}})

		if err != nil {
			log.Println(err)
		} else {
			for _, res := range results.Result {
				if res.UpdateId > lastIDUpdate {
					lastIDUpdate = res.UpdateId
				}

				switch res.Message.Text {
				case "/start":
					comandToMM <- msgmngr.CommandFromBot{
						UserID:      res.Message.From.Id,
						TextMessage: "Command is 'start'",
						TextCommand: "print",
					}
					if err := brds.UpdateRegUsers(client, ctx, regUsers); err != nil {
						log.Println(err)
					} else {
						regUsers.RWD.RLock()
						_, ok := regUsers.Users[res.Message.From.Username]
						regUsers.RWD.RUnlock()
						if ok {
							comandToMM <- msgmngr.CommandFromBot{
								UserID:      res.Message.From.Id,
								TextMessage: "<i>You are authenticated!</i>",
								TextCommand: "print",
							}

							regUsers.RWD.Lock()
							regUsers.Users[res.Message.From.Username] = brds.User{Id: res.Message.From.Id,
								IsBot:        res.Message.From.IsBot,
								FirstName:    res.Message.From.FirstName,
								LastName:     res.Message.From.LastName,
								Username:     res.Message.From.Username,
								LanguageCode: res.Message.From.LanguageCode,
							}
							regUsers.RWD.Unlock()

							if err := brds.SaveRegUsers(client, ctx, regUsers); err != nil {
								log.Println(err)
							}

						} else {
							comandToMM <- msgmngr.CommandFromBot{
								UserID:      res.Message.From.Id,
								TextMessage: "<i>You are not registered!</i>",
								TextCommand: "print",
							}

						}
					}
				case "/listp":
					regUsers.RWD.RLock()
					userBot, ok := regUsers.Users[res.Message.From.Username]
					regUsers.RWD.RUnlock()
					if !ok || userBot.Id != res.Message.From.Id {
						comandToMM <- msgmngr.CommandFromBot{
							UserID:      res.Message.From.Id,
							TextMessage: "<i>You aren't authenticated!\nUse the '/start' command.</i>",
							TextCommand: "print",
						}

					} else {
						comandToMM <- msgmngr.CommandFromBot{
							UserID:      res.Message.From.Id,
							TextCommand: "listp",
						}
					}
				case "/listr":
					regUsers.RWD.RLock()
					userBot, ok := regUsers.Users[res.Message.From.Username]
					regUsers.RWD.RUnlock()
					if !ok || userBot.Id != res.Message.From.Id {
						comandToMM <- msgmngr.CommandFromBot{
							UserID:      res.Message.From.Id,
							TextMessage: "<i>You aren't authenticated!\nUse the '/start' command.</i>",
							TextCommand: "print",
						}

					} else {
						comandToMM <- msgmngr.CommandFromBot{
							UserID:      res.Message.From.Id,
							TextCommand: "listr",
						}
					}
				case "/help":
					comandToMM <- msgmngr.CommandFromBot{
						UserID:      res.Message.From.Id,
						TextMessage: "<i>Use the following commands: /help | /start | /listp | /lisr.</i>",
						TextCommand: "print",
					}

				default:
					comandToMM <- msgmngr.CommandFromBot{
						UserID:      res.Message.From.Id,
						TextMessage: "<i>Unknow command.\nUse the '/help' command.</i>",
						TextCommand: "print",
					}

				}
			}
		}
	}
}

// botRun launches the bot.
func Run(botToken string, mQ chan msgmngr.MessageToBot, comandToMM chan<- msgmngr.CommandFromBot, rgdUsers *brds.RegesteredUsers) {
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
		poller(bot, comandToMM, rgdUsers)
	}()
	waitGroup.Wait()

}
