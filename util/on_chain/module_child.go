package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
	"time"
)

type AddChildArguments struct {
	Center                           string
	IdentityCode                     string
	FirstName                        string
	LastName                         string
	Gender                           string
	DateOfBirth                      string
	HomeAddress                      string
	Region                           string
	AvatarBlobId                     string
	HomeBlobID                       string
	FirstGuardianFullName            string
	FirstGuardianPhone               string
	FirstGuardianRelation            string
	FirstGuardianIdentityCardBlobID  string
	SecondGuardianFullName           string
	SecondGuardianPhone              string
	SecondGuardianRelation           string
	SecondGuardianIdentityCardBlobID string
}

type CreateCenterArguments struct {
	CapID       string
	Region      string
	Address     string
	PhoneNumber string
	ImageBlobID string
	Leaders     []string
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

type SupportChildHealthInsuranceNeedArguments struct {
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

type ConfirmProvideMealForChildArguments struct {
	ChildID     string
	NeedID      string
	StaffNft    string
	ImageBlobID string
	ProvideDate string
}

type ConfirmProvideMealForChildArgumentsV2 struct {
	ChildID     string
	NeedID      string
	StaffNft    string
	ImageBlobID string
	ProvideDate string
	Actor       string
}

type CreateChildNormalNeedWithdrawProposalArguments struct {
	NeedID      string
	ChildID     string
	LocalPool   string
	Description string
	ClosedAt    int64
}

type CreateChildNormalNeedWithdrawProposalArgumentsV2 struct {
	NeedID      string
	ChildID     string
	LocalPool   string
	Description string
	ProofBlobID *string
	ClosedAt    int64
	Creator     string
}

type CreateChildSpecialNeedProposalArguments struct {
	ChildID     string
	LocalPool   string
	Target      int64
	Description string
	ClosedAt    int64
}

type CreateChildSpecialNeedProposalArgumentsV2 struct {
	ChildID     string
	LocalPool   string
	Target      int64
	Description string
	ProofBlobID *string
	ClosedAt    int64
	Creator     string
}

type CreateChildSpecialNeedWithdrawProposalArguments struct {
	CampaignID     string
	LocalPool      string
	ChildID        string
	WithdrawAmount int64
	Description    string
	ClosedAt       int64
}

type CreateChildSpecialNeedWithdrawProposalArgumentsV2 struct {
	CampaignID     string
	LocalPool      string
	ChildID        string
	WithdrawAmount int64
	Description    string
	ProofBlobID    *string
	ClosedAt       int64
	Creator        string
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

type SubmitTaskArguments struct {
	Center      string
	StaffNft    string
	Description string
	ImageBlobID string
	Actor       string
}

type IModuleChild interface {
	GetModule() string
	GetChildObjectStruct() string
	ToAddChildArguments(args AddChildArguments) []interface{}
	ToCreateCenterArguments(args CreateCenterArguments) []interface{}
	ToSupportChildBooksNeedArguments(args SupportChildBooksNeedArguments) []interface{}
	ToSupportChildHealthInsuranceNeedArguments(args SupportChildHealthInsuranceNeedArguments) []interface{}
	ToSupportChildMealNeedArguments(args SupportChildMealNeedArguments) []interface{}
	ToSupportChildSpeicalNeedArguments(args SupportChildSpeicalNeedArguments) []interface{}
	ToConfirmProvideMealForChildArguments(args ConfirmProvideMealForChildArguments) []interface{}
	ToConfirmProvideMealForChildArgumentsV2(args ConfirmProvideMealForChildArgumentsV2) []interface{}
	ToCreateChildNormalNeedWithdrawProposalArguments(args CreateChildNormalNeedWithdrawProposalArguments) []interface{}
	ToCreateChildSpecialNeedWithdrawProposalArguments(args CreateChildSpecialNeedWithdrawProposalArguments) []interface{}
	ToCreateChildSpecialNeedProposalArguments(args CreateChildSpecialNeedProposalArguments) []interface{}
	ToCreateChildNormalNeedWithdrawProposalArgumentsV2(args CreateChildNormalNeedWithdrawProposalArgumentsV2) []interface{}
	ToCreateChildSpecialNeedWithdrawProposalArgumentsV2(args CreateChildSpecialNeedWithdrawProposalArgumentsV2) []interface{}
	ToCreateChildSpecialNeedProposalArgumentsV2(args CreateChildSpecialNeedProposalArgumentsV2) []interface{}
	ToConfirmChildSpecialNeedProposalArguments(args ConfirmChildSpecialNeedProposalArguments) []interface{}
	ToWithdrawFromNeedArguments(args WithdrawFromNeedArguments) []interface{}
	ToSubmitTaskArguments(args SubmitTaskArguments) []interface{}
	GetFunctionAddChild() string
	GetFunctionUploadCenter() string
	GetFunctionAddStringMetadata() string
	GetFunctionAddNumberMetadata() string
	GetFunctionUpdateStringMetadata() string
	GetFunctionUpdateNumberMetadata() string
	GetFunctionRemoveStringMetadata() string
	GetFunctionRemoveNumberMetadata() string
	GetFunctionCreateChildBooksNeedWithdrawProposal() string
	GetFunctionCreateChildBooksNeedWithdrawProposalV2() string
	GetFunctionCreateChildMealNeedWithdrawProposal() string
	GetFunctionCreateChildMealNeedWithdrawProposalV2() string
	GetFunctionCreateChildSpecialNeedWithdrawProposal() string
	GetFunctionCreateChildSpecialNeedWithdrawProposalV2() string
	GetFunctionCreateChildHealthInsuranceNeedWithdrawProposalV2() string
	GetFunctionCreateChildSpecialNeedProposal() string
	GetFunctionCreateChildSpecialNeedProposalV2() string
	GetFunctionConfirmChildSpecialNeedProposal() string
	GetFunctionWithdrawFromBooksNeedProposal() string
	GetFunctionWithdrawFromHealthInsuranceNeedProposal() string
	GetFunctionWithdrawFromMealNeedProposal() string
	GetFunctionWithdrawFromSpecialNeedCampaign() string
	GetFunctionSupportChildBooksNeed() string
	GetFunctionSupportChildHealthInsuranceNeed() string
	GetFunctionSupportChildMealNeed() string
	GetFunctionSupportChildSpecialNeedCampaign() string
	GetFunctionConfirmProvideMealForChild() string
	GetFunctionConfirmProvideMealForChildV2() string
	GetFunctionSubmitTask() string
}

type moduleChild struct{}

func InitializeModuleChild() IModuleChild {
	return &moduleChild{}
}

// ToCreateChildNormalNeedWithdrawProposalArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToCreateChildNormalNeedWithdrawProposalArgumentsV2(args CreateChildNormalNeedWithdrawProposalArgumentsV2) []interface{} {
	var proofBlobId string = ""
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.NeedID,
		args.ChildID,
		args.Description,
		proofBlobId,
		uint64(args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildSpecialNeedProposalArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedProposalArgumentsV2(args CreateChildSpecialNeedProposalArgumentsV2) []interface{} {
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.ChildID,
		args.LocalPool,
		uint64(args.Target),
		args.Description,
		proofBlobId,
		uint64(args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildSpecialNeedWithdrawProposalArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedWithdrawProposalArgumentsV2(args CreateChildSpecialNeedWithdrawProposalArgumentsV2) []interface{} {
	var proofBlobId string = ""
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.CampaignID,
		args.ChildID,
		uint64(args.WithdrawAmount),
		args.Description,
		proofBlobId,
		uint64(args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionSupportChildHealthInsuranceNeed implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildHealthInsuranceNeed() string {
	panic("unimplemented")
}

// ToSupportChildHealthInsuranceNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildHealthInsuranceNeedArguments(args SupportChildHealthInsuranceNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.NeedID,
		args.LocalPool,
		args.ChildID,
		args.DonorNft,
		uint64(args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSubmitTaskArguments implements IModuleChild.
func (m *moduleChild) ToSubmitTaskArguments(args SubmitTaskArguments) []interface{} {
	return []interface{}{
		args.Center,
		args.StaffNft,
		args.Description,
		args.ImageBlobID,
		args.Actor,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionConfirmProvideMealForChildV2 implements IModuleChild.
func (m *moduleChild) GetFunctionConfirmProvideMealForChildV2() string {
	return sui.CONFIRM_PROVIDE_MEAL_FOR_CHILD_FUNCTION_V2
}

// GetFunctionCreateChildBooksNeedWithdrawProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildBooksNeedWithdrawProposalV2() string {
	return sui.CREATE_CHILD_BOOKS_NEED_WITHDRAW_PROPOSAL_FUNCTION_V2
}

// GetFunctionCreateChildHealthInsuranceNeedWithdrawProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildHealthInsuranceNeedWithdrawProposalV2() string {
	return sui.CREATE_CHILD_HEALTH_INSURANCE_NEED_WITHDRAW_PROPOSAL_FUNCTION_V2
}

// GetFunctionCreateChildMealNeedWithdrawProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildMealNeedWithdrawProposalV2() string {
	return sui.CREATE_CHILD_MEAL_NEED_WITHDRAW_PROPOSAL_FUNCTION_V2
}

// GetFunctionCreateChildSpecialNeedProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildSpecialNeedProposalV2() string {
	return sui.CREATE_CHILD_SPECIAL_NEED_PROPOSAL_FUNCTION_V2
}

// GetFunctionCreateChildSpecialNeedWithdrawProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildSpecialNeedWithdrawProposalV2() string {
	return sui.CREATE_CHILD_SPECIAL_NEED_PROPOSAL_FUNCTION_V2
}

// GetFunctionWithdrawFromHealthInsuranceNeedProposal implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromHealthInsuranceNeedProposal() string {
	return sui.WITHDRAW_FROM_HEALTH_INSURANCE_NEED_PROPOSAL_FUNCTION
}

// GetFunctionSubmitTask implements IModuleChild.
func (m *moduleChild) GetFunctionSubmitTask() string {
	return sui.SUBMIT_TASK_FUNCTION
}

// ToConfirmProvideMealForChildArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToConfirmProvideMealForChildArgumentsV2(args ConfirmProvideMealForChildArgumentsV2) []interface{} {
	return []interface{}{
		args.ChildID,
		args.NeedID,
		args.StaffNft,
		args.ImageBlobID,
		args.ProvideDate,
		args.Actor,
		sui.CLOCK_OBJECT_ID,
	}
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
	return sui.SUPPORT_CHILD_SPECIAL_NEED_CAMPAIGN_FUNCTION
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

// GetFunctionConfirmProvideMealForChild implements IModuleChild.
func (m *moduleChild) GetFunctionConfirmProvideMealForChild() string {
	return sui.CONFIRM_PROVIDE_MEAL_FOR_CHILD_FUNCTION
}

// ToConfirmProvideMealForChildArguments implements IModuleChild.
func (m *moduleChild) ToConfirmProvideMealForChildArguments(args ConfirmProvideMealForChildArguments) []interface{} {
	return []interface{}{
		args.ChildID,
		args.NeedID,
		args.StaffNft,
		args.ImageBlobID,
		args.ProvideDate,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildSpecialNeedWithdrawProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedWithdrawProposalArguments(args CreateChildSpecialNeedWithdrawProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.CampaignID,
		args.ChildID,
		uint64(args.WithdrawAmount),
		args.Description,
		uint64(args.ClosedAt),
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
		uint64(args.Amount),
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
		os.Getenv(env.SPECIAL_NEED_DAO_ID), // SpecialNeedDao
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
		uint64(args.ClosedAt),
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildSpecialNeedProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedProposalArguments(args CreateChildSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.ChildID,
		args.LocalPool,
		uint64(args.Target),
		args.Description,
		uint64(args.ClosedAt),
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
		uint64(args.Amount),
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
		uint64(args.Amount),
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

// “ToWithdrawFromNeedArguments“ implements IModuleChild.
func (m *moduleChild) ToWithdrawFromNeedArguments(args WithdrawFromNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.TargetID,
		args.ProposalID,
		os.Getenv(env.POOL_WITHDRAW_DAO_OBJECT_ID),
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
		args.HomeAddress,
		args.Region,
		args.AvatarBlobId,
		args.HomeBlobID,
		uint64(time.Now().Year()),
		args.FirstGuardianFullName,
		args.FirstGuardianPhone,
		args.FirstGuardianRelation,
		args.FirstGuardianIdentityCardBlobID,
		args.SecondGuardianFullName,
		args.SecondGuardianPhone,
		args.SecondGuardianRelation,
		args.SecondGuardianIdentityCardBlobID,
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
		args.Leaders,
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
