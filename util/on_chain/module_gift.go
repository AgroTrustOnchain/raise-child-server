package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type CreateGiftArguments struct {
	SponsorID       string
	ChildID         string
	TrackingCode    string
	Carrier         string
	GiftImageBlobID string
	Category        string
	Amount          int64
	FirstName       string
	LastName        string
	Gender          string
	PhoneNumber     string
	Email           string
	Message         string
	Description     string
}

type CancelGiftArguments struct {
	GiftID       string
	CancelReason string
}

type ConfirmReceiveGiftArguments struct {
	GiftID      string
	ChildID     string
	StaffID     string
	ImageBlobID string
}

type IModuleGift interface {
	GetModule() string
	GetGiftObjectStruct() string
	ToCreateGiftArguments(args CreateGiftArguments) []interface{}
	ToCancelGiftArguments(args CancelGiftArguments) []interface{}
	ToConfirmReceiveGiftArguments(args ConfirmReceiveGiftArguments) []interface{}
	GetFunctionCreateGift() string
	GetFunctionCancelGift() string
	GetFunctionConfirmReceiveGift() string
}

type moduleGift struct{}

func InitializeModuleGift() IModuleGift {
	return &moduleGift{}
}

// GetFunctionCancelGift implements IModuleGift.
func (m *moduleGift) GetFunctionCancelGift() string {
	return sui.CANCEL_GIFT_FUNCTION
}

// GetFunctionConfirmReceiveGift implements IModuleGift.
func (m *moduleGift) GetFunctionConfirmReceiveGift() string {
	return sui.CONFIRM_RECIEVE_GIFT_FUNCTION
}

// GetFunctionCreateGift implements IModuleGift.
func (m *moduleGift) GetFunctionCreateGift() string {
	return sui.CREATE_GIFT_FUNCTION
}

// GetGiftObjectStruct implements IModuleGift.
func (m *moduleGift) GetGiftObjectStruct() string {
	return sui.GIFT_STRUCT
}

// GetModule implements IModuleGift.
func (m *moduleGift) GetModule() string {
	return sui.MODULE_GIFT
}

// ToCancelGiftArguments implements IModuleGift.
func (m *moduleGift) ToCancelGiftArguments(args CancelGiftArguments) []interface{} {
	return []interface{}{
		args.GiftID,
		args.CancelReason,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToConfirmReceiveGiftArguments implements IModuleGift.
func (m *moduleGift) ToConfirmReceiveGiftArguments(args ConfirmReceiveGiftArguments) []interface{} {
	return []interface{}{
		args.GiftID,
		args.ChildID,
		args.StaffID,
		args.ImageBlobID,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateGiftArguments implements IModuleGift.
func (m *moduleGift) ToCreateGiftArguments(args CreateGiftArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.SponsorID,
		args.ChildID,
		args.TrackingCode,
		args.Carrier,
		args.GiftImageBlobID,
		args.Category,
		args.Amount,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Description,
		sui.CLOCK_OBJECT_ID,
	}
}
