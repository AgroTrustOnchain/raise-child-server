package onchain

import (
	"fmt"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
	"time"
)

type AddChildArguments struct {
	Center       string
	IdentityCode string
	FirstName    string
	LastName     string
	Gender       string
	DateOfBirth  string
	Region       string
	AvatarBlobId string
}

type CreateCenterArguments struct {
	CapID       string
	Region      string
	Address     string
	PhoneNumber string
	ImageBlobID string
}

type SupportChildBooksNeedArguments struct {
	NeedID      string
	LocalPool   string
	ChildID     string
	DonorNft    string
	Amount      int64
	FirstName   string
	LastName    string
	Gender      string
	PhoneNumber string
	Email       string
	Message     string
}

type SupportChildMealNeedArguments struct {
	SupportChildBooksNeedArguments
	StartPeriod string
	EndPeriod   string
}

type SupportChildSpeicalNeedArguments struct {
	CampaignID  string
	LocalPool   string
	ChildID     string
	DonorNft    string
	Amount      int64
	FirstName   string
	LastName    string
	Gender      string
	PhoneNumber string
	Email       string
	Message     string
}

type CreateChildNormalNeedWithdrawProposalArguments struct {
	NeedID      string
	ChildID     string
	LocalPool   string
	Description string
	ClosedAt    int64
}

type CreateChildSpecialNeedProposalArguments struct {
	ChildID     string
	LocalPool   string
	Target      int64
	Description string
	ClosedAt    int64
}

type CreateChildSpecialNeedWithdrawProposalArguments struct {
	CampaignID     string
	LocalPool      string
	ChildID        string
	WithdrawAmount int64
	Description    string
	ClosedAt       int64
}

type ConfirmChildSpecialNeedProposalArguments struct {
	ProposalID string
	ChildID    string
}

type WithdrawFromNeedArguments struct {
	LocalPool  string
	TargetID   string
	ProposalID string
}

type IModuleChild interface {
	GetModule() string
	GetChildObjectStruct() string
	ToAddChildArguments(args AddChildArguments) []interface{}
	ToCreateCenterArguments(args CreateCenterArguments) []interface{}
	ToSupportChildBooksNeedArguments(args SupportChildBooksNeedArguments) []interface{}
	ToSupportChildMealNeedArguments(args SupportChildMealNeedArguments) []interface{}
	ToSupportChildSpeicalNeedArguments(args SupportChildSpeicalNeedArguments) []interface{}
	ToCreateChildNormalNeedWithdrawProposalArguments(args CreateChildNormalNeedWithdrawProposalArguments) []interface{}
	ToCreateChildSpecialNeedWithdrawProposalArguments(args CreateChildSpecialNeedWithdrawProposalArguments) []interface{}
	ToCreateChildSpecialNeedProposalArguments(args CreateChildSpecialNeedProposalArguments) []interface{}
	ToConfirmChildSpecialNeedProposalArguments(args ConfirmChildSpecialNeedProposalArguments) []interface{}
	ToWithdrawFromNeedArguments(args WithdrawFromNeedArguments) []interface{}
	GetFunctionAddChild() string
	GetFunctionUploadCenter() string
	GetFunctionAddStringMetadata() string
	GetFunctionAddNumberMetadata() string
	GetFunctionUpdateStringMetadata() string
	GetFunctionUpdateNumberMetadata() string
	GetFunctionRemoveStringMetadata() string
	GetFunctionRemoveNumberMetadata() string
	GetFunctionCreateChildBooksNeedWithdrawProposal() string
	GetFunctionCreateChildMealNeedWithdrawProposal() string
	GetFunctionCreateChildSpecialNeedWithdrawProposal() string
	GetFunctionCreateChildSpecialNeedProposal() string
	GetFunctionConfirmChildSpecialNeedProposal() string
	GetFunctionWithdrawFromBooksNeedProposal() string
	GetFunctionWithdrawFromMealNeedProposal() string
	GetFunctionWithdrawFromSpecialNeedCampaign() string
	GetFunctionSupportChildBooksNeed() string
	GetFunctionSupportChildMealNeed() string
	GetFunctionSupportChildSpecialNeedCampaign() string
}

type moduleChild struct{}

func InitializeModuleChild() IModuleChild {
	return &moduleChild{}
}

// GetFunctionConfirmChildSpecialNeedProposal implements IModuleChild.
func (m *moduleChild) GetFunctionConfirmChildSpecialNeedProposal() string {
	return sui.CONFIRM_CHILD_SPECIAL_NEED_PROPOSAL_FUNCTION
}

// GetFunctionCreateChildBooksNeedWithdrawProposal implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildBooksNeedWithdrawProposal() string {
	return sui.CREATE_CHILD_BOOKS_NEED_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionCreateChildMealNeedWithdrawProposal implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildMealNeedWithdrawProposal() string {
	return sui.CREATE_CHILD_MEAL_NEED_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionCreateChildSpecialNeedProposal implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildSpecialNeedProposal() string {
	return sui.CREATE_CHILD_SPECIAL_NEED_PROPOSAL_FUNCTION
}

// GetFunctionRemoveNumberMetadata implements IModuleChild.
func (m *moduleChild) GetFunctionRemoveNumberMetadata() string {
	panic("unimplemented")
}

// GetFunctionRemoveStringMetadata implements IModuleChild.
func (m *moduleChild) GetFunctionRemoveStringMetadata() string {
	panic("unimplemented")
}

// GetFunctionSupportChildBooksNeed implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildBooksNeed() string {
	return sui.SUPPORT_CHILD_BOOKS_NEED_FUNCTION
}

// GetFunctionSupportChildMealNeed implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildMealNeed() string {
	return sui.SUPPORT_CHILD_MEAL_NEED_FUNCTION
}

// GetFunctionSupportChildSpecialNeedCampaign implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildSpecialNeedCampaign() string {
	return sui.SUPPORT_CHILD_SPECIAL_NEED_CAMPAIGN
}

// GetFunctionWithdrawFromBooksNeedProposal implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromBooksNeedProposal() string {
	return sui.WITHDRAW_FROM_BOOKS_NEED_PROPOSAL_FUNCTION
}

// GetFunctionWithdrawFromMealNeedProposal implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromMealNeedProposal() string {
	return sui.WITHDRAW_FROM_MEAL_NEED_PROPOSAL_FUNCTION
}

// GetFunctionWithdrawFromSpecialNeedCampaign implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromSpecialNeedCampaign() string {
	return sui.WITHDRAW_FROM_SPECIAL_NEED_CAMPAIGN_FUNCTION
}

// GetFunctionCreateChildSpecialNeedWithdrawProposal implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildSpecialNeedWithdrawProposal() string {
	return sui.CREATE_CHILD_SPEICAL_NEED_WITHDRAW_PROPOSAL_FUNCTION
}

// ToCreateChildSpecialNeedWithdrawProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedWithdrawProposalArguments(args CreateChildSpecialNeedWithdrawProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.CampaignID,
		args.ChildID,
		args.WithdrawAmount,
		args.Description,
		args.ClosedAt,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildSpeicalNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildSpeicalNeedArguments(args SupportChildSpeicalNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.CampaignID,
		args.ChildID,
		args.LocalPool,
		args.DonorNft,
		args.Amount,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToConfirmChildSpecialNeedProposalArguments implements IModuleChild.
func (m *moduleChild) ToConfirmChildSpecialNeedProposalArguments(args ConfirmChildSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(""), // SpecialNeedDao
		args.ProposalID,
		args.ChildID,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildNormalNeedWithdrawProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildNormalNeedWithdrawProposalArguments(args CreateChildNormalNeedWithdrawProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.NeedID,
		args.ChildID,
		args.Description,
		args.ClosedAt,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildSpecialNeedProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedProposalArguments(args CreateChildSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.ChildID,
		args.LocalPool,
		args.Target,
		args.Description,
		args.ClosedAt,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildBooksNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildBooksNeedArguments(args SupportChildBooksNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.NeedID,
		args.LocalPool,
		args.ChildID,
		args.DonorNft,
		args.Amount,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildMealNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildMealNeedArguments(args SupportChildMealNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.NeedID,
		args.LocalPool,
		args.ChildID,
		args.DonorNft,
		args.Amount,
		args.StartPeriod,
		args.EndPeriod,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToWithdrawFromNeedArguments implements IModuleChild.
func (m *moduleChild) ToWithdrawFromNeedArguments(args WithdrawFromNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.TargetID,
		args.ProposalID,
		os.Getenv(""),
		sui.CLOCK_OBJECT_ID,
	}
}

// ToAddChildArguments implements IModuleChild.
func (m *moduleChild) ToAddChildArguments(args AddChildArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.Center,
		args.IdentityCode,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.DateOfBirth,
		args.Region,
		args.AvatarBlobId,
		fmt.Sprint(time.Now().Year()),
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
