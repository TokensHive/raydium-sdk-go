package module

import "github.com/shopspring/decimal"

type FormatOptions struct {
	GroupSeparator string
}

func FormatDecimal(value decimal.Decimal, _ FormatOptions) string {
	return value.String()
}
