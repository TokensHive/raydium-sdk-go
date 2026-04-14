package common

import (
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
)

type numberStringer struct{}

func (numberStringer) String() string { return "42" }

type badStringer struct{}

func (badStringer) String() string { return "bad" }

func TestBigNumberHelpers(t *testing.T) {
	n := big.NewInt(12)
	cloned := CloneBigInt(n)
	if cloned.Cmp(n) != 0 {
		t.Fatalf("clone mismatch")
	}
	cloned.SetInt64(20)
	if n.Int64() != 12 {
		t.Fatalf("clone should not mutate source")
	}
	if CloneBigInt(nil).Int64() != 0 {
		t.Fatalf("nil clone should be zero")
	}
	if TenExponential(big.NewInt(3)).Int64() != 1000 {
		t.Fatalf("invalid ten exponential")
	}
}

func TestParseBigNumberish(t *testing.T) {
	cases := []any{
		big.NewInt(2),
		*big.NewInt(3),
		"4",
		int(5),
		int8(6),
		int16(7),
		int32(8),
		int64(9),
		uint(10),
		uint8(11),
		uint16(12),
		uint32(13),
		uint64(14),
		float32(15),
		float64(16),
	}
	for _, tc := range cases {
		if _, err := ParseBigNumberish(tc); err != nil {
			t.Fatalf("unexpected parse error for %T: %v", tc, err)
		}
	}
	if _, err := ParseBigNumberish("1.2"); err == nil {
		t.Fatalf("expected decimal error")
	}
	if _, err := ParseBigNumberish("abc"); err == nil {
		t.Fatalf("expected invalid integer string error")
	}
	if _, err := ParseBigNumberish(float64(1.1)); err == nil {
		t.Fatalf("expected float underflow error")
	}
	if _, err := ParseBigNumberish(struct{}{}); err == nil {
		t.Fatalf("expected invalid type error")
	}
	if parsed, err := ParseBigNumberish(numberStringer{}); err != nil || parsed.Int64() != 42 {
		t.Fatalf("expected fmt.Stringer parse")
	}
	if _, err := ParseBigNumberish(uint(math.MaxUint)); err == nil {
		t.Fatalf("expected uint overflow error")
	}
	if _, err := ParseBigNumberish(uint64(math.MaxUint64)); err != nil {
		t.Fatalf("uint64 should be supported: %v", err)
	}
	if _, err := ParseBigNumberish(float64(math.MaxFloat64)); err == nil {
		t.Fatalf("expected float overflow error")
	}
	if _, err := ParseBigNumberish(float32(1.5)); err == nil {
		t.Fatalf("expected float32 underflow error")
	}
	var nilBigInt *big.Int
	if parsed, err := ParseBigNumberish(nilBigInt); err != nil || parsed.Sign() != 0 {
		t.Fatalf("nil *big.Int should parse as zero")
	}
	if _, err := ParseBigNumberish(badStringer{}); err == nil {
		t.Fatalf("bad fmt.Stringer should return parse error")
	}
	if _, err := ParseBigNumberish(""); err == nil {
		t.Fatalf("expected empty string error")
	}
}

func TestConstantHelpers(t *testing.T) {
	sign, intPart, fracPart, err := SplitSignedDecimal("-12.45")
	if err != nil || sign != "-" || intPart != "12" || fracPart != "45" {
		t.Fatalf("split signed decimal mismatch")
	}
	if _, _, _, err := SplitSignedDecimal(""); err == nil {
		t.Fatalf("expected empty split error")
	}
	if _, _, _, err := SplitSignedDecimal("-"); err == nil {
		t.Fatalf("expected invalid signed decimal")
	}
	if _, intPart, fracPart, err := SplitSignedDecimal("42"); err != nil || intPart != "42" || fracPart != "" {
		t.Fatalf("split integer mismatch")
	}
	if _, _, err := DecimalStringToRational("abc"); err == nil {
		t.Fatalf("expected decimal rational error")
	}
	if _, _, err := DecimalStringToRational(""); err == nil {
		t.Fatalf("expected decimal rational split error")
	}
	num, den, err := DecimalStringToRational("12.34")
	if err != nil {
		t.Fatalf("unexpected rational parse error: %v", err)
	}
	if num.String() != "1234" || den.String() != "100" {
		t.Fatalf("unexpected rational values")
	}
	negNum, negDen, err := DecimalStringToRational("-.5")
	if err != nil || negNum.String() != "-5" || negDen.String() != "10" {
		t.Fatalf("unexpected negative rational values")
	}
	intNum, intDen, err := DecimalStringToRational("5")
	if err != nil || intNum.String() != "5" || intDen.String() != "1" {
		t.Fatalf("unexpected integer rational values")
	}
	if BigIntToString(nil) != "0" {
		t.Fatalf("nil bigint string mismatch")
	}
	if BigIntToString(big.NewInt(7)) != "7" {
		t.Fatalf("non-nil bigint string mismatch")
	}
	if _, err := ParseInt("x"); err == nil {
		t.Fatalf("expected atoi error")
	}
	if parsed, err := ParseInt("12"); err != nil || parsed != 12 {
		t.Fatalf("expected parse int success")
	}
	if MustParseIntBase10("1").Int64() != 1 {
		t.Fatalf("must parse mismatch")
	}
	defer func() {
		if recover() == nil {
			t.Fatalf("must parse should panic")
		}
	}()
	MustParseIntBase10("x")
}

func TestLoggerAndLevels(t *testing.T) {
	SetLoggerLevel("foo", LogDebug)
	logger := CreateLogger("foo")
	if logger.ModuleName() != "foo" {
		t.Fatalf("module name mismatch")
	}
	if logger.Time() == "" {
		t.Fatalf("expected timestamp")
	}
	logger.Debug("debug").Info("info").Warning("warn").Error("error")
	logger.SetLevel(LogError)
	logger.Debug("not-logged").Info("not-logged").Warning("not-logged")
	SetLoggerLevel("foo", LogInfo)
	_ = CreateLogger("foo")
	_ = CreateLogger("new-module-with-default-level")
	_ = NewLogger("bar", LogDebug)
	_ = fmt.Sprint(logger)
	defer func() {
		if recover() == nil {
			t.Fatalf("expected logger panic")
		}
	}()
	logger.LogWithError("panic")
}

func TestPubKeyHelpers(t *testing.T) {
	parsed, err := ValidateAndParsePublicKey(WSOLMint.String(), false)
	if err != nil || !parsed.Equals(WSOLMint) {
		t.Fatalf("parse wspl failed")
	}
	transformed, err := ValidateAndParsePublicKey(SOLMint.String(), true)
	if err != nil || !transformed.Equals(WSOLMint) {
		t.Fatalf("sol transform failed")
	}
	if _, err := ValidateAndParsePublicKey("bad", false); err == nil {
		t.Fatalf("expected invalid key error")
	}
	key := WSOLMint
	if _, err := ValidateAndParsePublicKey(key, false); err != nil {
		t.Fatalf("public key input should succeed")
	}
	if _, err := ValidateAndParsePublicKey(&key, false); err != nil {
		t.Fatalf("public key pointer input should succeed")
	}
	if _, err := ValidateAndParsePublicKey(key, true); err != nil {
		t.Fatalf("public key transform path should succeed")
	}
	sol := SOLMint
	solPtr := &sol
	if transformedPtr, err := ValidateAndParsePublicKey(solPtr, true); err != nil || !transformedPtr.Equals(WSOLMint) {
		t.Fatalf("sol public key pointer transform should map to wsol")
	}
	transformedKey, err := ValidateAndParsePublicKey(sol, true)
	if err != nil || !transformedKey.Equals(WSOLMint) {
		t.Fatalf("sol public key transform should map to wsol")
	}
	var nilKey *solana.PublicKey
	if _, err := ValidateAndParsePublicKey(nilKey, false); err == nil {
		t.Fatalf("nil pointer should fail")
	}
	if _, err := ValidateAndParsePublicKey(99, false); err == nil {
		t.Fatalf("invalid public key type should fail")
	}
	if TryParsePublicKey("bad") != "bad" {
		t.Fatalf("try parse should return input string on failure")
	}
	if _, ok := TryParsePublicKey(WSOLMint.String()).(solana.PublicKey); !ok {
		t.Fatalf("try parse should return public key")
	}
	if _, err := SolToWSol(SOLMint); err != nil {
		t.Fatalf("sol to wsol failed: %v", err)
	}
}

func TestProgramIDsAndConstants(t *testing.T) {
	if len(AllProgramID) == 0 || len(DevnetProgramID) == 0 || len(IDOAllProgram) == 0 {
		t.Fatalf("program ID maps should be populated")
	}
	if EmptyConnection == "" || EmptyOwner == "" {
		t.Fatalf("error constants must be set")
	}
	if SOLInfo.Address == "" || TokenWSOLInfo.Address == "" {
		t.Fatalf("token infos must be set")
	}
}

func TestSleep(t *testing.T) {
	start := time.Now()
	Sleep(0)
	if time.Since(start) < 0 {
		t.Fatalf("invalid duration")
	}
}
