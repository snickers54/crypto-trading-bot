package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/imroc/req"
)

type ConversionPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type ConversionResponse struct {
	ID            string `json:"id"`
	Amount        string `json:"amount"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	From          string `json:"from"`
	To            string `json:"to"`
}

// ConvertTo is only able to convert Fiat to Stablecoin or Stablecoin to Fiat
func ConvertTo(currencyFromCode, currencyToCode string, amount float64) *ConversionResponse {
	payload := ConversionPayload{
		From:   currencyFromCode,
		To:     currencyToCode,
		Amount: fmt.Sprintf("%.2f", amount),
	}
	conversion := ConversionResponse{}
	path := "/conversions"
	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	headers := authHeaders(string(body), "POST", path)
	resp, err := req.Post(os.Getenv("API_URL")+path, headers, req.BodyJSON(body))
	if err != nil || resp.Response().StatusCode != 200 {
		fmt.Println(err)
		printError(resp)
		return nil
	}
	resp.ToJSON(&conversion)
	return &conversion
}
