package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type AccountHistoryDetails struct {
	OrderID   string `json:"order_id"`
	TradeID   string `json:"trade_id"`
	ProductID string `json:"product_id"`
}

type AccountHistory struct {
	ID        string                `json:"id"`
	CreatedAt string                `json:"created_at"`
	Amount    string                `json:"amount"`
	Balance   string                `json:"balance"`
	Type      string                `json:"type"`
	Details   AccountHistoryDetails `json:"details"`
}
type AccountHistories []AccountHistory

func GetAccountHistories(accounID string) *AccountHistories {
	path := fmt.Sprintf("/accounts/%s/ledger", accounID)
	headers := authHeaders("", "GET", path)
	accountHistories := AccountHistories{}
	resp, err := req.Get(os.Getenv("API_URL")+path, headers)
	if err != nil || resp.Response().StatusCode != 200 {
		fmt.Println(err)
		printError(resp)
		return nil
	}
	resp.ToJSON(&accountHistories)
	return &accountHistories
}
