package onchain

import (
	"fmt"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type CreateGiftForCenterArguments struct {
	DonorID         string
	CenterID        string
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
	Sender          string
}

type CreateGiftForChildArguments struct {
	ChildID string
	CreateGiftForCenterArguments
}
type CancelGiftArguments struct {
	GiftID       string
	CancelReason string
	Sender       string
}

type ConfirmReceiveGiftArguments struct {
	GiftID      string
	Recipient   string
	StaffID     string
	ImageBlobID string
	Sender      string
}

type IModuleGift interface {
	GetModule() string
	GetGiftObjectStruct() string
	ToCreateGiftForCenterArguments(args CreateGiftForCenterArguments) []interface{}
	ToCreateGiftForChildArguments(args CreateGiftForChildArguments) []interface{}
	ToCancelGiftArguments(args CancelGiftArguments) []interface{}
	ToConfirmReceiveGiftArguments(args ConfirmReceiveGiftArguments) []interface{}
	GetFunctionCreateGiftForChild() string
	GetFunctionCreateGiftForCenter() string
	GetFunctionConfirmReceiveChildGift() string
	GetFunctionConfirmReceiveCenterGift() string
	GetFunctionCancelGift() string
}

type moduleGift struct{}

func InitializeModuleGift() IModuleGift {
	return &moduleGift{}
}

// GetFunctionConfirmReceiveCenterGift implements IModuleGift.
func (m *moduleGift) GetFunctionConfirmReceiveCenterGift() string {
	return sui.CONFIRM_RECIEVE_CENTER_GIFT_FUNCTION
}

// GetFunctionCreateGiftForCenter implements IModuleGift.
func (m *moduleGift) GetFunctionCreateGiftForCenter() string {
	return sui.CREATE_GIFT_FOR_CENTER_FUNCTION
}

// GetFunctionCancelGift implements IModuleGift.
func (m *moduleGift) GetFunctionCancelGift() string {
	return sui.CANCEL_GIFT_FUNCTION
}

// GetFunctionConfirmReceiveChildGift implements IModuleGift.
func (m *moduleGift) GetFunctionConfirmReceiveChildGift() string {
	return sui.CONFIRM_RECIEVE_CHILD_GIFT_FUNCTION
}

// GetFunctionCreateGiftForChild implements IModuleGift.
func (m *moduleGift) GetFunctionCreateGiftForChild() string {
	return sui.CREATE_GIFT_FOR_CHILD_FUNCTION
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
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.GiftID,
		args.CancelReason,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToConfirmReceiveGiftArguments implements IModuleGift.
func (m *moduleGift) ToConfirmReceiveGiftArguments(args ConfirmReceiveGiftArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.GiftID,
		args.Recipient,
		args.StaffID,
		args.ImageBlobID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateGiftArguments implements IModuleGift.
func (m *moduleGift) ToCreateGiftForCenterArguments(args CreateGiftForCenterArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.DonorID,
		args.CenterID,
		args.TrackingCode,
		args.Carrier,
		args.GiftImageBlobID,
		args.Category,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Description,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateGiftForChildArguments implements IModuleGift.
func (m *moduleGift) ToCreateGiftForChildArguments(args CreateGiftForChildArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.DonorID,
		args.ChildID,
		args.CenterID,
		args.TrackingCode,
		args.Carrier,
		args.GiftImageBlobID,
		args.Category,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Description,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}
