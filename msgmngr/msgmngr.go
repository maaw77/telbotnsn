package msgmngr

import (
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/zbx"
)

type SavedHosts struct {
	RWD   sync.RWMutex
	Hosts map[string]zbx.ZabbixHost
}

type RegesteredUsers struct {
	RWD   sync.RWMutex
	Users map[string]bot.User
}

func sendMessage(mQ chan<- bot.MessageToBot, oZ zbx.ZabbixHost, rgdUsers *RegesteredUsers) {
	rgdUsers.RWD.RLock()
	for _, user := range rgdUsers.Users {
		text := fmt.Sprintf("<b>Host name:</b> %s, <b>problems:</b>%v", oZ.NameZ, oZ.ProblemZ)
		mQ <- bot.MessageToBot{
			ChatId:    user.Id,
			Text:      text,
			ParseMode: "HTML",
		}
	}
	rgdUsers.RWD.RUnlock()
}

func MessageManager(mQ chan<- bot.MessageToBot, fromZabbix <-chan zbx.ZabbixHost, rgdUsers *RegesteredUsers, svdHosts *SavedHosts) {
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
			if slices.Compare(hst.ProblemZ, oZ.ProblemZ) != 0 {
				sendMessage(mQ, oZ, rgdUsers)
				svdHosts.RWD.Lock()
				svdHosts.Hosts[oZ.HostidZ] = oZ
				svdHosts.RWD.Unlock()
			}
		}
	}
	log.Println(svdHosts.Hosts, len(svdHosts.Hosts))
	// close(mQ)
}
