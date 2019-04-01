package analysis

import "github.com/snickers54/crypto-trading-bot/api"

func UpdateProducts() {
	Knowledge.Mutex.Lock()
	products := api.GetListProducts()
	for _, product := range *products {
		Knowledge.Products[product.ID] = product
	}
	Knowledge.Mutex.Unlock()
}
