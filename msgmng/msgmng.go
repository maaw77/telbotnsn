package msgmng

import (
	"log"
	"time"

	"github.com/maaw77/telbotnsn/bot"
)

func MessageManager(mQ chan<- bot.MessageToBot) {
	for {
		select {
		case mQ <- bot.MessageToBot{
			ChatId: 80901973,
			Text:   "Hello! 0",
		}:
		case mQ <- bot.MessageToBot{
			ChatId: 80901973,
			Text:   "Hello! 1",
		}:
		}
		log.Println("eng select")
		log.Println(<-time.After(time.Second * 5))
	}
}
