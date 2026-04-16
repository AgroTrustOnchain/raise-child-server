package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type CreateCampaignForMainPoolArguments struct {
	Creator     string
	Target      int64
	Description string
	ProofBlobID *string
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
}

type IModuleCampaign interface {
	GetModule() string
	GetCampaignObjectStruct() string
	ToCreateCampaignForMainPoolArguments(args CreateCampaignForMainPoolArguments) []interface{}
	ToCreateCampaignForRegionPoolArguments(args CreateCampaignForRegionPoolArguments) []interface{}
	ToSupportCampaignArguments(args SupportCampaignArguments) []interface{}
	ToCreateCampaignWithdrawProposalArguments(args CreateCampaignWithdrawProposalArguments) []interface{}
	ToWithdrawFromCampaignArguments(args WithdrawFromCampaignArguments) []interface{}
	GetFunctionCreateCampaignForMainPool() string
	GetFunctionCreateCampaignForRegionPool() string
	GetFunctionSupportCampaign() string
	GetFunctionCreateCampaignWithdrawProposal() string
	GetFunctionWithdrawFromCampaign() string
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

// GetFunctionWithdrawFromCampaign implements IModuleCampaign.
func (m *moduleCampaign) GetFunctionWithdrawFromCampaign() string {
	return sui.WITHDRAW_FROM_CAMPAIGN_FUNCTION
}

// GetModule implements IModuleCampaign.
func (m *moduleCampaign) GetModule() string {
	return sui.MODULE_CAMPAIGN
}

// ToCreateCampaignForMainPoolArguments implements IModuleCampaign.
func (m *moduleCampaign) ToCreateCampaignForMainPoolArguments(args CreateCampaignForMainPoolArguments) []interface{} {
	var proofBlobId string
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.Creator,
		uint64(args.Target),
		args.Description,
		proofBlobId,
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
		uint64(args.Target),
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
		uint64(args.WithdrawAmount),
		args.Description,
		proofBlobId,
		uint64(args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToSupportCampaignArguments implements IModuleCampaign.
func (m *moduleCampaign) ToSupportCampaignArguments(args SupportCampaignArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolID,
		args.CampaignID,
		args.DonorNFT,
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

// ToWithdrawFromCampaignArguments implements IModuleCampaign.
func (m *moduleCampaign) ToWithdrawFromCampaignArguments(args WithdrawFromCampaignArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.CampaignID,
		args.ProposalID,
		os.Getenv(env.POOL_WITHDRAW_DAO_OBJECT_ID),
		sui.CLOCK_OBJECT_ID,
	}
}
