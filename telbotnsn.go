package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/msgmngr"
	"github.com/maaw77/telbotnsn/zbx"
)

// home
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	regUsers := msgmngr.RegesteredUsers{Users: map[string]bot.User{"maaw77": {Username: "maaw77", Id: 80901973}}}
	svdHosts := msgmngr.SavedHosts{Hosts: map[string]zbx.ZabbixHost{}}
	outZabbix := make(chan zbx.ZabbixHost)
	messageQueue := make(chan bot.MessageToBot, 10)
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		zbx.Run(os.Getenv("ZABBIX_USERNAME"), os.Getenv("ZABBIX_PASSWORD"), outZabbix)

	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		msgmngr.MessageManager(messageQueue, outZabbix, &regUsers, &svdHosts)

	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		bot.Run(os.Getenv("BOT_TOKEN"), messageQueue)
	}()

	waitGroup.Wait()

}
