package zbx

import (
	"bytes"
	"encoding/json"
	"errors"
	"html"
	"log"
	"net/http"
	"time"

	"github.com/maaw77/telbotnsn/brds"
)

const zabbixUrlAPI = "http://zabbix.gmkzoloto.ru/zabbix/api_jsonrpc.php"

// "*Микро*"

var WILDCARD = "*Березо*"

// type ZabbixHost struct {
// 	HostidZ  string
// 	HostZ    string
// 	NameZ    string
// 	StatusZ  string
// 	ProblemZ []string
// }

type ZabbixClient struct {
	Username string
	Password string
	Result   string
	Id       int
	URL      string
}

// For information about the Zabbix API, see the link https://www.zabbix.com/documentation/current/en/manual/api.
type ZabbixParams struct {
	Output                 []string          `json:"output,omitempty"`
	Search                 map[string]string `json:"search,omitempty"`
	SearchWildcardsEnabled bool              `json:"searchWildcardsEnabled,omitempty"`
	SearchByAny            bool              `json:"searchByAny,omitempty"`
	Username               string            `json:"username,omitempty"`
	Password               string            `json:"password,omitempty"`
	Host                   string            `json:"host,omitempty"`
	Monitored_hosts        int               `json:"monitored_hosts,omitempty"`
}
type ZabbixRequest struct {
	Jsonrpc string       `json:"jsonrpc,omitempty"`
	Method  string       `json:"method,omitempty"`
	Params  ZabbixParams `json:"params,omitempty"`
	Auth    string       `json:"auth,omitempty"`
	Id      int          `json:"id,omitempty"`
	Result  string       `json:"result,omitempty"`
}

//	type ZabbixParamsResponseHos {
//		Output   []string `json:"output,omitempty"`
//	}
type ZabbixResponse struct {
	Jsonrpc string              `json:"jsonrpc,omitempty"`
	Result  []map[string]string `json:"result,omitempty"`
	// Hosts   []map[string]string `json:"host,omitempty"`
	Id int `json:"id,omitempty"`
}

// Authentication authenticates the Zabbix API client.
func (c *ZabbixClient) Authentication() error {
	payload := ZabbixRequest{Jsonrpc: "2.0",
		Method: "user.login",
		Params: ZabbixParams{Username: c.Username, Password: c.Password},
		Id:     1}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(c.URL, "application/json-rpc", bytes.NewReader(b))
	switch {
	case err != nil:
		return err
	case resp.StatusCode != http.StatusOK:
		return errors.New(resp.Status)
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var zr ZabbixRequest

	if err := dec.Decode(&zr); err != nil {
		return err
	}
	if zr.Result == "" {
		return errors.New("authorization error")
	}
	c.Id = zr.Id
	c.Result = zr.Result
	return nil
}

// GetHost retrieves hosts according to the specified hostname pattern.
func (c *ZabbixClient) GetHost(pattern string) ([]map[string]string, error) {
	payload := ZabbixRequest{Jsonrpc: "2.0",
		Method: "host.get",
		Params: ZabbixParams{
			Output:                 []string{"hostid", "host", "name", "status"},
			Search:                 map[string]string{"name": pattern},
			SearchWildcardsEnabled: true,
			SearchByAny:            true,
			Monitored_hosts:        1,
		},
		Auth: c.Result,
		Id:   c.Id}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.URL, "application/json-rpc", bytes.NewReader(b))
	switch {
	case err != nil:
		return nil, err
	case resp.StatusCode != http.StatusOK:
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var zr ZabbixResponse

	if err := dec.Decode(&zr); err != nil {
		return nil, err
	}

	return zr.Result, nil
}

//	GetTrigger retrieves triggered triggers for the specified hosts.
//
// It then sends information about the hosts over the specified channel (outZabbix).
func (c *ZabbixClient) GetTrigger(hosts []map[string]string, svdZbxHosts *brds.SavedHosts) error {

	if hosts == nil {
		return errors.New("hosts is nil")

	}

	svdZbxHosts.RWD.Lock()
	defer svdZbxHosts.RWD.Unlock()

	svdZbxHosts.Hosts = map[string]brds.ZabbixHost{}
	for _, hst := range hosts {
		payload := ZabbixRequest{Jsonrpc: "2.0",
			Method: "trigger.get",

			Params: ZabbixParams{
				Output:                 []string{"status", "value", "description"},
				SearchWildcardsEnabled: true,
				SearchByAny:            true,
				Host:                   hst["host"],
			},
			Auth: c.Result,
			Id:   c.Id,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		resp, err := http.Post(c.URL, "application/json-rpc", bytes.NewReader(b))
		switch {
		case err != nil:
			return err
		case resp.StatusCode != http.StatusOK:
			return errors.New(resp.Status)
		}

		defer resp.Body.Close()
		// val, _ := io.ReadAll(resp.Body)
		// log.Println(string(val))
		dec := json.NewDecoder(resp.Body)
		var zr ZabbixResponse

		if err := dec.Decode(&zr); err != nil {
			return err
		}
		// fmt.Println(zr)
		tempProblems := []string{}
		for _, trgr := range zr.Result {
			if trgr["value"] == "1" {
				tempProblems = append(tempProblems, html.EscapeString(trgr["description"]))
				// svdZbxHosts.Hosts[hst["hostid"]].ProblemZ = append(svdZbxHosts.Hosts[hst["hostid"]].ProblemZ, html.EscapeString(trgr["description"]))
			}
		}
		if len(tempProblems) > 0 {
			svdZbxHosts.Hosts[hst["hostid"]] = brds.ZabbixHost{
				HostidZ:  hst["hostid"],
				HostZ:    hst["host"],
				NameZ:    hst["name"],
				StatusZ:  hst["status"],
				ProblemZ: tempProblems,
			}
		}
	}

	return nil
}

// Run launches the Zabbix API client.
func Run(username string, password string, svdZbxHosts *brds.SavedHosts) {
	for {
		client := ZabbixClient{Username: username, Password: password, URL: zabbixUrlAPI}

		if err := client.Authentication(); err != nil {
			log.Fatal(err)
		}

		hosts, err := client.GetHost(WILDCARD)
		if err != nil {
			log.Fatal(err)
		}

		if err := client.GetTrigger(hosts, svdZbxHosts); err != nil {
			log.Fatal(err)
		}

		time.Sleep(3 * time.Minute)
		log.Println("zbx is awake")

	}
	// close(outZabbix)
}
