package analysis

import (
	"encoding/json"
	"fmt"

	"github.com/sdcoffey/big"
	"github.com/snickers54/trading-bot/api"
)

func RealtimeFeed() {
	messages := make(chan []byte, 10)
	done, conn := api.InitWebsocketConnection(messages)
	defer conn.Close()
	if success := api.SubscribePrice([]string{
		"ETH-EUR", "LTC-EUR", "BCH-EUR", "ZRX-EUR", "BTC-EUR", "ETC-EUR", "XRP-EUR",
		"ETH-BTC", "LTC-BTC", "BCH-BTC", "ZRX-BTC", "ETC-BTC", "XRP-BTC", "XLM-BTC", "XLM-EUR",
	}, conn); success == false {
		return
	}
	for {
		select {
		case <-done:
			return
		case data := <-messages:
			level2Message := L2Message{}
			if err := json.Unmarshal(data, &level2Message); err != nil {
				fmt.Println(err)
				fmt.Println(string(data))
				continue
			}
			if level2Message.Type == "l2update" || level2Message.Type == "snapshot" {
				Knowledge.Mutex.Lock()
				UpdateBooks(level2Message)
				Knowledge.Mutex.Unlock()
			}
		}
	}
}

func UpdateBooks(l2 L2Message) {
	pair := l2.ProductID
	if _, ok := Knowledge.Rates[pair]; ok != true {
		Knowledge.Rates[pair] = Rate{
			Pair:    pair,
			BookBid: map[string]big.Decimal{},
			BookAsk: map[string]big.Decimal{},
			Ready:   false,
		}
	}
	r, _ := Knowledge.Rates[pair]
	if l2.Type == "snapshot" {
		for _, tuple := range l2.Bids {
			r.BookBid[tuple[0]] = big.NewFromString(tuple[0])
		}
		for _, tuple := range l2.Asks {
			r.BookAsk[tuple[0]] = big.NewFromString(tuple[0])
		}
	} else if l2.Type == "l2update" {
		for _, tuple := range l2.Changes {
			side := tuple[0]
			value := big.NewFromString(tuple[1])
			size := tuple[2]
			if side == "buy" {
				if size == "0" {
					delete(r.BookBid, tuple[1])
				} else {
					r.BookBid[tuple[1]] = value
				}
			} else if side == "sell" {
				if size == "0" {
					delete(r.BookAsk, tuple[1])
				} else {
					r.BookAsk[tuple[1]] = value
				}
			}
		}
	}
	r.UpdateAsk()
	r.UpdateBid()
	r.Ready = true
	Knowledge.Rates[pair] = r
}
