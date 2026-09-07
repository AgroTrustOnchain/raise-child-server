package onchain

import (
	"fmt"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type CreateCampaignForMainPoolArguments struct {
	Creator     string
	Target      int64
	Description string
	ProofBlobID *string
	Sender      string
}

type CreateCampaignForRegionPoolArguments struct {
	LocalPoolID string
	CreateCampaignForMainPoolArguments
}

type SupportCampaignArguments struct {
	LocalPoolID string
	CampaignID  string
	DonorNFT    string
	Amount      int64
	FirstName   string
	LastName    string
	Gender      string
	PhoneNumber string
	Email       string
	Message     string
	Sender      string
}

type SupportCampaignArgumentsV2 struct {
	LocalPoolID string
	CampaignID  string
	DonorNFT    string
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

type CreateCampaignWithdrawProposalArguments struct {
	LocalPoolID    string
	CampaignID     string
	WithdrawAmount int64
	Description    string
	ProofBlobID    *string
	ClosedAt       int64
	Creator        string
}

type WithdrawFromCampaignArguments struct {
	CampaignID string
	ProposalID string
	Sender     string
}

type WithdrawFromCampaignArgumentsV2 struct {
	CampaignID    string
	ProposalID    string
	TransferredAt int64
	Creator       string
	Sender        string
}

type VotePoolCampaignWithdrawProposalArguments struct {
	CampaignID string
	ProposalID string
	Sender     string
}

type IModuleCampaign interface {
	GetModule() string
	GetCampaignObjectStruct() string
	ToCreateCampaignForMainPoolArguments(args CreateCampaignForMainPoolArguments) []interface{}
	ToCreateCampaignForRegionPoolArguments(args CreateCampaignForRegionPoolArguments) []interface{}
	ToSupportCampaignArguments(args SupportCampaignArguments) []interface{}
	ToSupportCampaignArgumentsV2(args SupportCampaignArgumentsV2) []interface{}
	ToCreateCampaignWithdrawProposalArguments(args CreateCampaignWithdrawProposalArguments) []interface{}
	ToWithdrawFromCampaignArguments(args WithdrawFromCampaignArguments) []interface{}
	ToWithdrawFromCampaignArgumentsV2(args WithdrawFromCampaignArgumentsV2) []interface{}
	ToVotePoolCampaignWithdrawProposalArguments(args VotePoolCampaignWithdrawProposalArguments) []interface{}
	GetFunctionCreateCampaignForMainPool() string
	GetFunctionCreateCampaignForRegionPool() string
	GetFunctionSupportCampaign() string
	GetFunctionSupportCampaignV2() string
	GetFunctionCreateCampaignWithdrawProposal() string
	GetFunctionWithdrawFromCampaign() string
	GetFunctionWithdrawFromCampaignV2() string
	GetFunctionVotePoolCampaignWithdrawProposal() string
}

type moduleCampaign struct{}

func InitializeModuleCampaign() IModuleCampaign {
	return &moduleCampaign{}
}

// GetCampaignObjectStruct implements IModuleCampaign.
func (m *moduleCampaign) GetCampaignObjectStruct() string {
	return sui.CAMPAIGN_STRUCT
}

// GetFunctionCreateCampaignForMainPool implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionCreateCampaignForMainPool() string {
	return sui.CREATE_CAMPAIGN_FOR_MAIN_POOL_FUNCTION
}

// GetFunctionCreateCampaignForRegionPool implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionCreateCampaignForRegionPool() string {
	return sui.CREATE_CAMPAIGN_FOR_REGION_POOL_FUNCTION
}

// GetFunctionCreateCampaignWithdrawProposal implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionCreateCampaignWithdrawProposal() string {
	return sui.CREATE_CAMPAIGN_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionSupportCampaign implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionSupportCampaign() string {
	return sui.SUPPORT_CAMPAIGN_FUNCTION
}

// GetFunctionSupportCampaignV2 implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionSupportCampaignV2() string {
	return sui.SUPPORT_CAMPAIGN_FUNCTION_v2
}

// GetFunctionWithdrawFromCampaign implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionWithdrawFromCampaign() string {
	return sui.WITHDRAW_FROM_CAMPAIGN_FUNCTION
}

// GetFunctionVotePoolCampaignWithdrawProposal implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionVotePoolCampaignWithdrawProposal() string {
	return sui.VOTE_POOL_CAMPAIGN_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionWithdrawFromCampaignV2 implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionWithdrawFromCampaignV2() string {
	return sui.WITHDRAW_FROM_CAMPAIGN_FUNCTION_V2
}

// GetModule implements IModuleCampaign.
func (m *moduleCampaign) GetModule() string {
	return sui.MODULE_CAMPAIGN
}

// ToVotePoolCampaignWithdrawProposalArguments implements IModuleCampaign.
func (m *moduleCampaign) ToVotePoolCampaignWithdrawProposalArguments(args VotePoolCampaignWithdrawProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.CampaignID,
		args.ProposalID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToWithdrawFromCampaignArgumentsV2 implements IModuleCampaign.
func (m *moduleCampaign) ToWithdrawFromCampaignArgumentsV2(args WithdrawFromCampaignArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.POOL_ID),
		args.CampaignID,
		args.ProposalID,
		fmt.Sprintf("%d", args.TransferredAt),
		args.Creator,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateCampaignForMainPoolArguments implements IModuleCampaign.
func (m *moduleCampaign) ToCreateCampaignForMainPoolArguments(args CreateCampaignForMainPoolArguments) []interface{} {
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.Creator,
		fmt.Sprintf("%d", args.Target),
		args.Description,
		proofBlobId,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateCampaignForRegionPoolArguments implements IModuleCampaign.
func (m *moduleCampaign) ToCreateCampaignForRegionPoolArguments(args CreateCampaignForRegionPoolArguments) []interface{} {
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolID,
		args.Creator,
		fmt.Sprintf("%d", args.Target),
		args.Description,
		proofBlobId,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateCampaignWithdrawProposalArguments implements IModuleCampaign.
func (m *moduleCampaign) ToCreateCampaignWithdrawProposalArguments(args CreateCampaignWithdrawProposalArguments) []interface{} {
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolID,
		args.CampaignID,
		fmt.Sprintf("%d", args.WithdrawAmount),
		args.Description,
		proofBlobId,
		fmt.Sprintf("%d", args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportCampaignArguments implements IModuleCampaign.
func (m *moduleCampaign) ToSupportCampaignArguments(args SupportCampaignArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolID,
		args.CampaignID,
		args.DonorNFT,
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

// ToSupportCampaignArgumentsV2 implements IModuleCampaign.
func (m *moduleCampaign) ToSupportCampaignArgumentsV2(args SupportCampaignArgumentsV2) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolID,
		args.CampaignID,
		args.DonorNFT,
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

// ToWithdrawFromCampaignArguments implements IModuleCampaign.
func (m *moduleCampaign) ToWithdrawFromCampaignArguments(args WithdrawFromCampaignArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.CampaignID,
		args.ProposalID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}
