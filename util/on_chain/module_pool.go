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

type IModulePool interface {
	GetModule() string
	ToDonateToPoolArguements(args DonateToPoolArguements) []interface{}
	GetFunctionDonateToPool() string
	GetFunctionWithdrawFromPool() string
}

type modulePool struct{}

func InitializeModulePool() IModulePool {
	return &modulePool{}
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
