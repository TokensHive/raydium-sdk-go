package module

import "github.com/TokensHive/raydium-sdk-go/common"

type Currency struct {
	Symbol   string
	Name     string
	Decimals int
}

var SOL = Currency{
	Symbol:   common.SOLInfo.Symbol,
	Name:     common.SOLInfo.Name,
	Decimals: common.SOLInfo.Decimals,
}

func NewCurrency(decimals int, symbol string, name string) Currency {
	if symbol == "" {
		symbol = "UNKNOWN"
	}
	if name == "" {
		name = "UNKNOWN"
	}
	return Currency{
		Symbol:   symbol,
		Name:     name,
		Decimals: decimals,
	}
}

func (c Currency) Equals(other Currency) bool {
	return c == other
}

func CurrencyEquals(currencyA CurrencyLike, currencyB CurrencyLike) bool {
	tokenA, aIsToken := currencyA.(Token)
	tokenB, bIsToken := currencyB.(Token)
	if aIsToken && bIsToken {
		return tokenA.Equals(tokenB)
	}
	if aIsToken || bIsToken {
		return false
	}
	return currencyA.AsCurrency() == currencyB.AsCurrency()
}

type CurrencyLike interface {
	AsCurrency() Currency
}

func (c Currency) AsCurrency() Currency {
	return c
}
