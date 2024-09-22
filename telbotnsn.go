package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/maaw77/telbotnsn/brds"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile) //?????

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := brds.AddUsers(client, ctx, []string{"1", "2", "3"}); err != nil {
		log.Println(err)
	}

	users, err := brds.ListUsers(client, ctx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(users)
	}

	users["user:2"]["id"] = "54555555555"

	if err := brds.SaveUsers(client, ctx, users); err != nil {
		log.Println(err)
	}

	users, err = brds.ListUsers(client, ctx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(users)
	}

	// if err := delUsers(client, ctx, []string{"1", "2", "3"}); err != nil {
	// 	log.Println(err)
	// }

}
