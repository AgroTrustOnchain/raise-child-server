package onchain

import (
	"fmt"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type EditSpecialNeedProposalDaoArguments struct {
	MinVoters int
	MinRate   int64
	Sender    string
}

type VoteSpecialNeedProposalArguments struct {
	ProposalID   string
	DonorNft     string
	IsApprove    bool
	RefuseReason string
	Sender       string
}

type EditUpdateNeedDatesArguments struct {
	EditDatesID string
	StartDate   string
	EndDate     string
	Sender      string
}

type VoteChildNeedWithdrawProposalArguments struct {
	TargetID   string
	ProposalID string
	Sender     string
}

type IModuleNeed interface {
	GetModule() string
	GetSpecialNeedDaoObjectStruct() string
	ToEditSpecialNeedProposalDaoArguments(args EditSpecialNeedProposalDaoArguments) []interface{}
	ToVoteSpecialNeedProposalArguments(args VoteSpecialNeedProposalArguments) []interface{}
	ToEditUpdateNeedDatesArguments(args EditUpdateNeedDatesArguments) []interface{}
	ToVoteChildNeedWithdrawProposalArguments(args VoteChildNeedWithdrawProposalArguments) []interface{}
	GetFunctionEditSpecialNeedProposalDao() string
	GetFunctionVoteSpecialNeedProposal() string
	GetFunctionEditUpdateBooksNeedDates() string
	GetFunctionEditUpdateMealNeedDates() string
	GetFunctionEditUpdateHealthInsuranceNeedDates() string
	GetFunctionVoteBooksNeedWithdrawProposal() string
	GetFunctionVoteMealNeedWithdrawProposal() string
	GetFunctionVoteHealthInsuranceNeedWithdrawProposal() string
	GetFunctionVoteSpecialNeedCampaignWithdrawProposal() string
}

type moduleNeed struct{}

func InitializeModuleNeed() IModuleNeed {
	return &moduleNeed{}
}

// GetFunctionVoteBooksNeedWithdrawProposal implements IModuleNeed.
func (m *moduleNeed) GetFunctionVoteBooksNeedWithdrawProposal() string {
	return sui.VOTE_BOOKS_NEED_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionVoteHealthInsuranceNeedWithdrawProposal implements IModuleNeed.
func (m *moduleNeed) GetFunctionVoteHealthInsuranceNeedWithdrawProposal() string {
	return sui.VOTE_HEALTH_INSURANCE_NEED_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionVoteMealNeedWithdrawProposal implements IModuleNeed.
func (m *moduleNeed) GetFunctionVoteMealNeedWithdrawProposal() string {
	return sui.VOTE_MEAL_NEED_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionVoteSpecialNeedCampaignWithdrawProposal implements IModuleNeed.
func (m *moduleNeed) GetFunctionVoteSpecialNeedCampaignWithdrawProposal() string {
	return sui.VOTE_SPECIAL_NEED_CAMPAIGN_WITHDRAW_PROPOSAL_FUNCTION
}

// ToVoteChildNeedWithdrawProposalArguments implements IModuleNeed.
func (m *moduleNeed) ToVoteChildNeedWithdrawProposalArguments(args VoteChildNeedWithdrawProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		args.TargetID,
		args.ProposalID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// GetFunctionEditSpecialNeedProposalDao implements IModuleNeed.
func (m *moduleNeed) GetFunctionEditSpecialNeedProposalDao() string {
	return sui.EDIT_SPECIAL_NEED_DAO_RATE_FUNCTION
}

// GetSpecialNeedDaoObjectStruct implements IModuleNeed.
func (m *moduleNeed) GetSpecialNeedDaoObjectStruct() string {
	return sui.SPECIAL_NEED_DAO_STRUCT
}

// GetFunctionVoteSpecialNeedProposal implements IModuleNeed.
func (m *moduleNeed) GetFunctionVoteSpecialNeedProposal() string {
	return sui.VOTE_SPECIAL_NEED_PROPOSAL_FUNCTION
}

// GetModule implements IModuleNeed.
func (m *moduleNeed) GetModule() string {
	return sui.MODULE_NEED
}

// GetFunctionEditUpdateBooksNeedDates implements IModuleNeed.
func (m *moduleNeed) GetFunctionEditUpdateBooksNeedDates() string {
	return sui.EDIT_UPDATE_BOOKS_NEED_DATES_FUNCTION
}

// GetFunctionEditUpdateHealthInsuranceNeedDates implements IModuleNeed.
func (m *moduleNeed) GetFunctionEditUpdateHealthInsuranceNeedDates() string {
	return sui.EDIT_UPDATE_HEALTH_INSURANCE_NEED_DATES_FUNCTION
}

// GetFunctionEditUpdateMealNeedDates implements IModuleNeed.
func (m *moduleNeed) GetFunctionEditUpdateMealNeedDates() string {
	return sui.EDIT_UPDATE_MEAL_NEED_DATES_FUNCTION
}

// ToEditUpdateNeedDatesArguments implements IModuleNeed.
func (m *moduleNeed) ToEditUpdateNeedDatesArguments(args EditUpdateNeedDatesArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.EditDatesID,
		args.StartDate,
		args.EndDate,
		args.Sender,
	}
}

// ToEditSpecialNeedProposalDaoArguments implements IModuleNeed.
func (m *moduleNeed) ToEditSpecialNeedProposalDaoArguments(args EditSpecialNeedProposalDaoArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.SPECIAL_NEED_DAO_ID),
		fmt.Sprintf("%d", args.MinRate),
		fmt.Sprintf("%d", args.MinVoters),
		args.Sender,
	}
}

// ToVoteSpecialNeedProposalArguments implements IModuleNeed.
func (m *moduleNeed) ToVoteSpecialNeedProposalArguments(args VoteSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.SPECIAL_NEED_DAO_ID),
		args.ProposalID,
		args.DonorNft,
		args.IsApprove,
		args.RefuseReason,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}
