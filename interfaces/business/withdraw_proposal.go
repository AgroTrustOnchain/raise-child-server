package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IWithdrawProposalService interface {
	GetWithdrawProposal(id string, ctx context.Context) (response.WithdrawProposalResponse, error)
	GetWithdrawProposals(req request.GetWithdrawProposalsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	CreateWithdrawProposal(req request.CreateWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	VoteWithdrawProposal(id string, req request.VoteRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	ConfirmWithdrawProposal(id string, ctx context.Context) (map[string]interface{}, error)
	ConfirmMainPoolWithdrawProposal(id, capturedImgBlobId string, ctx context.Context) (response.BuildTransactionResponse, error)
}
