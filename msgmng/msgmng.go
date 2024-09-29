package msgmng

import (
	"fmt"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/zbx"
)

func MessageManager(mQ chan<- bot.MessageToBot, fromZabbix <-chan zbx.ZabbixHost) {

	usersId := []int{80901973}
	for oZ := range fromZabbix {
		for _, userId := range usersId {
			text := fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v", oZ.NameZ, oZ.ProblemZ)
			mQ <- bot.MessageToBot{
				ChatId:    userId,
				Text:      text,
				ParseMode: "HTML",
			}
		}

	}
	close(mQ)
}
