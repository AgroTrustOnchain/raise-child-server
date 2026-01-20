package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type UpdatePublisherNftArguements struct {
	IdentityCode       string
	IdentityCardBlobID string
	AvatarBlobID       string
	FirstName          string
	LastName           string
	Gender             string
	DateOfBirth        string
	PhoneNumber        string
	Email              string
}

type IModuleManage interface {
	GetModule() string
	ToUpdatePublisherNftArguements(args UpdatePublisherNftArguements) []interface{}
	GetManageObjectStruct() string
	GetFunctionDonateSuiPool() string
	GetFunctionWithdrawSuiPool() string
	GetFunctionUpdatePublisherNft() string
}

type moduleManage struct{}

func InitializeModuleManage() IModuleManage {
	return &moduleManage{}
}

// GetManageObjectStruct implements IModuleManage.
func (m *moduleManage) GetManageObjectStruct() string {
	panic("unimplemented")
}

// GetFunctionUpdatePublisherNft implements IModuleManage.
func (m *moduleManage) GetFunctionUpdatePublisherNft() string {
	return sui.UPDATE_PUBLISHER_NFT_FUNCTION
}

// ToUpdatePublisherNftArguements implements IModuleManage.
func (m *moduleManage) ToUpdatePublisherNftArguements(args UpdatePublisherNftArguements) []interface{} {
	return []interface{}{
		os.Getenv(env.UPDATE_ADMIN_INFO_CAP_ID),
		os.Getenv(env.PUBLISHER_NFT_ID),
		args.IdentityCode,
		args.IdentityCardBlobID,
		args.AvatarBlobID,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.PhoneNumber,
		args.Email,
		sui.CLOCK_OBJECT_ID,
	}
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
