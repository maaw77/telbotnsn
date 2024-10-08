package brds

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// User represents represents registered users.
type User struct {
	Id           int    `redis:"id"`
	IsBot        bool   `redis:"is_bot"`
	FirstName    string `redis:"first_name"`
	LastName     string `redis:"last_name"`
	Username     string `redis:"username"`
	LanguageCode string `redis:"language_code"`
}

// RegUsers registers the users to whom messages will be sent.
func RegUsers(client *redis.Client, ctx context.Context, users []string) error {
	if len(users) < 1 {
		return errors.New("the list of users is empty")
	}
	for i, user := range users {
		fmt.Println("User", i+1, ": ", user)
		if err := AddUser(client, ctx, User{Username: user}); err != nil {
			return err
		}

	}

	return nil

}

// addUsers adds user to whom messages will be sent.
func AddUser(client *redis.Client, ctx context.Context, user User) error {
	if err := client.HSet(ctx, "user:"+user.Username, user).Err(); err != nil {
		return err
	}

	// for _, user := range users {
	// 	if err := client.HSet(ctx, "user:"+user, "id", "0").Err(); err != nil {
	// 		return err
	// 	}
	// 	// log.Println("user:"+user+"=", client.HGetAll(ctx, "user:"+user).Val())
	// }

	return nil
}

// ListUsers returns a list(hash table) of registered users.
func ListUsers(client *redis.Client, ctx context.Context) (map[string]User, error) {
	var usr User
	users := make(map[string]User)

	itr := client.Scan(ctx, 0, "user:*", 0).Iterator()
	for itr.Next(ctx) {

		if err := client.HGetAll(ctx, itr.Val()).Scan(&usr); err != nil {
			return users, err
		}
		users[usr.Username] = usr
	}
	if err := itr.Err(); err != nil {
		return users, err
	}

	return users, nil
}

// DelUsers removes users from the mailing list.
func DelUsers(client *redis.Client, ctx context.Context, users []string) error {
	if len(users) < 1 {
		return errors.New("the list of users is empty")
	}
	for _, user := range users {

		if res, err := client.Del(ctx, "user:"+user).Result(); err != nil {
			return err
		} else {
			log.Println("res=", res)
		}
	}

	return nil
}

// SaveUsers  saves the list of users to the database.
func SaveUsers(client *redis.Client, ctx context.Context, users map[string]map[string]string) error {
	if len(users) < 1 {
		return errors.New("the list of users is empty")
	}
	for key, val := range users {
		if err := client.HSet(ctx, key, val).Err(); err != nil {
			return err
		}
	}

	return nil
}

// func main() {
// 	log.SetFlags(log.Ldate | log.Lshortfile) //?????

// 	if err := godotenv.Load(".env"); err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx := context.Background()
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6380",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})

// 	if err := addUsers(client, ctx, []string{"1", "2", "3"}); err != nil {
// 		log.Println(err)
// 	}

// 	users, err := listUsers(client, ctx)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Println(users)
// 	}

// 	users["user:2"]["id"] = "323434"

// 	if err := saveUsers(client, ctx, users); err != nil {
// 		log.Println(err)
// 	}

// 	users, err = listUsers(client, ctx)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		log.Println(users)
// 	}

// 	// if err := delUsers(client, ctx, []string{"1", "2", "3"}); err != nil {
// 	// 	log.Println(err)
// 	// }

// }
