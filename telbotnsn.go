package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"

	"github.com/maaw77/telbotnsn/bot"
	"github.com/maaw77/telbotnsn/brds"
	"github.com/maaw77/telbotnsn/msgmngr"
	"github.com/maaw77/telbotnsn/zbx"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("ZABBIX_WILDCARDSHOSTS") != "" {
		zbx.WILDCARD = os.Getenv("ZABBIX_WILDCARDSHOSTS")
	}

	if os.Getenv("ZABBIX_SLEEP") != "" {
		sd, err := strconv.ParseInt(os.Getenv("ZABBIX_SLEEP"), 10, 32)
		if err == nil {
			zbx.SleepDuration = int(sd)
		}
	}

	argumentsCLI := os.Args

	if len(argumentsCLI) < 2 {
		fmt.Printf("Usage: %s run|users <arguments>\n", os.Args[0])
		return
	}
	switch argumentsCLI[1] {
	case "run":
		var regUsers brds.RegesteredUsers
		client, ctx := brds.InitClient()
		if err := brds.UpdateRegUsers(client, ctx, &regUsers); err != nil {
			log.Fatal(err)
		}
		regUsers.RWD.RLock()
		log.Println(regUsers.Users)
		regUsers.RWD.RUnlock()

		svdHosts := brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{}}
		if err := brds.UpdateZabixHosts(client, ctx, &svdHosts); err != nil {
			log.Fatal(err)
		}
		client.Close()

		// svdHosts.RWD.Lock()
		// svdHosts.Hosts["Host_1"] = brds.ZabbixHost{HostIdZ: "111",
		// 	NameZ: "Host_1"}
		// svdHosts.Hosts["Host_2"] = brds.ZabbixHost{HostIdZ: "222",
		// 	NameZ: "Host_2"}
		// svdHosts.Hosts["Host_3"] = brds.ZabbixHost{HostIdZ: "333",
		// 	NameZ: "Host_3"}
		// svdHosts.RWD.Unlock()

		rstrdHosts := brds.SavedHosts{Hosts: map[string]brds.ZabbixHost{}}

		// outZabbix := make(chan zbx.ZabbixHost)
		messageQueue := make(chan msgmngr.MessageToBot, 5)
		commandQueueFromBot := make(chan msgmngr.CommandFromBot, 5)
		commandQueueFromZbx := make(chan msgmngr.CommandFromZbx, 5)

		var waitGroup sync.WaitGroup

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			zbx.Run(os.Getenv("ZABBIX_USERNAME"), os.Getenv("ZABBIX_PASSWORD"), commandQueueFromZbx, &svdHosts, &rstrdHosts)

		}()

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			msgmngr.MessageManager(messageQueue, commandQueueFromBot, commandQueueFromZbx, &regUsers, &svdHosts, &rstrdHosts)

		}()

		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			bot.Run(os.Getenv("BOT_TOKEN"), messageQueue, commandQueueFromBot, &regUsers)
		}()

		waitGroup.Wait()
	case "users":
		if len(argumentsCLI) < 3 {
			fmt.Printf("Usage: %s users -add|-del <username1> <username2>\n", os.Args[0])
			fmt.Printf("Usage: %s users -list\n", os.Args[0])
			return
		}

		client, ctx := brds.InitClient()
		switch argumentsCLI[2] {
		case "-list":
			users, err := brds.ListUsers(client, ctx)
			if err != nil {
				log.Fatal(err)
			}
			if len(users) == 0 {
				fmt.Println("There are no registered users here!")
				return
			}
			for key, value := range users {
				fmt.Printf("User: %s; data: %#v\n", key, value)
			}

		case "-add":
			if len(argumentsCLI) < 4 {
				fmt.Printf("Usage: %s users -add|-del <username1> <username2> ...\n", os.Args[0])
				return
			}
			// fmt.Println(argumentsCLI[2:])
			if err := brds.RegUsers(client, ctx, argumentsCLI[3:]); err != nil {
				log.Fatal(err)
			}
			fmt.Println(argumentsCLI[3:], "has been registered.")
		case "-del":
			if len(argumentsCLI) < 4 {
				fmt.Printf("Usage: %s users -add|-del <username1> <username2> ...\n", os.Args[0])
				return
			}
			countDelUsers, err := brds.DelUsers(client, ctx, argumentsCLI[3:])
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("%d user(s) deleted.\n", countDelUsers)
			}
		default:
			fmt.Println("Not a valid options")
			return
		}

	default:
		fmt.Println("Not a valid options")
	}

}
