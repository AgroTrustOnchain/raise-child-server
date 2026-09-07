package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type RegisterStaffArguments struct {
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

type RegisterAdminArguments struct {
	IdentityCode       string
	IdentityCardBlobID string
	AvatarBlobID       string
	FirstName          string
	LastName           string
	Gender             string
	DateOfBirth        string
	PhoneNumber        string
	Email              string
	Owner              string
	Sender             string
}

type RegisterVolunteerArguments struct {
	Region string
	RegisterAdminArguments
}

type RegisterLeaderArguments struct {
	CenterAddress     string
	CenterPhoneNumber string
	CenterImageBlobID string
	RegisterVolunteerArguments
}

type RegisterNormalStaffArguments struct {
	LocalPoolID string
	Region      string
	RegisterAdminArguments
}

type IModuleStaff interface {
	GetModule() string
	ToRegisterStaffArguments(args RegisterStaffArguments) []interface{}
	ToRegisterAdminArguments(args RegisterAdminArguments) []interface{}
	ToRegisterVolunteerArguments(args RegisterVolunteerArguments) []interface{}
	ToRegisterLeaderArguments(args RegisterLeaderArguments) []interface{}
	ToRegisterNormalStaffArguments(args RegisterNormalStaffArguments) []interface{} // Bỏ
	GetFunctionRegisterStaff() string                                               // Bỏ
	GetFunctionRegisterVolunteer() string
	GetFunctionRegisterLeader() string
	GetFunctionRegisterAdmin() string
	GetStaffNftObjectStruct() string
	GetStaffEventEmittedStruct() string
}

type moduleStaff struct{}

func InitializeModuleStaff() IModuleStaff {
	return &moduleStaff{}
}

// ToRegisterAdminArguments implements IModuleStaff.
func (m *moduleStaff) ToRegisterAdminArguments(args RegisterAdminArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.IdentityCode,
		args.AvatarBlobID,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.PhoneNumber,
		args.Email,
		args.Owner,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToRegisterLeaderArguments implements IModuleStaff.
func (m *moduleStaff) ToRegisterLeaderArguments(args RegisterLeaderArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.IdentityCode,
		args.IdentityCardBlobID,
		args.AvatarBlobID,
		args.Region,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.PhoneNumber,
		args.Email,
		args.CenterAddress,
		args.CenterPhoneNumber,
		args.CenterImageBlobID,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToRegisterVolunteerArguments implements IModuleStaff.
func (m *moduleStaff) ToRegisterVolunteerArguments(args RegisterVolunteerArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.IdentityCode,
		args.IdentityCardBlobID,
		args.AvatarBlobID,
		args.Region,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.PhoneNumber,
		args.Email,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToRegisterNormalStaffArguments implements IModuleStaff.
func (m *moduleStaff) ToRegisterNormalStaffArguments(args RegisterNormalStaffArguments) []interface{} {
	if args.LocalPoolID == "" {
		return []interface{}{
			os.Getenv(env.ADMIN_CAP_ID_1),
			os.Getenv(env.MANAGE_OBJECT_ID),
			args.IdentityCode,
			args.AvatarBlobID,
			args.Region,
			args.FirstName,
			args.LastName,
			args.Gender,
			args.DateOfBirth,
			args.PhoneNumber,
			args.Email,
			args.Owner,
			sui.CLOCK_OBJECT_ID,
		}
	}

	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.LocalPoolID,
		args.IdentityCode,
		args.AvatarBlobID,
		args.Region,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.PhoneNumber,
		args.Email,
		args.Owner,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToRegisterStaffArguments implements IModuleStaff.
func (m *moduleStaff) ToRegisterStaffArguments(args RegisterStaffArguments) []interface{} {
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

// GetFunctionRegisterVolunteer implements IModuleStaff.
func (m *moduleStaff) GetFunctionRegisterVolunteer() string {
	return sui.REGISTER_VOLUNTEER_FUNCTION
}

// GetFunctionRegisterLeader implements IModuleStaff.
func (m *moduleStaff) GetFunctionRegisterLeader() string {
	return sui.REGISTER_LOCAL_LEADER_FUNCTION
}

// GetFunctionRegisterAdmin implements IModuleStaff.
func (m *moduleStaff) GetFunctionRegisterAdmin() string {
	return sui.REGISTER_ADMIN_FUNCTION
}

// GetModule implements IModuleStaff.
func (m *moduleStaff) GetModule() string {
	return sui.MODULE_STAFF
}

// GetStaffNftObjectStruct implements IModuleStaff.
func (m *moduleStaff) GetStaffNftObjectStruct() string {
	return sui.STAFF_NFT_STRUCT
}

// GetStaffEventEmittedStruct implements IModuleStaff.
func (m *moduleStaff) GetStaffEventEmittedStruct() string {
	return sui.STAFF_CREATED_EVENT
}
