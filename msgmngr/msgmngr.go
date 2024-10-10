package msgmngr

import (
	"fmt"
	"slices"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/brds"
	"github.com/maaw77/telbotnsn/zbx"
)

// sendMessage sends messages to registered users.
func sendMessage(mQ chan<- bot.MessageToBot, oZ zbx.ZabbixHost, rgdUsers *brds.RegesteredUsers) {
	rgdUsers.RWD.RLock()
	for _, user := range rgdUsers.Users {
		if user.Id != 0 {
			text := fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v", oZ.NameZ, oZ.ProblemZ)
			mQ <- bot.MessageToBot{
				ChatId:    user.Id,
				Text:      text,
				ParseMode: "HTML",
			}
		}
	}
	rgdUsers.RWD.RUnlock()
}

// MessageManage controls the sending of messages from Zabbix to the Telegram bot
func MessageManager(mQ chan<- bot.MessageToBot, fromZabbix <-chan zbx.ZabbixHost, rgdUsers *brds.RegesteredUsers, svdHosts *brds.SavedHosts) {
	// usersId := []int{80901973}
	for oZ := range fromZabbix {
		svdHosts.RWD.RLock()
		hst, ok := svdHosts.Hosts[oZ.HostidZ]
		svdHosts.RWD.RUnlock()

		if !ok && len(oZ.ProblemZ) != 0 {

			sendMessage(mQ, oZ, rgdUsers)
			svdHosts.RWD.Lock()
			svdHosts.Hosts[oZ.HostidZ] = oZ
			svdHosts.RWD.Unlock()
		} else if ok {
			if len(oZ.ProblemZ) == 0 {
				oZ.ProblemZ = append(oZ.ProblemZ, "No problems at all!")
			}
			if slices.Compare(hst.ProblemZ, oZ.ProblemZ) != 0 {

				sendMessage(mQ, oZ, rgdUsers)
				svdHosts.RWD.Lock()
				svdHosts.Hosts[oZ.HostidZ] = oZ
				svdHosts.RWD.Unlock()
			}
		}
	}
	// log.Println(svdHosts.Hosts, len(svdHosts.Hosts))
	// close(mQ)
}
