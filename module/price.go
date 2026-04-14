package module

import (
	"fmt"
	"math/big"

	"github.com/TokensHive/raydium-sdk-go/common"
)

type Price struct {
	Fraction
	BaseToken  Token
	QuoteToken Token
	Scalar     Fraction
}

type PriceProps struct {
	BaseToken   Token
	QuoteToken  Token
	Numerator   any
	Denominator any
}

func NewPrice(params PriceProps) (Price, error) {
	f, err := NewFraction(params.Numerator, params.Denominator)
	if err != nil {
		return Price{}, err
	}
	scalar := Fraction{
		Numerator:   new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.BaseToken.Decimals)), nil),
		Denominator: new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(params.QuoteToken.Decimals)), nil),
	}
	return Price{
		Fraction:   f,
		BaseToken:  params.BaseToken,
		QuoteToken: params.QuoteToken,
		Scalar:     scalar,
	}, nil
}

func (p Price) Raw() Fraction {
	return Fraction{
		Numerator:   common.CloneBigInt(p.Numerator),
		Denominator: common.CloneBigInt(p.Denominator),
	}
}

func (p Price) Adjusted() Fraction {
	adjusted, _ := p.Fraction.Mul(p.Scalar)
	return adjusted
}

func (p Price) Invert() (Price, error) {
	return NewPrice(PriceProps{
		BaseToken:   p.QuoteToken,
		QuoteToken:  p.BaseToken,
		Denominator: p.Numerator,
		Numerator:   p.Denominator,
	})
}

func (p Price) Mul(other Price) (Price, error) {
	if !p.QuoteToken.Equals(other.BaseToken) {
		return Price{}, fmt.Errorf("mul token not equals")
	}
	fraction, _ := p.Fraction.Mul(other.Fraction)
	return NewPrice(PriceProps{
		BaseToken:   p.BaseToken,
		QuoteToken:  other.QuoteToken,
		Denominator: fraction.Denominator,
		Numerator:   fraction.Numerator,
	})
}

func (p Price) ToSignificant(significantDigits int32, rounding common.Rounding) (string, error) {
	adjusted := p.Adjusted()
	if significantDigits == 0 {
		significantDigits = int32(p.QuoteToken.Decimals)
	}
	return adjusted.ToSignificant(significantDigits, rounding)
}

func (p Price) ToFixed(decimalPlaces int32, rounding common.Rounding) (string, error) {
	adjusted := p.Adjusted()
	if decimalPlaces == 0 {
		decimalPlaces = int32(p.QuoteToken.Decimals)
	}
	return adjusted.ToFixed(decimalPlaces, rounding)
}
