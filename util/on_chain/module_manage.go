package onchain

import "raise-child/constants/on-chain/sui"

type IModuleManage interface {
	GetModule() string
	GetManageObjectStruct() string
	GetFunctionDonateSuiPool() string
	GetFunctionWithdrawSuiPool() string
}

type moduleManage struct{}

func InitializeModuleManage() IModuleManage {
	return &moduleManage{}
}

// GetManageObjectStruct implements IModuleManage.
func (m *moduleManage) GetManageObjectStruct() string {
	panic("unimplemented")
}

// GetFunctionDonateSuiPool implements IModuleManage.
func (m *moduleManage) GetFunctionDonateSuiPool() string {
	return sui.DONATE_SUI_POOL_FUNCTION
}

// GetFunctionWithdrawSuiPool implements IModuleManage.
func (m *moduleManage) GetFunctionWithdrawSuiPool() string {
	return sui.WITHDRAW_FROM_SUI_POOL_FUNCTION
}

// GetModule implements IModuleManage.
func (m *moduleManage) GetModule() string {
	return sui.MODULE_MANAGE
}
