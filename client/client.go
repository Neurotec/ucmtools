package client

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net/http"
	"net/url"
)

type Client struct {
	Server  string
	apiUrl  string
	request *gorequest.SuperAgent
}

func New(baseUrl string) (*Client, error) {
	request := gorequest.New()
	resp, _, _ := request.Post(baseUrl + "/cgi").
		Send(`{"action":"getInfo"}`).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		End()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("can't connect to ucm")
	}

	return &Client{resp.Header.Get("Server"), baseUrl + "/cgi", request}, nil
}

func (c *Client) Login(user, password string) bool {
	args := url.Values{}
	args.Set("user", user)

	resp, err := c.action("challenge", args)

	if err != nil {
		return false
	}

	if resp["status"].(float64) != 0 {
		return false
	}

	challenge := resp["response"].(map[string]interface{})["challenge"]
	args.Set("token", fmt.Sprintf("%x",
		md5.Sum([]byte(challenge.(string)+password))))
	respLogin, err := c.action("login", args)
	if err != nil {
		return false
	}
	if respLogin["status"].(float64) != 0 {
		return false
	}
	return true
}

func (c *Client) ListAccounts() []Account {
	var accounts []Account

	vals := url.Values{}
	vals.Set("item_num", "100")
	vals.Set("page", "1")
	vals.Set("sidx", "extension")
	vals.Set("sord", "asc")
	resp, err := c.action("listAccount", vals)
	if err != nil {
		return accounts
	}

	if resp["status"].(float64) != 0 {
		return accounts
	}

	accountsEncode := resp["response"].(map[string]interface{})["account"]
	accountsDecode := accountsEncode.([]interface{})
	for _, accountEncode := range accountsDecode {
		account := accountEncode.(map[string]interface{})
		accounts = append(accounts, Account{
			Type:         account["account_type"].(string),
			Extension:    account["extension"].(string),
			Status:       account["status"].(string),
			OutOfService: account["out_of_service"].(string),
		})
	}

	return accounts
}

func (c *Client) ListTrunks() []Trunk {
	var trunks []Trunk

	vals := url.Values{}
	vals.Set("item_num", "100")
	vals.Set("page", "1")
	vals.Set("sidx", "status")
	vals.Set("sord", "asc")
	resp, err := c.action("listAllTrunk", vals)
	if err != nil {
		return trunks
	}
	if resp["status"].(float64) != 0 {
		return trunks
	}

	trunksEncode := resp["response"].(map[string]interface{})["trunks"]
	trunksDecode := trunksEncode.([]interface{})
	for _, trunkEncode := range trunksDecode {
		trunk := trunkEncode.(map[string]interface{})
		trunks = append(trunks, Trunk{
			Type:         trunk["type"].(string),
			Username:     trunk["username"].(string),
			TrunkName:    trunk["trunk_name"].(string),
			Status:       trunk["status"].(string),
			OutOfService: trunk["out_of_service"].(string),
		})
	}

	return trunks
}

func (c *Client) action(name string, vals url.Values) (map[string]interface{}, error) {
	vals.Set("action", name)
	resp, _, errs := c.request.Post(c.apiUrl).SendString(vals.Encode()).End()
	if errs != nil {
		return nil, errors.New("failed do action " + name)
	}
	var decode map[string]interface{}

	err := json.NewDecoder(resp.Body).Decode(&decode)
	if err != nil {
		return nil, err
	}

	return decode, nil
}

type Trunk struct {
	Type         string `json:"type"`
	Username     string `json:"username"`
	TrunkName    string `json:"trunk_name"`
	Status       string `json:"status"`
	OutOfService string `json:"out_of_service"`
}

type Account struct {
	Extension    string `json:"extension"`
	Type         string `json:"account_type"`
	Status       string `json:"status"`
	OutOfService string `json:"out_of_service"`
}
