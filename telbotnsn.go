package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/msgmng"
	"github.com/maaw77/telbotnsn/zbx"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

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
		msgmng.MessageManager(messageQueue, outZabbix)

	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		bot.Run(os.Getenv("BOT_TOKEN"), messageQueue)
	}()

	waitGroup.Wait()
	// close(messageQueue)

}
