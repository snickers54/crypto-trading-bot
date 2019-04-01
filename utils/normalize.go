package utils

import (
	"strconv"
	"strings"

	"github.com/sdcoffey/big"
)

func NormalizeSize(quantity big.Decimal, increment string) string {
	places := 0
	i, _ := strconv.ParseInt(increment, 10, 64)
	if i < 1 {
		places = len(strings.Split(increment, ".")[1])
	}
	qty := quantity.Sub(big.NewFromString(increment).Div(big.NewDecimal(2)))
	return qty.FormattedString(places)
}
