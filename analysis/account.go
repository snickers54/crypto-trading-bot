package analysis

import (
	"os"
	"time"

	"github.com/snickers54/trading-bot/api"
)

func UpdateAccounts() {
	for {
		accounts := api.GetAccounts()
		if accounts == nil {
			os.Exit(0)
		}
		Knowledge.Mutex.Lock()
		for _, account := range *accounts {
			Knowledge.Accounts[account.Currency] = account
			if account.Available != "0" {
			}
		}
		Knowledge.Mutex.Unlock()
		time.Sleep(1000 * time.Millisecond)
	}
}
