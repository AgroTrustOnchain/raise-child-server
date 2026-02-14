package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IChildService interface {
	GetChildren(req request.GetChildrenRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetChild(id string, ctx context.Context) (response.ChildResponse, error)
	UploadChild(req request.UploadChildRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	AddStringMetada(id string, req request.AddChildStringMetadaRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	AddNumberMetada(id string, req request.AddChildNumberMetadaRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	SupportBooksNeed(id string, ctx context.Context) (response.UrlAPIResponse, error)
	SupportMealNeed(id string, req request.SupportMealNeadRequest, ctx context.Context) (response.UrlAPIResponse, error)
	SupportSpecialNeed(id string, req request.SupportSpecialNeedRequest, ctx context.Context) (response.UrlAPIResponse, error)
	CreateBooksNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	CreateMealNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	CreateSpecialNeedWithdrawProposal(req request.CreateSpecialNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	CreateSpecialNeedProposal(req request.CreateSpecialNeedProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	ConfirmSpecialNeedProposal(id string, ctx context.Context) (response.BuildTransactionResponse, error)
	VoteSpecialNeedProposal(id string, req request.VoteRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	EditSpecialNeedDao(req request.EditDaoRequest, ctx context.Context) (response.BuildTransactionResponse, error)
}
