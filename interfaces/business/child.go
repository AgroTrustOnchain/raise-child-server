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
}
