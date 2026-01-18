package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type RegisterStaffArguements struct {
	IdentityCode       string
	IdentityCardBlobID string
	Role               string
	AvatarBlobID       string
	Region             string
	FirstName          string
	LastName           string
	Gender             string
	PhoneNumber        string
	Email              string
}

type IModuleStaff interface {
	GetModule() string
	ToRegisterStaffArguements(args RegisterStaffArguements) []interface{}
	GetFunctionRegisterStaff() string
	GetStaffObjectStruct() string
	GetStaffNftObjectStruct() string
}

type moduleStaff struct{}

func InitializeModuleStaff() IModuleStaff {
	return &moduleStaff{}
}

// ToRegisterStaffArguements implements IModuleStaff.
func (m *moduleStaff) ToRegisterStaffArguements(args RegisterStaffArguements) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.IdentityCode,
		args.IdentityCardBlobID,
		args.Role,
		args.AvatarBlobID,
		args.Region,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionRegisterStaff implements IModuleStaff.
func (m *moduleStaff) GetFunctionRegisterStaff() string {
	return sui.REGISTER_STAFF_FUNCTION
}

// GetModule implements IModuleStaff.
func (m *moduleStaff) GetModule() string {
	return sui.MODULE_STAFF
}

// GetStaffNftObjectStruct implements IModuleStaff.
func (m *moduleStaff) GetStaffNftObjectStruct() string {
	return sui.STAFF_NFT_STRUCT
}

// GetStaffObjectStruct implements IModuleStaff.
func (m *moduleStaff) GetStaffObjectStruct() string {
	return sui.STAFF_STRUCT
}
