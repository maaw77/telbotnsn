package brds

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
)

// Options for configuring the connection to redis
var (
	AddrDef     = "db:6379"
	PasswordDef = "" // no password set
	DbDef       = 0  // use default DB
)

// Errors
var ErrEmptyInputData = errors.New("the input data is empty")

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
		return 0, ErrEmptyInputData
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
		return host, ErrEmptyInputData
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

// GetAllHosts returns all hosts from the database.
func GetAllHosts(client *redis.Client, ctx context.Context) (hosts map[string]ZabbixHost, err error) {

	if client == nil || ctx == nil {
		return hosts, ErrEmptyInputData
	}

	hosts = make(map[string]ZabbixHost)
	// var host ZabbixHost
	itr := client.Scan(ctx, 0, "host:*", 0).Iterator()
	for itr.Next(ctx) {
		// log.Println(itr.Val())
		after, _ := strings.CutPrefix(itr.Val(), "host:")
		host, err := GetHost(client, ctx, after)
		if err == nil {
			hosts[host.HostIdZ] = host
		}

	}
	// if err = itr.Err(); err != nil {
	// 	return hosts, err
	// }

	return hosts, itr.Err()
}

// DelHost removes ZabbixHost (with hostID) from the database.
func DelHost(client *redis.Client, ctx context.Context, hostID string) (res int64, err error) {
	if client == nil || ctx == nil || hostID == "" {
		return 0, ErrEmptyInputData
	}

	if res, err := client.Del(ctx, "problems:"+hostID).Result(); err != nil {
		return res, err
	}

	res, err = client.Del(ctx, "host:"+hostID).Result()
	// if err != nil {
	// 	return res, err
	// }

	return res, err
}

// AddMultHosts adds multiple hosts to the database.
func AddMultHosts(client *redis.Client, ctx context.Context, hosts map[string]ZabbixHost) (numHosts int64, err error) {
	if client == nil || ctx == nil || hosts == nil {
		return 0, ErrEmptyInputData
	}

	for _, v := range hosts {
		if res, err := AddHost(client, ctx, v); err != nil {
			return numHosts, err
		} else {
			numHosts += res
		}

	}

	return numHosts, nil
}

// DelAllHosts removes all ZabbixHosts from the database.
func DelAllHosts(client *redis.Client, ctx context.Context) (numHosts int64, err error) {
	if client == nil || ctx == nil {
		return 0, ErrEmptyInputData
	}

	itr := client.Scan(ctx, 0, "host:*", 0).Iterator()
	for itr.Next(ctx) {
		// log.Println(itr.Val())
		after, _ := strings.CutPrefix(itr.Val(), "host:")
		if res, err := DelHost(client, ctx, after); err != nil {
			return numHosts, err
		} else {
			numHosts += res
		}

	}
	// if err = itr.Err(); err != nil {
	// 	return hosts, err
	// }

	return numHosts, err
}

// UpdateZabixHosts saves the saved hosts (type SavedHosts) in the database if they are not empty.
// Otherwise, loads the saved hosts (type SavedHosts) from the database.
func UpdateZabixHosts(client *redis.Client, ctx context.Context, svdHosts *SavedHosts) error {
	if client == nil || ctx == nil || svdHosts == nil {
		return ErrEmptyInputData
	}

	svdHosts.RWD.Lock()
	defer svdHosts.RWD.Unlock()

	if len(svdHosts.Hosts) != 0 {
		if _, err := DelAllHosts(client, ctx); err != nil {
			return err
		}
		if _, err := AddMultHosts(client, ctx, svdHosts.Hosts); err != redis.Nil {
			return err
		}
	} else {
		hosts, err := GetAllHosts(client, ctx)
		if err != nil {
			return err
		}
		svdHosts.Hosts = hosts
	}

	return nil
}

// RegUsers registers the users to whom messages will be sent.
func RegUsers(client *redis.Client, ctx context.Context, users []string) error {
	if len(users) < 1 {
		return ErrEmptyInputData
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
	// if err := itr.Err(); err != nil {
	// 	return users, err
	// }

	return users, itr.Err()
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
