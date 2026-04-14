package module

import (
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/TokensHive/raydium-sdk-go/common"
)

type Token struct {
	Symbol      string
	Name        string
	Decimals    int
	IsToken2022 bool
	Mint        solana.PublicKey
}

type TokenProps struct {
	Mint        any
	Decimals    int
	Symbol      string
	Name        string
	SkipMint    bool
	IsToken2022 bool
}

var WSOL = Token{
	Symbol:      common.TokenWSOLInfo.Symbol,
	Name:        common.TokenWSOLInfo.Name,
	Decimals:    common.TokenWSOLInfo.Decimals,
	IsToken2022: false,
	Mint:        solana.MustPublicKeyFromBase58(common.TokenWSOLInfo.Address),
}

func NewToken(props TokenProps) (Token, error) {
	var (
		mint solana.PublicKey
		err  error
	)
	if props.SkipMint {
		mint = solana.PublicKey{}
	} else {
		mint, err = common.ValidateAndParsePublicKey(props.Mint, false)
		if err != nil {
			return Token{}, err
		}
	}
	if !props.SkipMint && (mint.Equals(common.SOLMint) || fmt.Sprint(props.Mint) == common.SOLMint.String()) {
		return WSOL, nil
	}
	if props.Symbol == "" {
		props.Symbol = mint.String()[:6]
	}
	if props.Name == "" {
		props.Name = mint.String()[:6]
	}
	return Token{
		Symbol:      props.Symbol,
		Name:        props.Name,
		Decimals:    props.Decimals,
		IsToken2022: props.IsToken2022,
		Mint:        mint,
	}, nil
}

func MustToken(props TokenProps) Token {
	t, err := NewToken(props)
	if err != nil {
		panic(err)
	}
	return t
}

func (t Token) Equals(other Token) bool {
	if t == other {
		return true
	}
	return t.Mint.Equals(other.Mint)
}

func (t Token) AsCurrency() Currency {
	return Currency{
		Symbol:   t.Symbol,
		Name:     t.Name,
		Decimals: t.Decimals,
	}
}
