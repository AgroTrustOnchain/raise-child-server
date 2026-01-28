package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type AddChildArguments struct {
	IdentityCode string
	FirstName    string
	LastName     string
	Gender       string
	DateOfBirth  string
	AvatarBlobId string
}

type CreateCenterArguments struct {
	CapID       string
	Region      string
	Address     string
	PhoneNumber string
	ImageBlobID string
}

type IModuleChild interface {
	GetModule() string
	GetChildObjectStruct() string
	ToAddChildArguments(args AddChildArguments) []interface{}
	ToCreateCenterArguments(args CreateCenterArguments) []interface{}
	GetFunctionAddChild() string
	GetFunctionUploadCenter() string
	GetFunctionAddStringMetadata() string
	GetFunctionAddNumberMetadata() string
	GetFunctionUpdateStringMetadata() string
	GetFunctionUpdateNumberMetadata() string
}

type moduleChild struct{}

func InitializeModuleChild() IModuleChild {
	return &moduleChild{}
}

// ToAddChildArguments implements IModuleChild.
func (m *moduleChild) ToAddChildArguments(args AddChildArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.IdentityCode,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.AvatarBlobId,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateCenterArguments implements IModuleChild.
func (m *moduleChild) ToCreateCenterArguments(args CreateCenterArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.CapID,
		args.Region,
		args.Address,
		args.PhoneNumber,
		args.ImageBlobID,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionUploadCenter implements IModuleChild.
func (m *moduleChild) GetFunctionUploadCenter() string {
	return sui.CREATE_CENTER_FUNCTION
}

// GetFunctionUpdateNumberMetadata implements IModuleChild.
func (m *moduleChild) GetFunctionUpdateNumberMetadata() string {
	return sui.UPDATE_NUMBER_METADATA_FUNCTION
}

// GetFunctionUpdateStringMetadata implements IModuleChild.
func (m *moduleChild) GetFunctionUpdateStringMetadata() string {
	return sui.UPDATE_STRING_METADATA_FUNCTION
}

// GetChildObjectStruct implements IModuleChild.
func (m *moduleChild) GetChildObjectStruct() string {
	return sui.CHILD_STRUCT
}

// GetFunctionAddChild implements IModuleChild.
func (m *moduleChild) GetFunctionAddChild() string {
	return sui.ADD_CHILD_FUNCTION
}

// GetFunctionAddNumberMetadata implements IModuleChild.
func (m *moduleChild) GetFunctionAddNumberMetadata() string {
	return sui.ADD_NUMBER_METADATA_FUNCTION
}

// GetFunctionAddStringMetadata implements IModuleChild.
func (m *moduleChild) GetFunctionAddStringMetadata() string {
	return sui.ADD_STRING_METADATA_FUNCTION
}

// GetModule implements IModuleChild.
func (m *moduleChild) GetModule() string {
	return sui.MODULE_CHILD
}
