package onchain

import (
	"fmt"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type AddChildArguments struct {
	Center                           string
	IdentityCode                     string
	BirthCertificateBlobID           string
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
	Sender                           string
}

type CreateCenterArguments struct {
	Region      string
	Address     string
	PhoneNumber string
	ImageBlobID string
	Leaders     []string
	Sender      string
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
	Sender      string
}

type SupportChildBooksNeedArgumentsV2 struct {
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
	Creator     string
	CreatedAt   int64
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
	Sender      string
}

type SupportChildHealthInsuranceNeedArgumentsV2 struct {
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
	Creator     string
	CreatedAt   int64
}

type SupportChildMealNeedArguments struct {
	SupportChildBooksNeedArguments
	StartPeriod string
	EndPeriod   string
}

type SupportChildMealNeedArgumentsV2 struct {
	SupportChildBooksNeedArgumentsV2
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
	Sender      string
}

type SupportChildSpeicalNeedArgumentsV2 struct {
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
	Creator     string
	CreatedAt   int64
}

type ConfirmProvideMealForChildArguments struct {
	ChildID     string
	NeedID      string
	StaffNft    string
	ImageBlobID string
	ProvideDate string
}

type ConfirmProvideNeedForChildArgumentsV2 struct {
	ChildID     string
	NeedID      string
	StaffNft    string
	ImageBlobID string
	ProvideDate string
	Actor       string
	Sender      string
}

type CreateChildNormalNeedWithdrawProposalArguments struct {
	NeedID      string
	ChildID     string
	LocalPool   string
	Description string
	ProofBlobID *string
	ClosedAt    int64
	Sender      string
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
	Sender      string
}

type CreateChildSpecialNeedProposalArgumentsV2 struct {
	ChildID     string
	LocalPool   string
	Target      int64
	Description string
	ProofBlobID *string
	ClosedAt    int64
	Creator     string
	Sender      string
}

type CreateChildSpecialNeedWithdrawProposalArguments struct {
	CampaignID     string
	LocalPool      string
	ChildID        string
	WithdrawAmount int64
	Description    string
	ProofBlobID    *string
	ClosedAt       int64
	Sender         string
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
	Sender     string
}

type WithdrawFromNeedArguments struct {
	LocalPool  string
	TargetID   string
	ProposalID string
	Sender     string
}

type WithdrawFromNeedArgumentsV2 struct {
	LocalPool     string
	TargetID      string
	ProposalID    string
	TransferredAt int64
	Creator       string
	Sender        string
}

type SubmitTaskArguments struct {
	Center      string
	StaffNft    string
	Description string
	ImageBlobID string
	Actor       string
}

type UpdateChildNeedArguments struct {
	StaffNft string
	ChildID  string
	NeedID   string
	Year     int
	Value    int64
	Sender   string
}

type IModuleChild interface {
	GetModule() string
	GetChildObjectStruct() string
	ToAddChildArguments(args AddChildArguments) []interface{}
	ToCreateCenterArguments(args CreateCenterArguments) []interface{}
	ToSupportChildBooksNeedArguments(args SupportChildBooksNeedArguments) []interface{}
	ToSupportChildBooksNeedArgumentsV2(args SupportChildBooksNeedArgumentsV2) []interface{}
	ToSupportChildHealthInsuranceNeedArguments(args SupportChildHealthInsuranceNeedArguments) []interface{}
	ToSupportChildHealthInsuranceNeedArgumentsV2(args SupportChildHealthInsuranceNeedArgumentsV2) []interface{}
	ToSupportChildMealNeedArguments(args SupportChildMealNeedArguments) []interface{}
	ToSupportChildMealNeedArgumentsV2(args SupportChildMealNeedArgumentsV2) []interface{}
	ToSupportChildSpeicalNeedArguments(args SupportChildSpeicalNeedArguments) []interface{}
	ToSupportChildSpeicalNeedArgumentsV2(args SupportChildSpeicalNeedArgumentsV2) []interface{}
	ToConfirmProvideMealForChildArguments(args ConfirmProvideMealForChildArguments) []interface{}
	ToConfirmProvideNeedForChildArgumentsV2(args ConfirmProvideNeedForChildArgumentsV2) []interface{}
	ToCreateChildNormalNeedWithdrawProposalArguments(args CreateChildNormalNeedWithdrawProposalArguments) []interface{}
	ToCreateChildSpecialNeedWithdrawProposalArguments(args CreateChildSpecialNeedWithdrawProposalArguments) []interface{}
	ToCreateChildSpecialNeedProposalArguments(args CreateChildSpecialNeedProposalArguments) []interface{}
	ToCreateChildNormalNeedWithdrawProposalArgumentsV2(args CreateChildNormalNeedWithdrawProposalArgumentsV2) []interface{}
	ToCreateChildSpecialNeedWithdrawProposalArgumentsV2(args CreateChildSpecialNeedWithdrawProposalArgumentsV2) []interface{}
	ToCreateChildSpecialNeedProposalArgumentsV2(args CreateChildSpecialNeedProposalArgumentsV2) []interface{}
	ToConfirmChildSpecialNeedProposalArguments(args ConfirmChildSpecialNeedProposalArguments) []interface{}
	ToWithdrawFromNeedArguments(args WithdrawFromNeedArguments) []interface{}
	ToWithdrawFromNeedArgumentsV2(args WithdrawFromNeedArgumentsV2) []interface{}
	ToUpdateChildNeedArguments(args UpdateChildNeedArguments) []interface{}
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
	GetFunctionCreateChildHealthInsuranceNeedWithdrawProposal() string
	GetFunctionCreateChildHealthInsuranceNeedWithdrawProposalV2() string
	GetFunctionCreateChildSpecialNeedProposal() string
	GetFunctionCreateChildSpecialNeedProposalV2() string
	GetFunctionConfirmChildSpecialNeedProposal() string
	GetFunctionWithdrawFromBooksNeedProposal() string
	GetFunctionWithdrawFromBooksNeedProposalV2() string
	GetFunctionWithdrawFromHealthInsuranceNeedProposal() string
	GetFunctionWithdrawFromHealthInsuranceNeedProposalV2() string
	GetFunctionWithdrawFromMealNeedProposal() string
	GetFunctionWithdrawFromMealNeedProposalV2() string
	GetFunctionWithdrawFromSpecialNeedCampaign() string
	GetFunctionWithdrawFromSpecialNeedCampaignV2() string
	GetFunctionSupportChildBooksNeed() string
	GetFunctionSupportChildBooksNeedV2() string
	GetFunctionSupportChildHealthInsuranceNeed() string
	GetFunctionSupportChildHealthInsuranceNeedV2() string
	GetFunctionSupportChildMealNeed() string
	GetFunctionSupportChildMealNeedV2() string
	GetFunctionSupportChildSpecialNeedCampaign() string
	GetFunctionSupportChildSpecialNeedCampaignV2() string
	GetFunctionConfirmProvideMealForChild() string
	GetFunctionConfirmProvideMealForChildV2() string
	GetFunctionConfirmProvideBooksForChildV2() string
	GetFunctionConfirmProvideHealthInsuranceForChildV2() string
	GetFunctionSubmitTask() string
	GetFunctionUpdateChildMealNeed() string
	GetFunctionUpdateChildBooksNeed() string
	GetFunctionUpdateChildHealthInsuranceNeed() string
	GetChildEventEmittedStruct() string
}

type moduleChild struct{}

func InitializeModuleChild() IModuleChild {
	return &moduleChild{}
}

// GetFunctionConfirmProvideBooksForChildV2 implements IModuleChild.
func (m *moduleChild) GetFunctionConfirmProvideBooksForChildV2() string {
	return sui.CONFIRM_PROVIDE_BOOKS_FOR_CHILD_FUNCTION_V2
}

// GetChildEventEmittedStruct implements IModuleChild.
func (m *moduleChild) GetChildEventEmittedStruct() string {
	return sui.CHILD_CREATED_EVENT
}

// GetFunctionConfirmProvideHealthInsuranceForChildV2 implements IModuleChild.
func (m *moduleChild) GetFunctionConfirmProvideHealthInsuranceForChildV2() string {
	return sui.CONFIRM_PROVIDE_HEALTH_INSURANCE_FOR_CHILD_FUNCTION_V2
}

// GetFunctionWithdrawFromBooksNeedProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromBooksNeedProposalV2() string {
	return sui.WITHDRAW_FROM_BOOKS_NEED_PROPOSAL_FUNCTION_V2
}

// GetFunctionWithdrawFromHealthInsuranceNeedProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromHealthInsuranceNeedProposalV2() string {
	return sui.WITHDRAW_FROM_HEALTH_INSURANCE_NEED_PROPOSAL_FUNCTION_V2
}

// GetFunctionWithdrawFromMealNeedProposalV2 implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromMealNeedProposalV2() string {
	return sui.WITHDRAW_FROM_MEAL_NEED_PROPOSAL_FUNCTION_V2
}

// GetFunctionWithdrawFromSpecialNeedCampaignV2 implements IModuleChild.
func (m *moduleChild) GetFunctionWithdrawFromSpecialNeedCampaignV2() string {
	return sui.WITHDRAW_FROM_SPECIAL_NEED_CAMPAIGN_FUNCTION_v2
}

// GetFunctionSupportChildBooksNeedV2 implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildBooksNeedV2() string {
	return sui.SUPPORT_CHILD_BOOKS_NEED_FUNCTION_V2
}

// GetFunctionSupportChildHealthInsuranceNeedV2 implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildHealthInsuranceNeedV2() string {
	return sui.SUPPORT_CHILD_HEALTH_INSURANCE_NEED_FUNCTION_V2
}

// GetFunctionSupportChildMealNeedV2 implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildMealNeedV2() string {
	return sui.SUPPORT_CHILD_MEAL_NEED_FUNCTION_V2
}

// GetFunctionSupportChildSpecialNeedCampaignV2 implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildSpecialNeedCampaignV2() string {
	return sui.SUPPORT_CHILD_SPECIAL_NEED_CAMPAIGN_FUNCTION_V2
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
		fmt.Sprintf("%d", args.ClosedAt),
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
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.ChildID,
		args.LocalPool,
		fmt.Sprintf("%d", args.Target),
		args.Description,
		proofBlobId,
		fmt.Sprintf("%d", args.ClosedAt),
		args.Creator,
		args.Sender,
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
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.CampaignID,
		args.ChildID,
		fmt.Sprintf("%d", args.WithdrawAmount),
		args.Description,
		proofBlobId,
		fmt.Sprintf("%d", args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionSupportChildHealthInsuranceNeed implements IModuleChild.
func (m *moduleChild) GetFunctionSupportChildHealthInsuranceNeed() string {
	return sui.SUPPORT_CHILD_HEALTH_INSURANCE_NEED_FUNCTION
}

// ToSupportChildHealthInsuranceNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildHealthInsuranceNeedArguments(args SupportChildHealthInsuranceNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.NeedID,
		args.ChildID,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildBooksNeedArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToSupportChildBooksNeedArgumentsV2(args SupportChildBooksNeedArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.NeedID,
		args.LocalPool,
		args.ChildID,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Creator,
		fmt.Sprintf("%d", args.CreatedAt),
	}
}

// ToSupportChildHealthInsuranceNeedArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToSupportChildHealthInsuranceNeedArgumentsV2(args SupportChildHealthInsuranceNeedArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.NeedID,
		args.LocalPool,
		args.ChildID,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Creator,
		fmt.Sprintf("%d", args.CreatedAt),
	}
}

// ToSupportChildMealNeedArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToSupportChildMealNeedArgumentsV2(args SupportChildMealNeedArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.NeedID,
		args.LocalPool,
		args.ChildID,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.StartPeriod,
		args.EndPeriod,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Creator,
		fmt.Sprintf("%d", args.CreatedAt),
	}
}

// ToSupportChildSpeicalNeedArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToSupportChildSpeicalNeedArgumentsV2(args SupportChildSpeicalNeedArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.CampaignID,
		args.ChildID,
		args.LocalPool,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Creator,
		fmt.Sprintf("%d", args.CreatedAt),
	}
}

// ToSubmitTaskArguments implements IModuleChild.
func (m *moduleChild) ToSubmitTaskArguments(args SubmitTaskArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.Center,
		args.StaffNft,
		args.Description,
		args.ImageBlobID,
		args.Actor,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToUpdateChildNeedArguments implements IModuleChild.
func (m *moduleChild) ToUpdateChildNeedArguments(args UpdateChildNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.StaffNft,
		args.ChildID,
		args.NeedID,
		fmt.Sprintf("%d", args.Year),
		fmt.Sprintf("%d", args.Value),
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionUpdateChildBooksNeed implements IModuleChild.
func (m *moduleChild) GetFunctionUpdateChildBooksNeed() string {
	return sui.UPDATE_CHILD_BOOKS_NEED_FUNCTION
}

// GetFunctionUpdateChildHealthInsuranceNeed implements IModuleChild.
func (m *moduleChild) GetFunctionUpdateChildHealthInsuranceNeed() string {
	return sui.UPDATE_CHILD_HEALTH_INSURANCE_NEED_FUNCTION
}

// GetFunctionUpdateChildMealNeed implements IModuleChild.
func (m *moduleChild) GetFunctionUpdateChildMealNeed() string {
	return sui.UPDATE_CHILD_MEAL_NEED_FUNCTION
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

// GetFunctionCreateChildHealthInsuranceNeedWithdrawProposal implements IModuleChild.
func (m *moduleChild) GetFunctionCreateChildHealthInsuranceNeedWithdrawProposal() string {
	return sui.CREATE_CHILD_HEALTH_INSURANCE_NEED_WITHDRAW_PROPOSAL_FUNCTION
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
func (m *moduleChild) ToConfirmProvideNeedForChildArgumentsV2(args ConfirmProvideNeedForChildArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.ChildID,
		args.NeedID,
		args.StaffNft,
		args.ImageBlobID,
		args.ProvideDate,
		args.Actor,
		args.Sender,
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
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.CampaignID,
		args.ChildID,
		fmt.Sprintf("%d", args.WithdrawAmount),
		args.Description,
		proofBlobId,
		fmt.Sprintf("%d", args.ClosedAt),
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildSpeicalNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildSpeicalNeedArguments(args SupportChildSpeicalNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.CampaignID,
		args.ChildID,
		args.LocalPool,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToConfirmChildSpecialNeedProposalArguments implements IModuleChild.
func (m *moduleChild) ToConfirmChildSpecialNeedProposalArguments(args ConfirmChildSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.SPECIAL_NEED_DAO_ID), // SpecialNeedDao
		args.ProposalID,
		args.ChildID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildNormalNeedWithdrawProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildNormalNeedWithdrawProposalArguments(args CreateChildNormalNeedWithdrawProposalArguments) []interface{} {
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.NeedID,
		args.ChildID,
		args.Description,
		proofBlobId,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateChildSpecialNeedProposalArguments implements IModuleChild.
func (m *moduleChild) ToCreateChildSpecialNeedProposalArguments(args CreateChildSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.ChildID,
		args.LocalPool,
		fmt.Sprintf("%d", args.Target),
		args.Description,
		fmt.Sprintf("%d", args.ClosedAt),
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildBooksNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildBooksNeedArguments(args SupportChildBooksNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.NeedID,
		args.ChildID,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportChildMealNeedArguments implements IModuleChild.
func (m *moduleChild) ToSupportChildMealNeedArguments(args SupportChildMealNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.NeedID,
		args.ChildID,
		args.DonorNft,
		fmt.Sprintf("%d", args.Amount),
		args.StartPeriod,
		args.EndPeriod,
		args.FirstName,
		args.LastName,
		args.Gender,
		args.PhoneNumber,
		args.Email,
		args.Message,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// “ToWithdrawFromNeedArguments“ implements IModuleChild.
func (m *moduleChild) ToWithdrawFromNeedArguments(args WithdrawFromNeedArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.TargetID,
		args.ProposalID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToWithdrawFromNeedArgumentsV2 implements IModuleChild.
func (m *moduleChild) ToWithdrawFromNeedArgumentsV2(args WithdrawFromNeedArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPool,
		args.TargetID,
		args.ProposalID,
		fmt.Sprintf("%d", args.TransferredAt),
		args.Creator,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToAddChildArguments implements IModuleChild.
func (m *moduleChild) ToAddChildArguments(args AddChildArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
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
		args.FirstGuardianFullName,
		args.FirstGuardianPhone,
		args.FirstGuardianRelation,
		args.SecondGuardianFullName,
		args.SecondGuardianPhone,
		args.SecondGuardianRelation,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateCenterArguments implements IModuleChild.
func (m *moduleChild) ToCreateCenterArguments(args CreateCenterArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.Region,
		args.Address,
		args.PhoneNumber,
		args.ImageBlobID,
		args.Leaders,
		args.Sender,
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
