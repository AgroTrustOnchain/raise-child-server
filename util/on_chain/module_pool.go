package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type DonateToPoolArguments struct {
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

type VoteWithdrawProposalArguments struct {
	ProposalId   string
	SponsorId    string
	IsApprove    bool
	RefuseReason string
}

type WithdrawFromPoolArguments struct {
	LocalPoolId        string
	WithdrawProposalId string
}

type IModulePool interface {
	GetModule() string
	ToDonateToPoolArguments(args DonateToPoolArguments) []interface{}
	ToDonateToLocalPoolArguments(args DonateToLocalPoolArguments) []interface{}
	ToCreateWithdrawProposalArguments(args CreateWithdrawProposalArguments) []interface{}
	ToVoteWithdrawProposalArguments(args VoteWithdrawProposalArguments) []interface{}
	ToWithdrawFromPoolArguments(args WithdrawFromPoolArguments) []interface{}
	GetFunctionDonateToPool() string
	GetFunctionDonateToLocalPool() string
	GetFunctionWithdrawFromPool() string
	GetFunctionCreateWithdrawProposal() string
	GetFunctionVoteWithdrawProposal() string
}

type modulePool struct{}

func InitializeModulePool() IModulePool {
	return &modulePool{}
}

// GetFunctionDonateToLocalPool implements IModulePool.
func (m *modulePool) GetFunctionDonateToLocalPool() string {
	return sui.DONATE_TO_LOCAL_POOL_FUNCTION
}

// GetFunctionCreateWithdrawProposal implements IModulePool.
func (m *modulePool) GetFunctionCreateWithdrawProposal() string {
	return sui.CREATE_WITHDRAW_PROPOSAL_FUNCTION
}

// GetFunctionVoteWithdrawProposal implements IModulePool.
func (m *modulePool) GetFunctionVoteWithdrawProposal() string {
	return sui.VOTE_WITHDRAW_PROPOSAL_FUNCTION
}

// ToWithdrawFromPoolArguments implements IModulePool.
func (m *modulePool) ToWithdrawFromPoolArguments(args WithdrawFromPoolArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.LocalPoolId,
		args.WithdrawProposalId,
		os.Getenv(env.TREASURY_CAP),
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
		args.WithdrawAmount,
		args.Description,
		args.IsFromLocalPool,
		args.ClosedAt,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToVoteWithdrawProposal implements IModulePool.
func (m *modulePool) ToVoteWithdrawProposalArguments(args VoteWithdrawProposalArguments) []interface{} {
	return []interface{}{
		args.ProposalId,
		args.SponsorId,
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
		os.Getenv(env.TREASURY_CAP),
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
		os.Getenv(env.TREASURY_CAP),
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
