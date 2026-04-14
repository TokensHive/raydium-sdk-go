package module

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/TokensHive/raydium-sdk-go/common"
)

type Fraction struct {
	Numerator   *big.Int
	Denominator *big.Int
}

func NewFraction(numerator any, denominator ...any) (Fraction, error) {
	num, err := common.ParseBigNumberish(numerator)
	if err != nil {
		return Fraction{}, err
	}
	den := big.NewInt(1)
	if len(denominator) > 0 {
		den, err = common.ParseBigNumberish(denominator[0])
		if err != nil {
			return Fraction{}, err
		}
	}
	if den.Sign() == 0 {
		return Fraction{}, fmt.Errorf("denominator cannot be zero")
	}
	return Fraction{
		Numerator:   num,
		Denominator: den,
	}, nil
}

func MustFraction(numerator any, denominator ...any) Fraction {
	f, err := NewFraction(numerator, denominator...)
	if err != nil {
		panic(err)
	}
	return f
}

func (f Fraction) Quotient() *big.Int {
	return new(big.Int).Quo(f.Numerator, f.Denominator)
}

func (f Fraction) Invert() Fraction {
	return Fraction{
		Numerator:   common.CloneBigInt(f.Denominator),
		Denominator: common.CloneBigInt(f.Numerator),
	}
}

func asFraction(other any) (Fraction, error) {
	if v, ok := other.(Fraction); ok {
		return v, nil
	}
	if v, ok := other.(*Fraction); ok {
		return *v, nil
	}
	parsed, err := common.ParseBigNumberish(other)
	if err != nil {
		return Fraction{}, err
	}
	return Fraction{Numerator: parsed, Denominator: big.NewInt(1)}, nil
}

func (f Fraction) Add(other any) (Fraction, error) {
	otherParsed, err := asFraction(other)
	if err != nil {
		return Fraction{}, err
	}
	if f.Denominator.Cmp(otherParsed.Denominator) == 0 {
		return Fraction{
			Numerator:   new(big.Int).Add(f.Numerator, otherParsed.Numerator),
			Denominator: common.CloneBigInt(f.Denominator),
		}, nil
	}

	left := new(big.Int).Mul(f.Numerator, otherParsed.Denominator)
	right := new(big.Int).Mul(otherParsed.Numerator, f.Denominator)
	return Fraction{
		Numerator:   new(big.Int).Add(left, right),
		Denominator: new(big.Int).Mul(f.Denominator, otherParsed.Denominator),
	}, nil
}

func (f Fraction) Sub(other any) (Fraction, error) {
	otherParsed, err := asFraction(other)
	if err != nil {
		return Fraction{}, err
	}
	if f.Denominator.Cmp(otherParsed.Denominator) == 0 {
		return Fraction{
			Numerator:   new(big.Int).Sub(f.Numerator, otherParsed.Numerator),
			Denominator: common.CloneBigInt(f.Denominator),
		}, nil
	}

	left := new(big.Int).Mul(f.Numerator, otherParsed.Denominator)
	right := new(big.Int).Mul(otherParsed.Numerator, f.Denominator)
	return Fraction{
		Numerator:   new(big.Int).Sub(left, right),
		Denominator: new(big.Int).Mul(f.Denominator, otherParsed.Denominator),
	}, nil
}

func (f Fraction) Mul(other any) (Fraction, error) {
	otherParsed, err := asFraction(other)
	if err != nil {
		return Fraction{}, err
	}
	return Fraction{
		Numerator:   new(big.Int).Mul(f.Numerator, otherParsed.Numerator),
		Denominator: new(big.Int).Mul(f.Denominator, otherParsed.Denominator),
	}, nil
}

func (f Fraction) Div(other any) (Fraction, error) {
	otherParsed, err := asFraction(other)
	if err != nil {
		return Fraction{}, err
	}
	if otherParsed.Numerator.Sign() == 0 {
		return Fraction{}, fmt.Errorf("division by zero")
	}
	return Fraction{
		Numerator:   new(big.Int).Mul(f.Numerator, otherParsed.Denominator),
		Denominator: new(big.Int).Mul(f.Denominator, otherParsed.Numerator),
	}, nil
}

func (f Fraction) ToSignificant(significantDigits int32, rounding common.Rounding) (string, error) {
	if significantDigits <= 0 {
		return "", fmt.Errorf("%d is not positive", significantDigits)
	}
	num := decimal.NewFromBigInt(f.Numerator, 0)
	den := decimal.NewFromBigInt(f.Denominator, 0)
	if den.IsZero() {
		return "", fmt.Errorf("denominator cannot be zero")
	}
	value := num.Div(den)
	floatValue, _ := value.Float64()
	return strconv.FormatFloat(floatValue, 'g', int(significantDigits), 64), nil
}

func (f Fraction) ToFixed(decimalPlaces int32, rounding common.Rounding) (string, error) {
	if decimalPlaces < 0 {
		return "", fmt.Errorf("%d is negative", decimalPlaces)
	}
	num := decimal.NewFromBigInt(f.Numerator, 0)
	den := decimal.NewFromBigInt(f.Denominator, 0)
	if den.IsZero() {
		return "", fmt.Errorf("denominator cannot be zero")
	}
	value := num.Div(den)
	switch rounding {
	case common.RoundDown:
		value = value.Truncate(decimalPlaces)
	case common.RoundUp:
		value = value.RoundCeil(decimalPlaces)
	default:
		value = value.Round(decimalPlaces)
	}
	return value.StringFixed(decimalPlaces), nil
}

func (f Fraction) IsZero() bool {
	return f.Numerator.Sign() == 0
}
