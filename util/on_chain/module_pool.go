package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type DonateToPoolArguments struct {
	DonorID     string
	Amount      int64
	FirstName   string
	LastName    string
	Gender      string
	PhoneNumber string
	Email       string
	Message     string
}

type DonateToLocalPoolArguments struct {
	LocalPoolId string
	DonateToPoolArguments
}

type CreateWithdrawProposalArguments struct {
	LocalPoolId     string
	WithdrawAmount  int64
	Description     string
	IsFromLocalPool bool
	ClosedAt        int64
}

type CreateWithdrawProposalV2Arguments struct {
	LocalPoolId     string
	WithdrawAmount  int64
	Description     string
	ProofBlobID     *string
	IsFromLocalPool bool
	ClosedAt        int64
	Creator         string
}

type VoteWithdrawProposalArguments struct {
	ProposalId   string
	DonorId      string
	IsApprove    bool
	RefuseReason string
}

type WithdrawFromPoolArguments struct {
	LocalPoolId        string
	WithdrawProposalId string
}

type EditWithdrawDaoRateArguements struct {
	MinRate   int64
	MinVoters int
}

type IModulePool interface {
	GetModule() string
	ToDonateToPoolArguments(args DonateToPoolArguments) []interface{}
	ToDonateToLocalPoolArguments(args DonateToLocalPoolArguments) []interface{}
	ToCreateWithdrawProposalArguments(args CreateWithdrawProposalArguments) []interface{}
	ToCreateWithdrawProposalV2Arguments(args CreateWithdrawProposalV2Arguments) []interface{}
	ToVoteWithdrawProposalArguments(args VoteWithdrawProposalArguments) []interface{}
	ToWithdrawFromPoolArguments(args WithdrawFromPoolArguments) []interface{}
	ToEditWithdrawDaoRateArguements(args EditWithdrawDaoRateArguements) []interface{}
	GetWithdrawProposalEventEmittedStruct() string
	GetFunctionDonateToPool() string
	GetFunctionDonateToLocalPool() string
	GetFunctionWithdrawFromPool() string
	GetFunctionCreateWithdrawProposal() string
	GetFunctionCreateWithdrawProposalV2() string
	GetFunctionVoteWithdrawProposal() string
	GetFunctionEditWithdrawDaoRate() string
}

type modulePool struct{}

func InitializeModulePool() IModulePool {
	return &modulePool{}
}

// GetFunctionEditWithdrawDaoRate implements IModulePool.
func (m *modulePool) GetFunctionEditWithdrawDaoRate() string {
	return sui.EDIT_WITHDRAW_DAO_RATE_FUNCTION
}

// GetFunctionDonateToLocalPool implements IModulePool.
func (m *modulePool) GetFunctionDonateToLocalPool() string {
	return sui.DONATE_TO_LOCAL_POOL_FUNCTION
}

// GetFunctionCreateWithdrawProposal implements IModulePool.
func (m *modulePool) GetFunctionCreateWithdrawProposal() string {
	return sui.CREATE_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionCreateWithdrawProposalV2 implements IModulePool.
func (m *modulePool) GetFunctionCreateWithdrawProposalV2() string {
	return sui.CREATE_WITHDRAW_PROPOSAL_V2_FUNCTION
}

// GetFunctionVoteWithdrawProposal implements IModulePool.
func (m *modulePool) GetFunctionVoteWithdrawProposal() string {
	return sui.VOTE_WITHDRAW_PROPOSAL_FUNCTION
}

// GetWithdrawProposalEventEmittedStruct implements IModulePool.
func (m *modulePool) GetWithdrawProposalEventEmittedStruct() string {
	return sui.WITHDRAW_PROPOSAL_EVENT
}

// ToEditWithdrawDaoRateArguements implements IModulePool.
func (m *modulePool) ToEditWithdrawDaoRateArguements(args EditWithdrawDaoRateArguements) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.POOL_WITHDRAW_DAO_OBJECT_ID),
		uint64(args.MinRate),
		args.MinVoters,
	}
}

// ToWithdrawFromPoolArguments implements IModulePool.
func (m *modulePool) ToWithdrawFromPoolArguments(args WithdrawFromPoolArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolId,
		args.WithdrawProposalId,
		os.Getenv(env.POOL_WITHDRAW_DAO_OBJECT_ID),
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateWithdrawProposalArguments implements IModulePool.
func (m *modulePool) ToCreateWithdrawProposalArguments(args CreateWithdrawProposalArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolId,
		uint64(args.WithdrawAmount),
		args.Description,
		args.IsFromLocalPool,
		uint64(args.ClosedAt),
		sui.CLOCK_OBJECT_ID,
	}
}

// ToCreateWithdrawProposalV2Arguments implements IModulePool.
func (m *modulePool) ToCreateWithdrawProposalV2Arguments(args CreateWithdrawProposalV2Arguments) []interface{} {
	var proofBlobId string = ""
	if args.ProofBlobID != nil {
		proofBlobId = *args.ProofBlobID
	}

	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolId,
		uint64(args.WithdrawAmount),
		args.Description,
		proofBlobId,
		args.IsFromLocalPool,
		uint64(args.ClosedAt),
		args.Creator,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToVoteWithdrawProposal implements IModulePool.
func (m *modulePool) ToVoteWithdrawProposalArguments(args VoteWithdrawProposalArguments) []interface{} {
	return []interface{}{
		args.ProposalId,
		args.DonorId,
		os.Getenv(env.POOL_WITHDRAW_DAO_OBJECT_ID),
		args.IsApprove,
		args.RefuseReason,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToDonateToLocalPoolArguments implements IModulePool.
func (m *modulePool) ToDonateToLocalPoolArguments(args DonateToLocalPoolArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolId,
		args.DonorID,
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

// ToDonateToPoolArguments implements IModulePool.
func (m *modulePool) ToDonateToPoolArguments(args DonateToPoolArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.DonorID,
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

// GetFunctionDonateToPool implements IModulePool.
func (m *modulePool) GetFunctionDonateToPool() string {
	return sui.DONATE_TO_POOL_FUNCTION
}

// GetFunctionWithdrawFromPool implements IModulePool.
func (m *modulePool) GetFunctionWithdrawFromPool() string {
	return sui.WITHDRAW_FROM_POOL_FUNCTION
}

// GetModule implements IModulePool.
func (m *modulePool) GetModule() string {
	return sui.MODULE_POOL
}
