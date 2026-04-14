package base

import "github.com/TokensHive/raydium-sdk-go/common"

type ModuleBase struct {
	ModuleName string
	Logger     *common.Logger
}

func NewModuleBase(moduleName string) ModuleBase {
	return ModuleBase{
		ModuleName: moduleName,
		Logger:     common.CreateLogger(moduleName),
	}
}
