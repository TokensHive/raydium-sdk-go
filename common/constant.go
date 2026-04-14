package common

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

type Rounding int

const (
	RoundDown Rounding = iota
	RoundHalfUp
	RoundUp
)

func ParseBigNumberish(value any) (*big.Int, error) {
	switch v := value.(type) {
	case *big.Int:
		return CloneBigInt(v), nil
	case big.Int:
		return new(big.Int).Set(&v), nil
	case string:
		return parseSignedDecimalString(v)
	case int:
		return big.NewInt(int64(v)), nil
	case int8:
		return big.NewInt(int64(v)), nil
	case int16:
		return big.NewInt(int64(v)), nil
	case int32:
		return big.NewInt(int64(v)), nil
	case int64:
		return big.NewInt(v), nil
	case uint:
		if v > uint(math.MaxInt64) {
			return nil, fmt.Errorf("overflow for uint: %d", v)
		}
		return big.NewInt(int64(v)), nil
	case uint8:
		return big.NewInt(int64(v)), nil
	case uint16:
		return big.NewInt(int64(v)), nil
	case uint32:
		return big.NewInt(int64(v)), nil
	case uint64:
		if v > math.MaxInt64 {
			n := new(big.Int)
			n.SetUint64(v)
			return n, nil
		}
		return big.NewInt(int64(v)), nil
	case float32:
		if math.Trunc(float64(v)) != float64(v) {
			return nil, fmt.Errorf("BigNumberish number underflow: %v", v)
		}
		return ParseBigNumberish(int64(v))
	case float64:
		if math.Trunc(v) != v {
			return nil, fmt.Errorf("BigNumberish number underflow: %v", v)
		}
		if v >= float64(math.MaxInt64) || v <= float64(math.MinInt64) {
			return nil, fmt.Errorf("BigNumberish number overflow: %v", v)
		}
		return big.NewInt(int64(v)), nil
	case fmt.Stringer:
		return parseSignedDecimalString(v.String())
	default:
		return nil, fmt.Errorf("invalid BigNumberish value: %T", value)
	}
}

func SplitSignedDecimal(value string) (sign string, intPart string, fracPart string, err error) {
	if value == "" {
		return "", "", "", fmt.Errorf("empty value")
	}
	sign = ""
	rest := value
	if value[0] == '-' {
		sign = "-"
		rest = value[1:]
	}
	if rest == "" {
		return "", "", "", fmt.Errorf("invalid signed decimal: %s", value)
	}
	for i, c := range rest {
		if c == '.' {
			intPart = rest[:i]
			fracPart = rest[i+1:]
			if intPart == "" {
				intPart = "0"
			}
			return sign, intPart, fracPart, nil
		}
	}
	intPart = rest
	fracPart = ""
	return sign, intPart, fracPart, nil
}

func ParseIntBase10(value string) (*big.Int, error) {
	n := new(big.Int)
	if _, ok := n.SetString(value, 10); !ok {
		return nil, fmt.Errorf("invalid integer string: %s", value)
	}
	return n, nil
}

func MustParseIntBase10(value string) *big.Int {
	v, err := ParseIntBase10(value)
	if err != nil {
		panic(err)
	}
	return v
}

func DecimalStringToRational(value string) (*big.Int, *big.Int, error) {
	sign, intPart, fracPart, err := SplitSignedDecimal(value)
	if err != nil {
		return nil, nil, err
	}
	denominator := big.NewInt(1)
	if fracPart != "" {
		denominator = new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(fracPart))), nil)
	}
	numeratorString := sign + intPart + fracPart
	numerator, err := ParseIntBase10(numeratorString)
	if err != nil {
		return nil, nil, err
	}
	return numerator, denominator, nil
}

func BigIntToString(v *big.Int) string {
	if v == nil {
		return "0"
	}
	return v.Text(10)
}

func ParseInt(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %w", value, err)
	}
	return i, nil
}
