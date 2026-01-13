package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type AddChildArguements struct {
	IdentityCode string
	FirstName    string
	LastName     string
	Gender       string
	DateOfBirth  string
	AvatarBlobId string
}

type IModuleChild interface {
	GetModule() string
	GetChildObjectStruct() string
	ToAddChildArguements(args AddChildArguements) []interface{}
	GetFunctionAddChild() string
	GetFunctionAddStringMetadata() string
	GetFunctionAddNumberMetadata() string
	GetFunctionUpdateStringMetadata() string
	GetFunctionUpdateNumberMetadata() string
}

type moduleChild struct{}

func InitializeModuleChild() IModuleChild {
	return &moduleChild{}
}

// ToAddChildArguements implements IModuleChild.
func (m *moduleChild) ToAddChildArguements(args AddChildArguements) []interface{} {
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
