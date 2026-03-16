package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type UpdatePublisherNftArguments struct {
	AdminCap           string
	AdminNft           string
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

type MintCapArguments struct {
	Recipient string
}

type IModuleManage interface {
	GetModule() string
	ToUpdatePublisherNftArguments(args UpdatePublisherNftArguments) []interface{}
	ToMintCapArguments(args MintCapArguments) []interface{}
	GetAdminNftStruct() string
	GetManageObjectStruct() string
	GetAdminCapStruct() string
	GetUpdateAdminInfoCapStruct() string
	GetRegisterVolunteerCapStruct() string
	GetRegisterLeaderCapStruct() string
	GetRegisterAdminCapStruct() string
	GetUploadCenterCapStruct() string
	GetFunctionDonateSuiPool() string
	GetFunctionWithdrawSuiPool() string
	GetFunctionUpdatePublisherNft() string
	GetFunctionMintRegisterVolunteerCap() string
	GetFunctionMintRegisterLeaderCap() string
	GetFunctionMintRegisterAdminCap() string
	GetFunctionMintUploadCenterCap() string
}

type moduleManage struct{}

func InitializeModuleManage() IModuleManage {
	return &moduleManage{}
}

// GetAdminNftStruct implements IModuleManage.
func (m *moduleManage) GetAdminNftStruct() string {
	return sui.ADMIN_NFT_STRUCT
}

// GetManageObjectStruct implements IModuleManage.
func (m *moduleManage) GetManageObjectStruct() string {
	panic("unimplemented")
}

// GetFunctionMintRegisterAdminCap implements IModuleManage.
func (m *moduleManage) GetFunctionMintRegisterAdminCap() string {
	return sui.MINT_REGISTER_ADMIN_CAP_FUNCTION
}

// GetFunctionMintRegisterLeaderCap implements IModuleManage.
func (m *moduleManage) GetFunctionMintRegisterLeaderCap() string {
	return sui.MINT_REGISTER_LEADER_CAP_FUNCTION
}

// GetFunctionMintRegisterVolunteerCap implements IModuleManage.
func (m *moduleManage) GetFunctionMintRegisterVolunteerCap() string {
	return sui.MINT_REGISTER_VOLUNTEER_CAP_FUNCTION
}

// GetUploadCenterCapStruct implements IModuleManage.
func (m *moduleManage) GetUploadCenterCapStruct() string {
	return sui.UPLOAD_CENTER_CAP_STRUCT
}

// GetFunctionMintUploadCenterCap implements IModuleManage.
func (m *moduleManage) GetFunctionMintUploadCenterCap() string {
	return sui.MINT_UPLOAD_CENTER_CAP_FUNCTION
}

// ToMintCapArguments implements IModuleManage.
func (m *moduleManage) ToMintCapArguments(args MintCapArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.Recipient,
	}
}

// GetAdminCapStruct implements IModuleManage.
func (m *moduleManage) GetAdminCapStruct() string {
	return sui.ADMIN_CAP_STRUCT
}

// GetRegisterAdminCapStruct implements ImoduleManage.
func (m *moduleManage) GetRegisterAdminCapStruct() string {
	return sui.REGISTER_ADMIN_CAP_STRUCT
}

// GetRegisterLeaderCapStruct implements ImoduleManage.
func (m *moduleManage) GetRegisterLeaderCapStruct() string {
	return sui.REGISTER_LEADER_CAP_STRUCT
}

// GetRegisterVolunteerCapStruct implements ImoduleManage.
func (m *moduleManage) GetRegisterVolunteerCapStruct() string {
	return sui.REGISTER_VOLUNTEER_CAP_STRUCT
}

// GetUpdateAdminInfoCapStruct implements ImoduleManage.
func (m *moduleManage) GetUpdateAdminInfoCapStruct() string {
	return sui.UPDATE_ADMIN_INFO_CAP_STRUCT
}

// GetFunctionUpdatePublisherNft implements IModuleManage.
func (m *moduleManage) GetFunctionUpdatePublisherNft() string {
	return sui.UPDATE_PUBLISHER_NFT_FUNCTION
}

// ToUpdatePublisherNftArguments implements IModuleManage.
func (m *moduleManage) ToUpdatePublisherNftArguments(args UpdatePublisherNftArguments) []interface{} {
	return []interface{}{
		args.AdminCap,
		args.AdminNft,
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
