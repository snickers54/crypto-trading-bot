package api

import (
	"fmt"
	"os"

	"github.com/imroc/req"
)

type Product struct {
	ID             string `json:"id"`
	BaseCurrency   string `json:"base_currency"`
	QuoteCurrency  string `json:"quote_currency"`
	BaseMinSize    string `json:"base_min_size"`
	BaseMaxSize    string `json:"base_max_size"`
	QuoteIncrement string `json:"quote_increment"`

	DisplayName   string `json:"display_name"`
	Status        string `json:"status"`
	MarginEnabled bool   `json:"margin_enabled"`
	StatusMessage string `json:"status_message,omitempty"`
}

type Products []Product

func GetListProducts() *Products {
	path := "/products"
	products := Products{}
	headers := authHeaders("", "GET", path)
	resp, err := req.Get(os.Getenv("API_URL")+path, headers)
	if err != nil || resp.Response().StatusCode != 200 {
		fmt.Println(err)
		printError(resp)
		return nil
	}
	resp.ToJSON(&products)
	return &products
}
