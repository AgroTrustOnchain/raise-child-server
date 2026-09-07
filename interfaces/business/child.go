package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

// CreateBooksNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// 	CreateMealNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// 	CreateSpecialNeedWithdrawProposal(req request.CreateSpecialNeedWithdrawProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
//	ConfirmProvideMealForChild(id string, req request.ConfirmProvideMealForChildRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// ConfirmSpecialNeedProposal(id string, ctx context.Context) (response.BuildTransactionResponse, error)
// 	VoteSpecialNeedProposal(id string, req request.VoteRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// 	EditSpecialNeedDao(req request.EditDaoRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// 	UpdateBooksNeed(req request.UpdateChildNeedRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// 	UpdateMealNeed(req request.UpdateChildNeedRequest, ctx context.Context) (response.BuildTransactionResponse, error)
// 	UpdateHealthInsuranceNeed(req request.UpdateChildNeedRequest, ctx context.Context) (response.BuildTransactionResponse, error)

type IChildService interface {
	GetUserSupportedChildren(wallet string, req request.GetChildrenRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetChildren(req request.GetChildrenRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetChild(id string, ctx context.Context) (response.ChildResponse, error)
	UploadChild(req request.UploadChildRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	AddStringMetadata(id string, req request.AddChildStringMetadataRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	AddNumberMetadata(id string, req request.AddChildNumberMetadataRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	SupportBooksNeed(id string, ctx context.Context) (response.PaymentUrlResponse, error)
	SupportHealthInsuranceNeed(id string, ctx context.Context) (response.PaymentUrlResponse, error)
	SupportMealNeed(id string, req request.SupportMealNeadRequest, ctx context.Context) (response.PaymentUrlResponse, error)
	SupportSpecialNeed(id string, req request.SupportSpecialNeedRequest, ctx context.Context) (response.PaymentUrlResponse, error)
	ConfirmProvideMealForChild(id string, req request.ConfirmProvideMealForChildRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	CreateBooksNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) error
	CreateMealNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) error
	CreateSpecialNeedWithdrawProposal(req request.CreateSpecialNeedWithdrawProposalRequest, ctx context.Context) error
	CreateHealthInsuranceNeedWithdrawProposal(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) error
	CreateBooksNeedWithdrawProposalV2(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	CreateMealNeedWithdrawProposalV2(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	CreateSpecialNeedWithdrawProposalV2(req request.CreateSpecialNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	CreateHealthInsuranceNeedWithdrawProposalV2(req request.CreateNormalNeedWithdrawProposalRequest, ctx context.Context) (*entities.PendingWithdrawProposal, error)
	CreateSpecialNeedProposal(req request.CreateSpecialNeedProposalRequest, ctx context.Context) (response.BuildTransactionResponse, error)
	CreateSpecialNeedProposalV2(req request.CreateSpecialNeedProposalRequest, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error)
	ConfirmSpecialNeedProposal(id string, ctx context.Context) error
	VoteSpecialNeedProposal(id string, req request.VoteRequest, ctx context.Context) error
	UpdateBooksNeed(req request.UpdateChildNeedRequest, ctx context.Context) error
	UpdateMealNeed(req request.UpdateChildNeedRequest, ctx context.Context) error
	UpdateHealthInsuranceNeed(req request.UpdateChildNeedRequest, ctx context.Context) error
}
