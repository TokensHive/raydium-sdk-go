package common

import (
	"fmt"
	"math/big"
	"strings"
)

var (
	BNZero  = big.NewInt(0)
	BNOne   = big.NewInt(1)
	BNTwo   = big.NewInt(2)
	BNThree = big.NewInt(3)
	BNFive  = big.NewInt(5)
	BNTen   = big.NewInt(10)
	BN100   = big.NewInt(100)
	BN1000  = big.NewInt(1000)
	BN10000 = big.NewInt(10000)
)

func CloneBigInt(v *big.Int) *big.Int {
	if v == nil {
		return big.NewInt(0)
	}
	return new(big.Int).Set(v)
}

func TenExponential(shift *big.Int) *big.Int {
	return new(big.Int).Exp(BNTen, shift, nil)
}

func parseSignedDecimalString(value string) (*big.Int, error) {
	if value == "" {
		return nil, fmt.Errorf("empty numeric string")
	}
	if strings.Contains(value, ".") {
		return nil, fmt.Errorf("decimal point not supported: %s", value)
	}
	n := new(big.Int)
	if _, ok := n.SetString(value, 10); !ok {
		return nil, fmt.Errorf("invalid integer string: %s", value)
	}
	return n, nil
}
