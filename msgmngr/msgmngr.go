package msgmngr

import (
	"errors"
	"fmt"
	"log"

	"github.com/maaw77/telbotnsn/brds"
	// "github.com/maaw77/telbotnsn/zbx"
)

// CommandFromBot represents an instance of a command from BOT to the message manager..
type CommandFromBot struct {
	User brds.User
	// UserID      int
	// Username    string
	TextCommand string
	TextMessage string
}

// CommandFromZbx represents an instance of a command from ZBX to the message manager.

type CommandFromZbx struct {
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
// func formatHostZbx(svdHosts *brds.SavedHosts) (outHosts string) {
// 	svdHosts.RWD.RLock()
// 	defer svdHosts.RWD.RUnlock()
// 	if len(svdHosts.Hosts) < 1 {
// 		return "There are no problematic hosts!"
// 	}
// 	for _, host := range svdHosts.Hosts {
// 		outHosts += fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v\n", host.NameZ, host.ProblemZ)
// 	}
// 	outHosts += fmt.Sprintf("\n<b>The number of problematic hosts is %d.</b>", len(svdHosts.Hosts))
// 	return
// }

// formatProblemHostZbx returns a list of problematic hosts formatted as a string
func formatProblemHostZbx(prblmHost *brds.SavedHosts) (outHosts string, err error) {

	if prblmHost == nil || prblmHost.Hosts == nil {
		return outHosts, errors.New("input data is nil")
	}

	prblmHost.RWD.RLock()
	defer prblmHost.RWD.RUnlock()

	for _, host := range prblmHost.Hosts {
		marker := ""
		if host.ItChanged {
			marker = "ch_"
		} else if host.ItNew {
			marker = "new_"
		}
		outHosts += fmt.Sprintf("<b>%sHost name:</b> %s, <b>problems:</b>%v\n", marker, host.NameZ, host.ProblemZ)
	}

	outHosts += fmt.Sprintf("\n<b>The number of problematic hosts is %d.</b>", len(prblmHost.Hosts))
	return
}

// ormatRestoredHostZbx returns a list of restored hosts formatted as a string
func formatRestoredHostZbx(rstrdHost *brds.SavedHosts) (outHosts string, err error) {

	if rstrdHost == nil || rstrdHost.Hosts == nil {
		return outHosts, errors.New("input data is nil")
	}

	rstrdHost.RWD.RLock()
	defer rstrdHost.RWD.RUnlock()

	for _, host := range rstrdHost.Hosts {
		outHosts += fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v\n", host.NameZ, host.ProblemZ)
	}

	outHosts += fmt.Sprintf("\n<b>The number of restored hosts is %d.</b>", len(rstrdHost.Hosts))
	return
}

// sendMsgAllUsers sends messages to all registered users of the bot.
func sendMsgAllUsers(text string, mQ chan<- MessageToBot, rgdUsers *brds.RegesteredUsers) {
	rgdUsers.RWD.RLock()
	defer rgdUsers.RWD.RUnlock()

	for _, user := range rgdUsers.Users {
		if user.Id != 0 {
			mQ <- MessageToBot{
				ChatId:    user.Id,
				Text:      text,
				ParseMode: "HTML",
			}
		}
	}

}

// MessageManage controls the sending of messages to the Telegram
func MessageManager(mQ chan<- MessageToBot, fromBot <-chan CommandFromBot, fromZbx <-chan CommandFromZbx,
	regUsers *brds.RegesteredUsers, prblmHosts, rstrdHosts *brds.SavedHosts) {

	client, ctx := brds.InitClient()

	for {
		select {
		case cmd := <-fromBot:
			switch cmd.TextCommand {
			case "/start":

				// mQ <- MessageToBot{
				// 	ChatId:    cmd.UserID,
				// 	Text:      "START",
				// 	ParseMode: "HTML",
				// }

				if err := brds.UpdateRegUsers(client, ctx, regUsers); err != nil {
					log.Println(err)
				} else {
					regUsers.RWD.RLock()
					_, ok := regUsers.Users[cmd.User.Username]
					regUsers.RWD.RUnlock()
					if ok {
						mQ <- MessageToBot{
							ChatId:    cmd.User.Id,
							Text:      "<i>You are authenticated!</i>",
							ParseMode: "HTML",
						}

						regUsers.RWD.Lock()
						regUsers.Users[cmd.User.Username] = cmd.User
						regUsers.RWD.Unlock()

						if err := brds.SaveRegUsers(client, ctx, regUsers); err != nil {
							log.Println(err)
						}

					} else {
						mQ <- MessageToBot{
							ChatId:    cmd.User.Id,
							Text:      "<i>You are not registered!</i>",
							ParseMode: "HTML",
						}

					}
				}
			case "/listp":
				regUsers.RWD.RLock()
				userBot, ok := regUsers.Users[cmd.User.Username]
				regUsers.RWD.RUnlock()
				if !ok || userBot.Id != cmd.User.Id {
					mQ <- MessageToBot{
						ChatId:    cmd.User.Id,
						Text:      "<i>You aren't authenticated!\nUse the '/start' command.</i>",
						ParseMode: "HTML",
					}

				} else {
					outSring, _ := formatProblemHostZbx(prblmHosts)
					// log.Println(outSring)
					mQ <- MessageToBot{
						ChatId:    cmd.User.Id,
						Text:      outSring,
						ParseMode: "HTML",
					}

				}
			case "/listr":
				regUsers.RWD.RLock()
				userBot, ok := regUsers.Users[cmd.User.Username]
				regUsers.RWD.RUnlock()
				if !ok || userBot.Id != cmd.User.Id {
					mQ <- MessageToBot{
						ChatId:    cmd.User.Id,
						Text:      "<i>You aren't authenticated!\nUse the '/start' command.</i>",
						ParseMode: "HTML",
					}

				} else {
					outSring, _ := formatRestoredHostZbx(rstrdHosts)
					mQ <- MessageToBot{
						ChatId:    cmd.User.Id,
						Text:      outSring,
						ParseMode: "HTML",
					}
				}
			case "/help":

				mQ <- MessageToBot{
					ChatId:    cmd.User.Id,
					Text:      "<i>Use the following commands: /help | /start | /listp | /listr.</i>",
					ParseMode: "HTML",
				}
			default:
				mQ <- MessageToBot{
					ChatId:    cmd.User.Id,
					Text:      "<i>Unknow command.\nUse the '/help' command.</i>",
					ParseMode: "HTML",
				}
				// 	mQ <- MessageToBot{
				// 		ChatId:    cmd.UserID,
				// 		Text:      cmd.TextMessage,
				// 		ParseMode: "HTML",
				// 	}
				// case "listp":
				// 	outSring, _ := formatProblemHostZbx(prblmHosts)
				// 	log.Println(outSring)
				// 	mQ <- MessageToBot{
				// 		ChatId:    cmd.UserID,
				// 		Text:      outSring,
				// 		ParseMode: "HTML",
				// 	}
				// case "listr":
				// 	outSring, _ := formatRestoredHostZbx(rstrdHosts)
				// 	mQ <- MessageToBot{
				// 		ChatId:    cmd.UserID,
				// 		Text:      outSring,
				// 		ParseMode: "HTML",
				// 	}
			}
		case cmd := <-fromZbx:
			log.Println(cmd)
			go sendMsgAllUsers(cmd.TextMessage, mQ, regUsers)
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
