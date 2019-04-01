package utils

import (
	"testing"

	"github.com/sdcoffey/big"
)

func TestNormalizeSize(t *testing.T) {
	norm := NormalizeSize("1.00", "0.001")
	if norm.EQ(big.ONE) == false {
		t.Errorf("Expecting %s, got %s", big.ONE.String(), norm.String())
	}
}
