package common

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type PublicKeyish interface{}

func mustPublicKey(value string) solana.PublicKey {
	return solana.MustPublicKeyFromBase58(value)
}

var (
	MemoProgramID        = mustPublicKey("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
	MemoProgramID2       = mustPublicKey("MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr")
	RentProgramID        = mustPublicKey("SysvarRent111111111111111111111111111111111")
	ClockProgramID       = mustPublicKey("SysvarC1ock11111111111111111111111111111111")
	MetadataProgramID    = mustPublicKey("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	InstructionProgramID = mustPublicKey("Sysvar1nstructions1111111111111111111111111")
	SystemProgramID      = solana.SystemProgramID

	RAYMint  = mustPublicKey("4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R")
	PAIMint  = mustPublicKey("Ea5SjE2Y6yvCeW5dYTn7PYMuW5ikXkvbGdcmSnXeaLjS")
	SRMMint  = mustPublicKey("SRMuApVNdxXokk5GT7XD5cUUgXMBCoAz2LHeuAoKWRt")
	USDCMint = mustPublicKey("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	USDTMint = mustPublicKey("Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB")
	MSOLMint = mustPublicKey("mSoLzYCxHdYgdzU16g5QSh3i5K3z3KZK7ytfqcJm7So")
	StSOLMint = mustPublicKey("7dHbWXmci3dT8UFYWYZweBLXgycu7Y3iL6trKn1Y7ARj")
	USDHMint = mustPublicKey("USDH1SM1ojwWUga67PGrgFWUHibbjqMvuMaDkRJTgkX")
	NRVMint  = mustPublicKey("NRVwhjBQiUPYtfDT5zRBVJajzFQHaBUNtC7SNVvqRFa")
	ANAMint  = mustPublicKey("ANAxByE6G2WjFp7A4NqtWYXb3mgruyzZYg3spfxe6Lbo")
	ETHMint  = mustPublicKey("7vfCXTUXx5WJV5JADk17DUJ4ksgau7utNKj4b963voxs")
	WSOLMint = mustPublicKey("So11111111111111111111111111111111111111112")
	SOLMint  = solana.PublicKey{}
)

func ValidateAndParsePublicKey(publicKey PublicKeyish, transformSol bool) (solana.PublicKey, error) {
	switch key := publicKey.(type) {
	case solana.PublicKey:
		if transformSol && key.Equals(SOLMint) {
			return WSOLMint, nil
		}
		return key, nil
	case *solana.PublicKey:
		if key == nil {
			return solana.PublicKey{}, fmt.Errorf("invalid public key")
		}
		if transformSol && key.Equals(SOLMint) {
			return WSOLMint, nil
		}
		return *key, nil
	case string:
		return parsePubKeyString(key, transformSol)
	default:
		return solana.PublicKey{}, fmt.Errorf("invalid public key")
	}
}

func parsePubKeyString(v string, transformSol bool) (solana.PublicKey, error) {
	if transformSol && v == SOLMint.String() {
		return WSOLMint, nil
	}
	parsed, err := solana.PublicKeyFromBase58(v)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid public key")
	}
	return parsed, nil
}

func TryParsePublicKey(v string) any {
	key, err := solana.PublicKeyFromBase58(v)
	if err != nil {
		return v
	}
	return key
}

func SolToWSol(mint PublicKeyish) (solana.PublicKey, error) {
	return ValidateAndParsePublicKey(mint, true)
}
