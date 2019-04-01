package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sdcoffey/big"
	"github.com/snickers54/trading-bot/utils"

	"github.com/imroc/req"
)

type OrderCommonPayload struct {
	ClientOID string `json:"client_oid"`
	Type      string `json:"type"`
	Side      string `json:"side"`
	ProductID string `json:"product_id"`
	STP       string `json:"stp"`
	Stop      string `json:"stop,omitempty"`
	StopPrice string `json:"stop_price,omitempty"`
}

type LimitOrderPayload struct {
	OrderCommonPayload
	// Limit order params
	Price       string `json:"price"`
	Size        string `json:"size"`
	TimeInForce string `json:"time_in_force,omitempty"`
	CancelAfter string `json:"cancel_after,omitempty"`
	PostOnly    bool   `json:"post_only,omitempty"`
}

func (lop MarketOrderPayload) Print() {
	curr := strings.Split(lop.OrderCommonPayload.ProductID, "-")
	fmt.Printf("[%s] %s %s\n",
		strings.ToUpper(lop.OrderCommonPayload.Side),
		lop.Size, curr[0])
}

type MarketOrderPayload struct {
	OrderCommonPayload
	Size  string `json:"size,omitempty"`
	Funds string `json:"funds,omitempty"`
}

type Order struct {
	ID            string `json:"id"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	ProductID     string `json:"product_id"`
	Side          string `json:"side"`
	STP           string `json:"stp"`
	Type          string `json:"type"`
	DoneAt        string `json:"done_at,omitEmpty"`
	DoneReason    string `json:"done_reason"`
	TimeInForce   string `json:"time_in_force,omitempty"`
	PostOnly      string `json:"post_only,omitempty"`
	CreatedAt     string `json:"created_at"`
	FillFees      string `json:"fill_fees"`
	FilledSize    string `json:"filled_size"`
	ExecutedValue string `json:"executed_value"`
	Status        string `json:"status"`
	Settled       bool   `json:"settled"`
}

func MakeOrder(orderPayload interface{}) *Order {
	defer utils.TimeTrack(time.Now(), "MakeOrder")
	orderResponse := Order{}
	path := "/orders"
	body, err := json.Marshal(orderPayload)
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
	resp.ToJSON(&orderResponse)
	return &orderResponse

}

func CancelAllOrderPerProduct(productID string) {
	defer utils.TimeTrack(time.Now(), "DeleteOrderPerProduct")
	path := fmt.Sprintf("/orders")
	headers := authHeaders("", "DELETE", path)
	resp, err := req.Delete(os.Getenv("API_URL")+path, headers)
	if err != nil || resp.Response().StatusCode != 200 {
		fmt.Println(err)
		printError(resp)
		return
	}
	return
}

func NewOrderBuy(quantity big.Decimal, product Product) MarketOrderPayload {
	return MarketOrderPayload{
		OrderCommonPayload: OrderCommonPayload{
			ClientOID: uuid.New().String(),
			Type:      "market",
			Side:      "buy",
			ProductID: product.ID,
		},
		Size: utils.NormalizeSize(quantity, product.BaseMinSize),
	}
}

func NewOrderSell(quantity big.Decimal, product Product) MarketOrderPayload {
	return MarketOrderPayload{
		OrderCommonPayload: OrderCommonPayload{
			ClientOID: uuid.New().String(),
			Type:      "market",
			Side:      "sell",
			ProductID: product.ID,
		},
		Size: utils.NormalizeSize(quantity, product.BaseMinSize),
	}
}
