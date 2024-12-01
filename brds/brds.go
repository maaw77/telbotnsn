package brds

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var AddrDef = "localhost:6380"
var PasswordDef = "" // no password set
var DbDef = 0        // use default DB

type ZabbixHost struct {
	HostIdZ   string   `redis:"hostid"`
	HostZ     string   `redis:"host"`
	NameZ     string   `redis:"name"`
	StatusZ   string   `redis:"status"`
	ProblemZ  []string `redis:"problem"`
	ItNew     bool     `redis:"new"`
	ItChanged bool     `redis:"changed"`
}

type SavedHosts struct {
	RWD   sync.RWMutex
	Hosts map[string]ZabbixHost //Hosts[ZabbixHost.HostidZ]ZabbixHost

}

type RegesteredUsers struct {
	RWD   sync.RWMutex
	Users map[string]User // Uesrs[User.Username]User
}

// User represents registered users.
type User struct {
	Id           int    `redis:"id"`
	IsBot        bool   `redis:"is_bot"`
	FirstName    string `redis:"first_name"`
	LastName     string `redis:"last_name"`
	Username     string `redis:"username"`
	LanguageCode string `redis:"language_code"`
}

// InitClient initializes the Redis client.
func InitClient() (client *redis.Client, ctx context.Context) {
	ctx = context.Background()
	client = redis.NewClient(&redis.Options{
		Addr:     AddrDef,
		Password: PasswordDef,
		DB:       DbDef,
	})
	return
}

// AddHost adds the host to the database.
func AddHost(client *redis.Client, ctx context.Context, host ZabbixHost) (int64, error) {
	if client == nil || ctx == nil || host.HostIdZ == "" || host.ProblemZ == nil {
		return 0, errors.New("the input data is empty")
	}

	for _, v := range host.ProblemZ {
		res, err := client.RPush(ctx, "problems:"+host.HostIdZ, v).Result()
		if err != nil {
			return res, err
		}
	}

	res, err := client.HSet(ctx, "host:"+host.HostIdZ,
		"hostid", host.HostIdZ,
		"host", host.HostZ,
		"name", host.NameZ,
		"status", host.StatusZ,
		"new", host.ItNew,
		"changed", host.ItChanged,
		"problems", "problems:"+host.HostIdZ).Result()
	if err != nil {
		return res, err
	}
	return res, nil

}

// GetHost returns ZabbixHost from the database.
func GetHost(client *redis.Client, ctx context.Context, hostID string) (host ZabbixHost, err error) {
	if client == nil || ctx == nil || hostID == "" {
		return host, errors.New("the input data is empty")
	}

	res1, err := client.HGetAll(ctx, "host:"+hostID).Result()
	switch {
	case err != nil:
		return host, err
	case len(res1) == 0:
		return host, errors.New("the data was not found")

	}

	// log.Println(res1)
	host.HostIdZ = res1["hostid"]
	host.HostZ = res1["host"]
	host.NameZ = res1["name"]
	host.StatusZ = res1["status"]
	switch {
	case res1["new"] == "1":
		host.ItNew = true
	case res1["changed"] == "1":
		host.ItChanged = true
	}

	res2, err := client.LRange(ctx, res1["problems"], 0, -1).Result()
	if err != nil {
		return host, err
	}
	host.ProblemZ = res2

	return host, nil
}

// DelHost removes ZabbixHost from the database.
func DelHost(client *redis.Client, ctx context.Context, hostID string) (res int64, err error) {
	if client == nil || ctx == nil || hostID == "" {
		return 0, errors.New("the input data is empty")
	}
	if res, err := client.Del(ctx, "problems:"+hostID).Result(); err != nil {
		return res, err
	}

	res, err = client.Del(ctx, "host:"+hostID).Result()
	if err != nil {
		return res, err
	}

	return res, nil
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

// AddUsers adds the user to the database
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
			// log.Printf("%s  has been deleted\n", user)
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
		// log.Println("len(users) < 1 ")
		regUsers.RWD.Lock()
		regUsers.Users = make(map[string]User)
		regUsers.RWD.Unlock()
		// return nil
	} else {
		// log.Println("len(users) > 1 ")
		regUsers.RWD.Lock()

		if len(regUsers.Users) == 0 {
			regUsers.Users = users
		} else {
			for k, v := range regUsers.Users {
				_, ok := users[k]
				if ok && v.Id != 0 {
					users[k] = v
				}
			}
			regUsers.Users = users
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
