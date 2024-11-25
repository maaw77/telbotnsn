package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	for textMessage := range sliceMessage(msg.Text, 4000) {
		msg.Text = textMessage
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
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errors.New(resp.Status)
		}
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

// sliceMessage  slices the incoming text according to the limit.
// Sends chunks in turn to the outString channel.
func sliceMessage(incomingText string, limit int) (outString chan string) {
	outString = make(chan string)
	go func() {
		defer close(outString)
		if limit <= 0 {
			return
		}
		// log.Println("incomingText= ", incomingText)
		lenIncText := len(incomingText)
		// log.Println("lenIncText= ", lenIncText)
		if lenIncText <= limit {
			outString <- incomingText
		} else {
			p := make([]byte, limit)
			var chunkIncString strings.Builder
			readerString := strings.NewReader(incomingText)
			var counterBytes int

			for {

				n, err := readerString.Read(p)

				if n != 0 {
					lastIndexNl := strings.LastIndexByte(string(p[:n]), '\n')
					if lastIndexNl > 0 {
						if lastIndexNl+1 < n {
							offset := (lastIndexNl + 1) - n
							readerString.Seek(int64(offset), io.SeekCurrent)
							n = lastIndexNl + 1
						}
						if counterBytes != 0 {
							chunkIncString.WriteString("...")
						}
						counterBytes += n

						chunkIncString.Write(p[:n-1])

						if int64(counterBytes) < readerString.Size() {
							chunkIncString.WriteString("...")
						}
						outString <- chunkIncString.String()
						chunkIncString.Reset()

					}
				}

				if err != nil {
					return
				}
			}

		}

	}()

	return outString
}

// poller receives incoming updates using long polling (https://core.telegram.org/bots/api#getting-updates).
// Sends the received data via comandToMM channel.
func poller(bot Bot, comandToMM chan<- msgmngr.CommandFromBot) {
	var lastIDUpdate int

	// client, ctx := brds.InitClient()

	for {
		results, err := bot.GetUpdates(ParamGetUpdates{Offset: lastIDUpdate + 1, Timeout: 3, Allowed_updates: []string{"message"}})

		if err != nil {
			log.Println(err)
		} else {
			for _, res := range results.Result {
				if res.UpdateId > lastIDUpdate {
					lastIDUpdate = res.UpdateId
				}
				comandToMM <- msgmngr.CommandFromBot{
					User: brds.User{Id: res.Message.From.Id,
						IsBot:        res.Message.From.IsBot,
						FirstName:    res.Message.From.FirstName,
						LastName:     res.Message.From.LastName,
						Username:     res.Message.From.Username,
						LanguageCode: res.Message.From.LanguageCode,
					},
					TextCommand: res.Message.Text,
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
		poller(bot, comandToMM)
	}()
	waitGroup.Wait()

}
