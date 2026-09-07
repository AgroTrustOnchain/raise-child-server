package onchain

import (
	"os"
	"raise-child/constants/env"
	"raise-child/constants/on-chain/sui"
)

type RefundVotePowerArguments struct {
	TargetID   string
	ProposalID string
}

type AdminRefundVotePowerArguments struct {
	RefundVotePowerArguments
	Sender string
}

type IModuleRefund interface {
	GetModule() string
	ToRefundVotePowerArguments(args RefundVotePowerArguments) []interface{}
	ToAdminRefundVotePowerArguments(args AdminRefundVotePowerArguments) []interface{}
	GetFunctionRefundPoolVotePower() string
	GetFunctionAdminRefundPoolVotePower() string
	GetFunctionRefundPoolCampaignVotePower() string
	GetFunctionAdminRefundPoolCampaignVotePower() string
	GetFunctionRefundChildSpecialNeedCampaignVotePower() string
	GetFunctionAdminRefundChildSpecialNeedCampaignVotePower() string
	GetFunctionRefundChildMealNeedVotePower() string
	GetFunctionAdminRefundChildMealNeedVotePower() string
	GetFunctionRefundChildBooksNeedVotePower() string
	GetFunctionAdminRefundChildBooksNeedVotePower() string
	GetFunctionRefundChildHealthInsuranceNeedVotePower() string
	GetFunctionAdminRefundChildHealthInsuranceNeedVotePower() string
}

type moduleRefund struct{}

func InitializeModuleRefund() IModuleRefund {
	return &moduleRefund{}
}

// GetModule implements IModuleRefund.
func (m *moduleRefund) GetModule() string {
	return sui.MODULE_BACKGROUND_REFUND
}

// GetFunctionAdminRefundChildBooksNeedVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionAdminRefundChildBooksNeedVotePower() string {
	return sui.ADMIN_REFUND_CHILD_BOOKS_NEED_VOTE_POWER_FUNCTION
}

// GetFunctionAdminRefundChildHealthInsuranceNeedVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionAdminRefundChildHealthInsuranceNeedVotePower() string {
	return sui.ADMIN_REFUND_CHILD_HEALTH_INSURANCE_NEED_VOTE_POWER_FUNCTION
}

// GetFunctionAdminRefundChildMealNeedVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionAdminRefundChildMealNeedVotePower() string {
	return sui.ADMIN_REFUND_CHILD_MEAL_NEED_VOTE_POWER_FUNCTION
}

// GetFunctionAdminRefundChildSpecialNeedCampaignVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionAdminRefundChildSpecialNeedCampaignVotePower() string {
	return sui.ADMIN_REFUND_CHILD_SPECIAL_NEED_CAMPAIGN_VOTE_POWER_FUNCTION
}

// GetFunctionRefundPoolVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionRefundPoolVotePower() string {
	return sui.REFUND_POOL_VOTE_POWER_FUNCTION
}

// GetFunctionAdminRefundPoolCampaignVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionAdminRefundPoolCampaignVotePower() string {
	return sui.ADMIN_REFUND_POOL_CAMPAIGN_VOTE_POWER_FUNCTION
}

// GetFunctionAdminRefundPoolVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionAdminRefundPoolVotePower() string {
	return sui.ADMIN_REFUND_POOL_VOTE_POWER_FUNCTION
}

// GetFunctionRefundChildBooksNeedVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionRefundChildBooksNeedVotePower() string {
	return sui.REFUND_CHILD_BOOKS_NEED_VOTE_POWER_FUNCTION
}

// GetFunctionRefundChildHealthInsuranceNeedVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionRefundChildHealthInsuranceNeedVotePower() string {
	return sui.REFUND_CHILD_HEALTH_INSURANCE_NEED_VOTE_POWER_FUNCTION
}

// GetFunctionRefundChildMealNeedVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionRefundChildMealNeedVotePower() string {
	return sui.REFUND_CHILD_MEAL_NEED_VOTE_POWER_FUNCTION
}

// GetFunctionRefundChildSpecialNeedCampaignVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionRefundChildSpecialNeedCampaignVotePower() string {
	return sui.REFUND_CHILD_SPECIAL_NEED_CAMPAIGN_VOTE_POWER_FUNCTION
}

// GetFunctionRefundPoolCampaignVotePower implements IModuleRefund.
func (m *moduleRefund) GetFunctionRefundPoolCampaignVotePower() string {
	return sui.REFUND_POOL_CAMPAIGN_VOTE_POWER_FUNCTION
}

// ToAdminRefundVotePowerArguments implements IModuleRefund.
func (m *moduleRefund) ToAdminRefundVotePowerArguments(args AdminRefundVotePowerArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.MANAGE_OBJECT_ID),
		os.Getenv(env.POOL_ID),
		args.TargetID,
		args.ProposalID,
		args.Sender,
		sui.CLOCK_OBJECT_ID,
	}
}

// ToRefundVotePowerArguments implements IModuleRefund.
func (m *moduleRefund) ToRefundVotePowerArguments(args RefundVotePowerArguments) []interface{} {
	return []interface{}{
		os.Getenv(env.ADMIN_CAP_ID_1),
		os.Getenv(env.POOL_ID),
		args.TargetID,
		args.ProposalID,
		sui.CLOCK_OBJECT_ID,
	}
}
