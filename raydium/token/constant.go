package token

import "github.com/gagliardetto/solana-go"

type TokenInfo struct {
	ChainID   int
	Address   string
	ProgramID string
	Decimals  int
	Symbol    string
	Name      string
	LogoURI   string
}

var (
	SOLInfo = TokenInfo{
		ChainID:   101,
		Address:   solana.PublicKey{}.String(),
		ProgramID: solana.TokenProgramID.String(),
		Decimals:  9,
		Symbol:    "SOL",
		Name:      "solana",
		LogoURI:   "https://img-v1.raydium.io/icon/So11111111111111111111111111111111111111112.png",
	}
	TokenWSOL = TokenInfo{
		ChainID:   101,
		Address:   "So11111111111111111111111111111111111111112",
		ProgramID: solana.TokenProgramID.String(),
		Decimals:  9,
		Symbol:    "WSOL",
		Name:      "Wrapped SOL",
		LogoURI:   "https://img-v1.raydium.io/icon/So11111111111111111111111111111111111111112.png",
	}
)
