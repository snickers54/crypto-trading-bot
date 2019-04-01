package analysis

import (
	"math"
	"sync"

	"github.com/sdcoffey/big"
	"github.com/snickers54/crypto-trading-bot/api"
)

var MARKET_FEE = big.NewDecimal(0.25).Div(big.NewDecimal(100))

type L2Message struct {
	Type      string     `json:"type"`
	ProductID string     `json:"product_id"`
	Bids      [][]string `json:"bids,omitempty"`
	Asks      [][]string `json:"asks,omitempty"`
	Changes   [][]string `json:"changes,omitempty"`
}

type Rate struct {
	Pair     string
	RateBuy  big.Decimal
	RateSell big.Decimal

	BestBid big.Decimal
	BestAsk big.Decimal

	BookBid map[string]big.Decimal
	BookAsk map[string]big.Decimal
	Ready   bool
}

// the one who want to sell the Base currency
func (r *Rate) UpdateAsk() {
	r.BestAsk = big.NewDecimal(math.Inf(1))
	for _, ask := range r.BookAsk {
		if ask.LT(r.BestAsk) {
			r.BestAsk = ask
			r.RateBuy = big.ONE.Div(ask).Mul(big.ONE.Sub(MARKET_FEE))
		}
	}
	// fmt.Printf("[%s] Best Ask: %s, rate: %s\n", r.Pair, r.BestAsk, r.RateSell)
}

// the one who want to buy the Base currency
func (r *Rate) UpdateBid() {
	r.BestBid = big.NewDecimal(math.Inf(-1))
	for _, bid := range r.BookBid {
		if bid.GT(r.BestBid) {
			r.BestBid = bid
			r.RateSell = bid.Mul(big.ONE.Sub(MARKET_FEE))
		}
	}
	// fmt.Printf("[%s] Best Bid: %s, rate: %s\n", r.Pair, r.BestBid, r.RateBuy)
}

type KnowledgeStruct struct {
	Accounts map[string]api.Account
	Rates    map[string]Rate
	Products map[string]api.Product
	Mutex    sync.Mutex
}

var Knowledge = KnowledgeStruct{
	Accounts: map[string]api.Account{},
	Rates:    map[string]Rate{},
	Products: map[string]api.Product{},
}
