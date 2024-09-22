package brds

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// addUsers adds users to whom messages will be sent.
func AddUsers(client *redis.Client, ctx context.Context, users []string) error {
	if len(users) > 0 {
		for _, user := range users {
			if err := client.HSet(ctx, "user:"+user, "id", "0").Err(); err != nil {
				return err
			}
			// log.Println("user:"+user+"=", client.HGetAll(ctx, "user:"+user).Val())
		}
	}
	return nil
}

// DelUsers removes users from the mailing list.
func DelUsers(client *redis.Client, ctx context.Context, users []string) error {
	if len(users) > 0 {
		for _, user := range users {

			if res, err := client.Del(ctx, "user:"+user).Result(); err != nil {
				return err
			} else {
				log.Println("res=", res)
			}
		}
	}
	return nil
}

// SaveUsers  saves the list of users to the database.
func SaveUsers(client *redis.Client, ctx context.Context, users map[string]map[string]string) error {
	if len(users) > 0 {
		for key, val := range users {
			if err := client.HSet(ctx, key, val).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// ListUsers returns a list(hash table) of registered users.
func ListUsers(client *redis.Client, ctx context.Context) (map[string]map[string]string, error) {
	users := make(map[string]map[string]string)

	itr := client.Scan(ctx, 0, "user*", 0).Iterator()
	for itr.Next(ctx) {
		users[itr.Val()] = client.HGetAll(ctx, itr.Val()).Val()
	}
	if err := itr.Err(); err != nil {
		return users, err
	}

	return users, nil
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
