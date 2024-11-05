package msgmngr

import (
	"fmt"
	"log"
	"time"

	"github.com/maaw77/telbotnsn/brds"
	// "github.com/maaw77/telbotnsn/zbx"
)

// CommandFromBot is designed to store the data of commands transmitted by the bot to the message manager.
type CommandFromBot struct {
	UserID      int
	TextCommand string
	TextMessage string
}

// A MessageToBot is a message sent to the user.
type MessageToBot struct {
	ChatId    int    `json:"chat_id"` // User or chat ID.
	Text      string `json:"text"`    // The text of the message.
	ParseMode string `json:"parse_mode,omitempty"`
}

// formatHostZbx returns a list of hosts formatted as a string
func formatHostZbx(svdHosts *brds.SavedHosts) (outHosts string) {
	svdHosts.RWD.RLock()
	defer svdHosts.RWD.RUnlock()
	if len(svdHosts.Hosts) < 1 {
		return "There are no problematic hosts!"
	}
	for _, host := range svdHosts.Hosts {
		outHosts += fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v\n", host.NameZ, host.ProblemZ)
	}
	return
}

// MessageManage controls the sending of messages to the Telegram
func MessageManager(mQ chan<- MessageToBot, fromBot <-chan CommandFromBot, rgdUsers *brds.RegesteredUsers, svdHosts *brds.SavedHosts) {
	for {
		select {
		case cmd := <-fromBot:
			switch cmd.TextCommand {
			case "print":

				mQ <- MessageToBot{
					ChatId:    cmd.UserID,
					Text:      cmd.TextMessage,
					ParseMode: "HTML",
				}
			case "list":
				mQ <- MessageToBot{
					ChatId:    cmd.UserID,
					Text:      formatHostZbx(svdHosts),
					ParseMode: "HTML",
				}
			}
		case <-time.After(5 * time.Second):
			log.Println("Default select!!")
		}
	}
}

// // sendMessage sends messages to registered users.
// func sendMessage(mQ chan<- bot.MessageToBot, oZ zbx.ZabbixHost, rgdUsers *brds.RegesteredUsers) {
// 	rgdUsers.RWD.RLock()
// 	for _, user := range rgdUsers.Users {
// 		if user.Id != 0 {
// 			text := fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v", oZ.NameZ, oZ.ProblemZ)
// 			mQ <- bot.MessageToBot{
// 				ChatId:    user.Id,
// 				Text:      text,
// 				ParseMode: "HTML",
// 			}
// 		}
// 	}
// 	rgdUsers.RWD.RUnlock()
// }

// // MessageManage controls the sending of messages from Zabbix to the Telegram bot
// func MessageManager(mQ chan<- bot.MessageToBot, fromZabbix <-chan zbx.ZabbixHost, rgdUsers *brds.RegesteredUsers, svdHosts *brds.SavedHosts) {
// 	// usersId := []int{80901973}
// 	for oZ := range fromZabbix {
// 		svdHosts.RWD.RLock()
// 		hst, ok := svdHosts.Hosts[oZ.HostidZ]
// 		svdHosts.RWD.RUnlock()

// 		if !ok && len(oZ.ProblemZ) != 0 {

// 			sendMessage(mQ, oZ, rgdUsers)
// 			svdHosts.RWD.Lock()
// 			svdHosts.Hosts[oZ.HostidZ] = oZ
// 			svdHosts.RWD.Unlock()
// 		} else if ok {
// 			if len(oZ.ProblemZ) == 0 {
// 				oZ.ProblemZ = append(oZ.ProblemZ, "No problems at all!")
// 			}
// 			if slices.Compare(hst.ProblemZ, oZ.ProblemZ) != 0 {

// 				sendMessage(mQ, oZ, rgdUsers)
// 				svdHosts.RWD.Lock()
// 				svdHosts.Hosts[oZ.HostidZ] = oZ
// 				svdHosts.RWD.Unlock()
// 			}
// 		}
// 	}
// 	// log.Println(svdHosts.Hosts, len(svdHosts.Hosts))
// 	// close(mQ)
// }
