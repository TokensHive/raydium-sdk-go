package module

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/TokensHive/raydium-sdk-go/common"
)

func SplitNumber(num string, decimals int) (string, string, error) {
	integral := "0"
	fractional := "0"
	if strings.Contains(num, ".") {
		parts := strings.Split(num, ".")
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid number string: %s", num)
		}
		integral = parts[0]
		fractional = parts[1]
		if len(fractional) < decimals {
			fractional += strings.Repeat("0", decimals-len(fractional))
		}
	} else {
		integral = num
	}
	if decimals == 0 {
		return integral, "0", nil
	}
	if len(fractional) > decimals {
		fractional = fractional[:decimals]
	}
	return integral, fractional, nil
}

type TokenAmount struct {
	Fraction
	Token Token
}

func NewTokenAmount(token Token, amount any, isRaw bool) (TokenAmount, error) {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(token.Decimals)), nil)
	parsedAmount := big.NewInt(0)
	if isRaw {
		v, err := common.ParseBigNumberish(amount)
		if err != nil {
			return TokenAmount{}, err
		}
		parsedAmount = v
	} else {
		integral, fractional, err := SplitNumber(fmt.Sprint(amount), token.Decimals)
		if err != nil {
			return TokenAmount{}, err
		}
		integralAmount, err := common.ParseBigNumberish(integral)
		if err != nil {
			return TokenAmount{}, err
		}
		fractionalAmount, err := common.ParseBigNumberish(fractional)
		if err != nil {
			return TokenAmount{}, err
		}
		integralAmount = integralAmount.Mul(integralAmount, multiplier)
		parsedAmount = parsedAmount.Add(integralAmount, fractionalAmount)
	}
	f := Fraction{
		Numerator:   parsedAmount,
		Denominator: multiplier,
	}
	return TokenAmount{
		Fraction: f,
		Token:    token,
	}, nil
}

func (a TokenAmount) Raw() *big.Int {
	return a.Numerator
}

func (a TokenAmount) IsZero() bool {
	return a.Raw().Sign() == 0
}

func (a TokenAmount) GT(other TokenAmount) error {
	if !a.Token.Equals(other.Token) {
		return fmt.Errorf("gt token not equals")
	}
	if a.Raw().Cmp(other.Raw()) <= 0 {
		return fmt.Errorf("not greater than")
	}
	return nil
}

func (a TokenAmount) LT(other TokenAmount) error {
	if !a.Token.Equals(other.Token) {
		return fmt.Errorf("lt token not equals")
	}
	if a.Raw().Cmp(other.Raw()) >= 0 {
		return fmt.Errorf("not less than")
	}
	return nil
}

func (a TokenAmount) Add(other TokenAmount) (TokenAmount, error) {
	if !a.Token.Equals(other.Token) {
		return TokenAmount{}, fmt.Errorf("add token not equals")
	}
	return NewTokenAmount(a.Token, new(big.Int).Add(a.Raw(), other.Raw()), true)
}

func (a TokenAmount) Subtract(other TokenAmount) (TokenAmount, error) {
	if !a.Token.Equals(other.Token) {
		return TokenAmount{}, fmt.Errorf("sub token not equals")
	}
	return NewTokenAmount(a.Token, new(big.Int).Sub(a.Raw(), other.Raw()), true)
}

func (a TokenAmount) ToSignificant(significantDigits int32, rounding common.Rounding) (string, error) {
	if significantDigits == 0 {
		significantDigits = int32(a.Token.Decimals)
	}
	return a.Fraction.ToSignificant(significantDigits, rounding)
}

func (a TokenAmount) ToFixed(decimalPlaces int32, rounding common.Rounding) (string, error) {
	if decimalPlaces == 0 {
		decimalPlaces = int32(a.Token.Decimals)
	}
	if decimalPlaces > int32(a.Token.Decimals) {
		return "", fmt.Errorf("decimals overflow")
	}
	return a.Fraction.ToFixed(decimalPlaces, rounding)
}

func (a TokenAmount) ToExact() string {
	num := decimal.NewFromBigInt(a.Numerator, 0)
	den := decimal.NewFromBigInt(a.Denominator, 0)
	return num.Div(den).String()
}

type CurrencyAmount struct {
	Fraction
	Currency Currency
}

func NewCurrencyAmount(currency Currency, amount any, isRaw bool) (CurrencyAmount, error) {
	tokenLike := Token{
		Symbol:   currency.Symbol,
		Name:     currency.Name,
		Decimals: currency.Decimals,
	}
	tokenAmount, err := NewTokenAmount(tokenLike, amount, isRaw)
	if err != nil {
		return CurrencyAmount{}, err
	}
	return CurrencyAmount{
		Fraction: tokenAmount.Fraction,
		Currency: currency,
	}, nil
}

func (a CurrencyAmount) Raw() *big.Int {
	return a.Numerator
}

func (a CurrencyAmount) IsZero() bool {
	return a.Raw().Sign() == 0
}

func (a CurrencyAmount) GT(other CurrencyAmount) error {
	if !a.Currency.Equals(other.Currency) {
		return fmt.Errorf("gt currency not equals")
	}
	if a.Raw().Cmp(other.Raw()) <= 0 {
		return fmt.Errorf("not greater than")
	}
	return nil
}

func (a CurrencyAmount) LT(other CurrencyAmount) error {
	if !a.Currency.Equals(other.Currency) {
		return fmt.Errorf("lt currency not equals")
	}
	if a.Raw().Cmp(other.Raw()) >= 0 {
		return fmt.Errorf("not less than")
	}
	return nil
}

func (a CurrencyAmount) Add(other CurrencyAmount) (CurrencyAmount, error) {
	if !a.Currency.Equals(other.Currency) {
		return CurrencyAmount{}, fmt.Errorf("add currency not equals")
	}
	return NewCurrencyAmount(a.Currency, new(big.Int).Add(a.Raw(), other.Raw()), true)
}

func (a CurrencyAmount) Sub(other CurrencyAmount) (CurrencyAmount, error) {
	if !a.Currency.Equals(other.Currency) {
		return CurrencyAmount{}, fmt.Errorf("sub currency not equals")
	}
	return NewCurrencyAmount(a.Currency, new(big.Int).Sub(a.Raw(), other.Raw()), true)
}

func (a CurrencyAmount) ToSignificant(significantDigits int32, rounding common.Rounding) (string, error) {
	if significantDigits == 0 {
		significantDigits = int32(a.Currency.Decimals)
	}
	return a.Fraction.ToSignificant(significantDigits, rounding)
}

func (a CurrencyAmount) ToFixed(decimalPlaces int32, rounding common.Rounding) (string, error) {
	if decimalPlaces == 0 {
		decimalPlaces = int32(a.Currency.Decimals)
	}
	if decimalPlaces > int32(a.Currency.Decimals) {
		return "", fmt.Errorf("decimals overflow")
	}
	return a.Fraction.ToFixed(decimalPlaces, rounding)
}

func (a CurrencyAmount) ToExact() string {
	num := decimal.NewFromBigInt(a.Numerator, 0)
	den := decimal.NewFromBigInt(a.Denominator, 0)
	return num.Div(den).String()
}
