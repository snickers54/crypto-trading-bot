package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type Account struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Hold      string `json:"hold"`
	ProfileID string `json:"profile_id"`
}
type Accounts []Account

func GetAccounts() *Accounts {
	path := "/accounts"
	headers := authHeaders("", "GET", path)
	accounts := Accounts{}
	resp, err := req.Get(os.Getenv("API_URL")+path, headers)
	if err != nil || resp.Response().StatusCode != 200 {
		fmt.Println(err)
		printError(resp)
		return nil
	}
	resp.ToJSON(&accounts)
	return &accounts
}
