package ido

import "github.com/TokensHive/raydium-sdk-go/raydium/base"

type Module struct {
	base.ModuleBase
}

func New(moduleName string) *Module {
	return &Module{ModuleBase: base.NewModuleBase(moduleName)}
}
