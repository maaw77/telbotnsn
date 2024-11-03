package brds

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/maaw77/telbotnsn/zbx"
)

type SavedHosts struct {
	RWD   sync.RWMutex
	Hosts map[string]zbx.ZabbixHost
}

type RegesteredUsers struct {
	RWD   sync.RWMutex
	Users map[string]User // Uesrs[User.Username]User
}

// User represents represents registered users.
type User struct {
	Id           int    `redis:"id"`
	IsBot        bool   `redis:"is_bot"`
	FirstName    string `redis:"first_name"`
	LastName     string `redis:"last_name"`
	Username     string `redis:"username"`
	LanguageCode string `redis:"language_code"`
}

// InitClient initializes the Redis client.
func InitClient() (*redis.Client, context.Context) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client, ctx
}

// RegUsers registers the users to whom messages will be sent.
func RegUsers(client *redis.Client, ctx context.Context, users []string) error {
	if len(users) < 1 {
		return errors.New("the list of users is empty")
	}
	for _, user := range users {
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
func DelUsers(client *redis.Client, ctx context.Context, users []string) (int32, error) {
	var countDelUsers int32
	if len(users) < 1 {
		return countDelUsers, errors.New("the list of users is empty")
	}
	for _, user := range users {

		if res, err := client.Del(ctx, "user:"+user).Result(); err != nil {
			return countDelUsers, err
		} else if res > 0 {
			countDelUsers++
			log.Printf("%s  has been deleted\n", user)
		}

	}

	return countDelUsers, nil
}

// UpdateRegUsers updates the list of registered users (type RegesteredUsers) from the database.
func UpdateRegUsers(client *redis.Client, ctx context.Context, regUsers *RegesteredUsers) error {
	users, err := ListUsers(client, ctx)
	log.Println(users)
	if err != nil {
		return err
	}
	if len(users) < 1 {
		log.Println("len(users) < 1 ")
		regUsers.RWD.Lock()
		regUsers.Users = make(map[string]User)
		regUsers.RWD.Unlock()
		// return nil
	} else {
		log.Println("len(users) > 1 ")
		regUsers.RWD.Lock()

		if len(regUsers.Users) == 0 {
			regUsers.Users = users
		} else {
			for k, v := range regUsers.Users {
				_, ok := users[k]
				if ok && v.Id != 0 {
					users[k] = v
				}
				regUsers.Users = users
			}

		}
		regUsers.RWD.Unlock()
	}

	// if err := SaveRegUsers(client, ctx, regUsers); err != nil {
	// 	return err
	// }

	return nil
}

// SaveUsers saves a list of registered users (type RegesteredUsers) in the database.
func SaveRegUsers(client *redis.Client, ctx context.Context, regUsers *RegesteredUsers) error {
	regUsers.RWD.RLock()
	defer regUsers.RWD.RUnlock()
	for _, v := range regUsers.Users {
		if err := AddUser(client, ctx, v); err != nil {
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
