package module

import (
	"math/big"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"

	"github.com/TokensHive/raydium-sdk-go/common"
)

func mustToken(t *testing.T, symbol string, decimals int, mint string) Token {
	t.Helper()
	if mint == "" {
		mint = common.WSOLMint.String()
	}
	token, err := NewToken(TokenProps{
		Mint:     mint,
		Decimals: decimals,
		Symbol:   symbol,
		Name:     symbol,
	})
	if err != nil {
		t.Fatalf("unexpected token error: %v", err)
	}
	return token
}

func TestCurrencyAndToken(t *testing.T) {
	currency := NewCurrency(6, "", "")
	if currency.Symbol != "UNKNOWN" || currency.Name != "UNKNOWN" {
		t.Fatalf("default currency fields mismatch")
	}
	if !currency.Equals(currency) {
		t.Fatalf("currency should equal itself")
	}
	if !CurrencyEquals(currency, currency) {
		t.Fatalf("currency equals helper mismatch")
	}
	token := mustToken(t, "WSOL", 9, common.WSOLMint.String())
	if CurrencyEquals(currency, token) {
		t.Fatalf("mixed currency/token should not be equal")
	}
	if !CurrencyEquals(token, token) {
		t.Fatalf("token equality helper mismatch")
	}

	if !token.Equals(token) {
		t.Fatalf("token equals mismatch")
	}
	if token.AsCurrency().Decimals != token.Decimals {
		t.Fatalf("token currency conversion mismatch")
	}
	wsol, err := NewToken(TokenProps{Mint: common.SOLMint.String(), Decimals: 9})
	if err != nil || !wsol.Mint.Equals(WSOL.Mint) {
		t.Fatalf("sol should map to wsol")
	}
	if _, err := NewToken(TokenProps{Mint: "bad", Decimals: 9}); err == nil {
		t.Fatalf("expected invalid token mint error")
	}
	skipMint, err := NewToken(TokenProps{Mint: "bad", Decimals: 9, SkipMint: true})
	if err != nil || skipMint.Mint != (solana.PublicKey{}) {
		t.Fatalf("skip mint mismatch")
	}
	if MustToken(TokenProps{Mint: common.WSOLMint.String(), Decimals: 9}).Decimals != 9 {
		t.Fatalf("must token success mismatch")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Fatalf("must token should panic for invalid mint")
			}
		}()
		_ = MustToken(TokenProps{Mint: "bad", Decimals: 9})
	}()
}

func TestFractionOps(t *testing.T) {
	f := MustFraction(10, 2)
	if f.Quotient().Int64() != 5 {
		t.Fatalf("quotient mismatch")
	}
	inverted := f.Invert()
	if inverted.Numerator.Int64() != 2 || inverted.Denominator.Int64() != 10 {
		t.Fatalf("invert mismatch")
	}
	if _, err := NewFraction(1, 0); err == nil {
		t.Fatalf("expected zero denominator error")
	}
	if _, err := NewFraction(1, "bad"); err == nil {
		t.Fatalf("expected denominator parse error")
	}
	if _, err := f.Add(MustFraction(1, 2)); err != nil {
		t.Fatalf("add fraction error: %v", err)
	}
	if _, err := f.Add(MustFraction(1, 3)); err != nil {
		t.Fatalf("add fraction diff denominator error: %v", err)
	}
	if _, err := f.Add(struct{}{}); err == nil {
		t.Fatalf("expected add invalid type error")
	}
	if _, err := f.Sub(MustFraction(1, 2)); err != nil {
		t.Fatalf("sub fraction error: %v", err)
	}
	if _, err := f.Sub(MustFraction(1, 3)); err != nil {
		t.Fatalf("sub fraction diff denominator error: %v", err)
	}
	if _, err := f.Sub(struct{}{}); err == nil {
		t.Fatalf("expected sub invalid type error")
	}
	if _, err := f.Mul(MustFraction(2, 1)); err != nil {
		t.Fatalf("mul fraction error: %v", err)
	}
	if _, err := f.Mul(&Fraction{Numerator: big.NewInt(1), Denominator: big.NewInt(2)}); err != nil {
		t.Fatalf("mul pointer fraction error: %v", err)
	}
	if _, err := f.Mul(struct{}{}); err == nil {
		t.Fatalf("expected mul invalid type error")
	}
	if _, err := f.Div(MustFraction(2, 1)); err != nil {
		t.Fatalf("div fraction error: %v", err)
	}
	if _, err := f.Div(0); err == nil {
		t.Fatalf("expected divide by zero error")
	}
	if _, err := f.Div(struct{}{}); err == nil {
		t.Fatalf("expected divide invalid type error")
	}
	if _, err := f.ToSignificant(0, common.RoundDown); err == nil {
		t.Fatalf("expected significant digits error")
	}
	if _, err := f.ToFixed(-1, common.RoundDown); err == nil {
		t.Fatalf("expected fixed decimal error")
	}
	if _, err := f.ToSignificant(4, common.RoundHalfUp); err != nil {
		t.Fatalf("to significant failed: %v", err)
	}
	if _, err := f.ToFixed(2, common.RoundHalfUp); err != nil {
		t.Fatalf("to fixed failed: %v", err)
	}
	if _, err := f.ToFixed(2, common.RoundUp); err != nil {
		t.Fatalf("to fixed round up failed: %v", err)
	}
	zeroDenominator := Fraction{Numerator: big.NewInt(1), Denominator: big.NewInt(0)}
	if _, err := zeroDenominator.ToSignificant(2, common.RoundDown); err == nil {
		t.Fatalf("expected zero denominator significant error")
	}
	if _, err := zeroDenominator.ToFixed(2, common.RoundDown); err == nil {
		t.Fatalf("expected zero denominator fixed error")
	}
	if _, err := NewFraction("bad", 1); err == nil {
		t.Fatalf("expected new fraction parse error")
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Fatalf("must fraction should panic")
			}
		}()
		_ = MustFraction("bad")
	}()
	if !MustFraction(0).IsZero() {
		t.Fatalf("zero check failed")
	}
}

func TestAmountHelpers(t *testing.T) {
	token := mustToken(t, "RAY", 6, common.WSOLMint.String())
	intPart, fracPart, err := SplitNumber("1.23", 6)
	if err != nil || intPart != "1" || fracPart != "230000" {
		t.Fatalf("split number mismatch")
	}
	if _, _, err := SplitNumber("1.2.3", 6); err == nil {
		t.Fatalf("expected split error")
	}
	if integral, fractional, err := SplitNumber("2", 0); err != nil || integral != "2" || fractional != "0" {
		t.Fatalf("split number zero decimals mismatch")
	}
	if integral, fractional, err := SplitNumber("2", 6); err != nil || integral != "2" || fractional != "0" {
		t.Fatalf("split number integer with decimals mismatch")
	}
	if integral, fractional, err := SplitNumber("1.234", 2); err != nil || integral != "1" || fractional != "23" {
		t.Fatalf("split number truncate mismatch")
	}
	rawAmount, err := NewTokenAmount(token, 1000, true)
	if err != nil {
		t.Fatalf("new raw amount failed: %v", err)
	}
	if _, err := NewTokenAmount(token, struct{}{}, true); err == nil {
		t.Fatalf("expected invalid raw token amount")
	}
	humanAmount, err := NewTokenAmount(token, "1.25", false)
	if err != nil {
		t.Fatalf("new human amount failed: %v", err)
	}
	if _, err := NewTokenAmount(token, "bad", false); err == nil {
		t.Fatalf("expected invalid token amount input")
	}
	if _, err := NewTokenAmount(token, "1.2.3", false); err == nil {
		t.Fatalf("expected invalid split token amount input")
	}
	if _, err := NewTokenAmount(token, ".5", false); err == nil {
		t.Fatalf("expected invalid integral token amount input")
	}
	if _, err := NewTokenAmount(token, "1.a", false); err == nil {
		t.Fatalf("expected invalid fractional token amount input")
	}
	if rawAmount.IsZero() {
		t.Fatalf("raw amount should not be zero")
	}
	if humanAmount.Raw().Cmp(big.NewInt(1250000)) != 0 {
		t.Fatalf("unexpected raw amount")
	}
	if err := humanAmount.GT(rawAmount); err != nil {
		t.Fatalf("gt should pass")
	}
	if err := rawAmount.LT(humanAmount); err != nil {
		t.Fatalf("lt should pass")
	}
	if _, err := rawAmount.Add(rawAmount); err != nil {
		t.Fatalf("add amount failed")
	}
	if _, err := humanAmount.Subtract(rawAmount); err != nil {
		t.Fatalf("sub amount failed")
	}
	if _, err := humanAmount.ToSignificant(0, common.RoundDown); err != nil {
		t.Fatalf("to significant amount failed")
	}
	if _, err := humanAmount.ToFixed(0, common.RoundDown); err != nil {
		t.Fatalf("to fixed amount failed")
	}
	if humanAmount.ToExact() == "" {
		t.Fatalf("to exact should not be empty")
	}

	otherToken := mustToken(t, "USD", 6, common.USDCMint.String())
	otherAmount, _ := NewTokenAmount(otherToken, 1, true)
	if _, err := rawAmount.Add(otherAmount); err == nil {
		t.Fatalf("expected token mismatch")
	}
	if _, err := rawAmount.Subtract(otherAmount); err == nil {
		t.Fatalf("expected subtract token mismatch")
	}
	if _, err := humanAmount.ToFixed(7, common.RoundDown); err == nil {
		t.Fatalf("expected decimals overflow")
	}
	if err := rawAmount.GT(humanAmount); err == nil {
		t.Fatalf("expected not greater than error")
	}
	if err := rawAmount.GT(otherAmount); err == nil {
		t.Fatalf("expected gt token mismatch")
	}
	if err := humanAmount.LT(rawAmount); err == nil {
		t.Fatalf("expected not less than error")
	}
	if err := humanAmount.LT(otherAmount); err == nil {
		t.Fatalf("expected lt token mismatch")
	}
}

func TestCurrencyAmountAndPercentAndPrice(t *testing.T) {
	currency := NewCurrency(6, "USD", "USD")
	if _, err := NewCurrencyAmount(currency, struct{}{}, false); err == nil {
		t.Fatalf("expected invalid currency amount input")
	}
	currencyAmount, err := NewCurrencyAmount(currency, "2.5", false)
	if err != nil {
		t.Fatalf("new currency amount failed: %v", err)
	}
	other, _ := NewCurrencyAmount(currency, "1.0", false)
	if err := currencyAmount.GT(other); err != nil {
		t.Fatalf("currency gt failed")
	}
	if err := other.LT(currencyAmount); err != nil {
		t.Fatalf("currency lt failed")
	}
	if _, err := currencyAmount.Add(other); err != nil {
		t.Fatalf("currency add failed")
	}
	if _, err := currencyAmount.Sub(other); err != nil {
		t.Fatalf("currency sub failed")
	}
	if currencyAmount.IsZero() {
		t.Fatalf("currency amount should not be zero")
	}
	if _, err := currencyAmount.ToSignificant(0, common.RoundDown); err != nil {
		t.Fatalf("currency to significant failed")
	}
	if _, err := currencyAmount.ToFixed(0, common.RoundDown); err != nil {
		t.Fatalf("currency to fixed failed")
	}
	if currencyAmount.ToExact() == "" {
		t.Fatalf("currency exact should not be empty")
	}
	otherCurrency := NewCurrency(9, "BTC", "BTC")
	otherCurrencyAmount, _ := NewCurrencyAmount(otherCurrency, 1, true)
	if _, err := currencyAmount.Add(otherCurrencyAmount); err == nil {
		t.Fatalf("expected currency mismatch error")
	}
	if _, err := currencyAmount.Sub(otherCurrencyAmount); err == nil {
		t.Fatalf("expected currency sub mismatch error")
	}
	if _, err := currencyAmount.ToFixed(7, common.RoundDown); err == nil {
		t.Fatalf("expected currency decimals overflow")
	}
	if err := other.GT(currencyAmount); err == nil {
		t.Fatalf("expected not greater than")
	}
	if err := other.GT(otherCurrencyAmount); err == nil {
		t.Fatalf("expected currency gt mismatch")
	}
	if err := currencyAmount.LT(other); err == nil {
		t.Fatalf("expected not less than")
	}
	if err := currencyAmount.LT(otherCurrencyAmount); err == nil {
		t.Fatalf("expected currency lt mismatch")
	}

	percent, err := NewPercent(1, 4)
	if err != nil {
		t.Fatalf("new percent failed: %v", err)
	}
	if _, err := percent.ToSignificant(4, common.RoundHalfUp); err != nil {
		t.Fatalf("percent significant failed")
	}
	if _, err := percent.ToFixed(2, common.RoundHalfUp); err != nil {
		t.Fatalf("percent fixed failed")
	}
	if _, err := NewPercent("bad", 1); err == nil {
		t.Fatalf("expected invalid percent error")
	}

	base := mustToken(t, "A", 6, common.WSOLMint.String())
	quote := mustToken(t, "B", 6, common.USDCMint.String())
	price, err := NewPrice(PriceProps{
		BaseToken:   base,
		QuoteToken:  quote,
		Numerator:   2,
		Denominator: 1,
	})
	if err != nil {
		t.Fatalf("new price failed: %v", err)
	}
	if price.Raw().Numerator.Int64() != 2 {
		t.Fatalf("price raw mismatch")
	}
	if price.Adjusted().Numerator.Sign() == 0 {
		t.Fatalf("price adjusted should not be zero")
	}
	if _, err := price.Invert(); err != nil {
		t.Fatalf("price invert failed")
	}
	if _, err := price.ToSignificant(0, common.RoundHalfUp); err != nil {
		t.Fatalf("price significant failed")
	}
	if _, err := price.ToFixed(0, common.RoundHalfUp); err != nil {
		t.Fatalf("price fixed failed")
	}
	if _, err := NewPrice(PriceProps{BaseToken: base, QuoteToken: quote, Numerator: "bad", Denominator: 1}); err == nil {
		t.Fatalf("expected invalid price parse error")
	}
	next, _ := NewPrice(PriceProps{
		BaseToken:   quote,
		QuoteToken:  base,
		Numerator:   1,
		Denominator: 2,
	})
	if _, err := price.Mul(next); err != nil {
		t.Fatalf("price mul failed")
	}
	bad, _ := NewPrice(PriceProps{
		BaseToken:   base,
		QuoteToken:  base,
		Numerator:   1,
		Denominator: 1,
	})
	if _, err := price.Mul(bad); err == nil {
		t.Fatalf("expected price token mismatch")
	}

	if FormatDecimal(decimalMust(t, "12.34"), FormatOptions{}) != "12.34" {
		t.Fatalf("format decimal mismatch")
	}
}

func decimalMust(t *testing.T, value string) decimal.Decimal {
	t.Helper()
	d, err := decimal.NewFromString(value)
	if err != nil {
		t.Fatalf("failed parsing decimal: %v", err)
	}
	return d
}
