package module

import "github.com/TokensHive/raydium-sdk-go/common"

var OneHundredPercent = MustFraction(100)

type Percent struct {
	Fraction
}

func NewPercent(numerator any, denominator ...any) (Percent, error) {
	f, err := NewFraction(numerator, denominator...)
	if err != nil {
		return Percent{}, err
	}
	return Percent{Fraction: f}, nil
}

func (p Percent) ToSignificant(significantDigits int32, rounding common.Rounding) (string, error) {
	multiplied, _ := p.Fraction.Mul(OneHundredPercent)
	return multiplied.ToSignificant(significantDigits, rounding)
}

func (p Percent) ToFixed(decimalPlaces int32, rounding common.Rounding) (string, error) {
	multiplied, _ := p.Fraction.Mul(OneHundredPercent)
	return multiplied.ToFixed(decimalPlaces, rounding)
}
