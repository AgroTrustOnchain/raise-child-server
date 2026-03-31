package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type EditSpecialNeedProposalDaoArguments struct {
	MinVoters int
	MinRate   int64
}

type VoteSpecialNeedProposalArguments struct {
	ProposalID   string
	DonorNft     string
	IsApprove    bool
	RefuseReason string
}

type EditUpdateNeedDatesArguments struct {
	EditDatesID string
	StartDate   string
	EndDate     string
}

type IModuleNeed interface {
	GetModule() string
	GetSpecialNeedDaoObjectStruct() string
	ToEditSpecialNeedProposalDaoArguments(args EditSpecialNeedProposalDaoArguments) []interface{}
	ToVoteSpecialNeedProposalArguments(args VoteSpecialNeedProposalArguments) []interface{}
	ToEditUpdateNeedDatesArguments(args EditUpdateNeedDatesArguments) []interface{}
	GetFunctionEditSpecialNeedProposalDao() string
	GetFunctionVoteSpecialNeedProposal() string
	GetFunctionEditUpdateBooksNeedDates() string
	GetFunctionEditUpdateMealNeedDates() string
	GetFunctionEditUpdateHealthInsuranceNeedDates() string
}

type moduleNeed struct{}

func InitializeModuleNeed() IModuleNeed {
	return &moduleNeed{}
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
		os.Getenv(env.MANAGE_OBJECT_ID),
		args.EditDatesID,
		args.StartDate,
		args.EndDate,
	}
}

// ToEditSpecialNeedProposalDaoArguments implements IModuleNeed.
func (m *moduleNeed) ToEditSpecialNeedProposalDaoArguments(args EditSpecialNeedProposalDaoArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.SPECIAL_NEED_DAO_ID),
		uint64(args.MinRate),
		args.MinVoters,
	}
}

// ToVoteSpecialNeedProposalArguments implements IModuleNeed.
func (m *moduleNeed) ToVoteSpecialNeedProposalArguments(args VoteSpecialNeedProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.SPECIAL_NEED_DAO_ID),
		args.ProposalID,
		args.DonorNft,
		args.IsApprove,
		args.RefuseReason,
		sui.CLOCK_OBJECT_ID,
	}
}
