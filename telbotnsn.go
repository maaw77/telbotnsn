package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/brds"
	"github.com/maaw77/telbotnsn/msgmngr"
	"github.com/maaw77/telbotnsn/zbx"
)

func main() {

	// implement problem counting???

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	argumentsCLI := os.Args
	if len(argumentsCLI) < 2 {
		fmt.Println("Usage: run|users <arguments>")
		return
	}
	switch argumentsCLI[1] {
	case "run":
		regUsers := msgmngr.RegesteredUsers{Users: map[string]brds.User{"maaw77": {Username: "maaw77", Id: 80901973}}}
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
	case "users":
		if len(argumentsCLI) < 3 {
			fmt.Println("Usage: users -add|-del <username1> <username2>")
			fmt.Println("Usage: users -list")
			return
		}
		ctx := context.Background()
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6380",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		switch argumentsCLI[2] {
		case "-list":
			users, err := brds.ListUsers(client, ctx)
			if err != nil {
				log.Fatal(err)
			}
			for key, value := range users {
				fmt.Printf("User: %s; data: %#v\n", key, value)
			}

		case "-add":
			if len(argumentsCLI) < 4 {
				fmt.Println("Usage: users -add|-del <username1> <username2>")
				return
			}
			fmt.Println(argumentsCLI[2:])
			if err := brds.RegUsers(client, ctx, argumentsCLI[3:]); err != nil {
				log.Fatal(err)
			}
		case "-del":
			if len(argumentsCLI) < 4 {
				fmt.Println("Usage: users -add|-del <username1> <username2>")
				return
			}
			fmt.Println(argumentsCLI[2:])
		default:
			fmt.Println("Not a valid options")
			return
		}

	default:
		fmt.Println("Not a valid options")
	}

}
