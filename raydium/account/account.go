package account

import "github.com/TokensHive/raydium-sdk-go/raydium/base"

type TokenAccountData struct {
	Mint   string
	Amount string
}

type Module struct {
	base.ModuleBase
	TokenAccounts []TokenAccountData
}

func New(moduleName string, tokenAccounts []TokenAccountData) *Module {
	return &Module{
		ModuleBase:    base.NewModuleBase(moduleName),
		TokenAccounts: tokenAccounts,
	}
}

func (m *Module) ResetTokenAccounts() {
	m.TokenAccounts = []TokenAccountData{}
}
