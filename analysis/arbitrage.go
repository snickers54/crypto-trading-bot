package analysis

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sdcoffey/big"
	"github.com/snickers54/crypto-trading-bot/api"
)

type ratesMap map[string]map[string]big.Decimal

func checkRatesReady(pair1, pair2, pair3 string) bool {
	r1, ok1 := Knowledge.Rates[pair1]
	r2, ok2 := Knowledge.Rates[pair2]
	r3, ok3 := Knowledge.Rates[pair3]
	return ok1 && ok2 && ok3 && r1.Ready && r2.Ready && r3.Ready
}

func ListenDiscrepancy(pair string, ordersTube chan []api.MarketOrderPayload) {
	fiat := "EUR"
	arr := append(strings.Split(pair, "-"), fiat)
	pair1 := fmt.Sprintf("%s-%s", arr[0], arr[1])
	pair2 := fmt.Sprintf("%s-%s", arr[0], arr[2])
	pair3 := fmt.Sprintf("%s-%s", arr[1], arr[2])
	for {
		Knowledge.Mutex.Lock()
		if len(arr) != 3 || checkRatesReady(pair1, pair2, pair3) == false {
			Knowledge.Mutex.Unlock()
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		// let's try for the normal way
		rate1 := Knowledge.Rates[pair2].RateBuy
		rate2 := Knowledge.Rates[pair1].RateSell
		rate3 := Knowledge.Rates[pair3].RateSell
		rate := rate1.Mul(rate2).Mul(rate3)
		if rate.GT(big.ONE) {
			fmt.Printf("[%s (%s) / %s (%s) / %s (%s)] RATE(%s)\n", pair2, Knowledge.Rates[pair2].BestBid, pair1, Knowledge.Rates[pair1].BestAsk, pair3, Knowledge.Rates[pair3].BestAsk, rate.String())
			// ETH - BTC - EUR which gives us: ETH-EUR, ETH-BTC, BTC-EUR
			orders := prepareOrders(arr)
			for _, order := range orders {
				order.Print()
			}
			// ordersTube <- orders
		}
		Knowledge.Mutex.Unlock()
		time.Sleep(500 * time.Millisecond)
	}
}

func prepareOrders(curr []string) []api.MarketOrderPayload {
	// ETH-EUR so I BUY ETH FROM EUR
	pair := fmt.Sprintf("%s-%s", curr[0], curr[2])
	product := Knowledge.Products[pair]
	rate := Knowledge.Rates[pair]
	qty := big.NewFromString(os.Getenv("INITIAL_BET")).Mul(rate.RateSell)
	initialOrder := api.NewOrderBuy(qty, product)

	// ETH-BTC so I SELL ETH TO BUY BTC
	pair = fmt.Sprintf("%s-%s", curr[0], curr[1])
	product = Knowledge.Products[pair]
	rate = Knowledge.Rates[pair]
	cryptoOrder := api.NewOrderSell(qty, product)
	qty = qty.Mul(rate.RateBuy)

	// BTC-EUR so I SELL BTC TO BUY EUR
	pair = fmt.Sprintf("%s-%s", curr[1], curr[2])
	product = Knowledge.Products[pair]
	finalOrder := api.NewOrderSell(qty, product)

	return []api.MarketOrderPayload{
		initialOrder,
		cryptoOrder,
		finalOrder,
	}
}
