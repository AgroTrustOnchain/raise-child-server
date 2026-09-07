package business

import "context"

type IBackgroundService interface {
	ProcessBackgroundCenterRequests(ctx context.Context)
	ProcessBackgroundRegistrationRequests(ctx context.Context)
	ProcessRefundVotePower(ctx context.Context)
	ProcessCreateChildrenWithdrawProposals(ctx context.Context)
	//ProcessBackgroundUploadChildRequests(ctx context.Context)
}
