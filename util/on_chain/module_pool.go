package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type DonateToPoolArguements struct {
	Amount      int64
	FirstName   string
	LastName    string
	Gender      string
	PhoneNumber string
	Email       string
	Message     string
}

type DonateToLocalPoolArguements struct {
	LocalPoolId string
	DonateToPoolArguements
}

type CreateWithdrawProposalArguements struct {
	LocalPoolId     string
	WithdrawAmount  int64
	Description     string
	IsFromLocalPool bool
	ClosedAt        int64
}

type VoteWithdrawProposalArguements struct {
	ProposalId   string
	SponsorId    string
	IsApprove    bool
	RefuseReason string
}

type WithdrawFromPoolArguements struct {
	LocalPoolId        string
	WithdrawProposalId string
}

type IModulePool interface {
	GetModule() string
	ToDonateToPoolArguements(args DonateToPoolArguements) []interface{}
	ToDonateToLocalPoolArguements(args DonateToLocalPoolArguements) []interface{}
	ToCreateWithdrawProposalArguements(args CreateWithdrawProposalArguements) []interface{}
	ToVoteWithdrawProposalArguments(args VoteWithdrawProposalArguements) []interface{}
	ToWithdrawFromPoolArguements(args WithdrawFromPoolArguements) []interface{}
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

// ToWithdrawFromPoolArguements implements IModulePool.
func (m *modulePool) ToWithdrawFromPoolArguements(args WithdrawFromPoolArguements) []interface{} {
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

// ToCreateWithdrawProposalArguements implements IModulePool.
func (m *modulePool) ToCreateWithdrawProposalArguements(args CreateWithdrawProposalArguements) []interface{} {
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
func (m *modulePool) ToVoteWithdrawProposalArguments(args VoteWithdrawProposalArguements) []interface{} {
	return []interface{}{
		args.ProposalId,
		args.SponsorId,
		os.Getenv(env.POOL_WITHDRAW_DAO_OBJECT_ID),
		args.IsApprove,
		args.RefuseReason,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToDonateToLocalPoolArguements implements IModulePool.
func (m *modulePool) ToDonateToLocalPoolArguements(args DonateToLocalPoolArguements) []interface{} {
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

// ToDonateToPoolArguements implements IModulePool.
func (m *modulePool) ToDonateToPoolArguements(args DonateToPoolArguements) []interface{} {
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
