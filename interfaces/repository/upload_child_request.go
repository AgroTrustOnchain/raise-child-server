package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
)

type IUploadChildRequestRepository interface {
	GetUploadChildRequest(id string, ctx context.Context) (*entities.UploadChildRequest, error)
	GetUploadChildRequests(req request.GetUploadChildRequests, ctx context.Context) ([]entities.UploadChildRequest, int, error)
	GetWalletUploadChildRequests(id string, page int, ctx context.Context) ([]entities.UploadChildRequest, int, error)
	CreateUploadChildRequest(req entities.UploadChildRequest, ctx context.Context) error
	UpdateUploadChildRequest(req entities.UploadChildRequest, ctx context.Context) error
	IsChildRequested(identityCode string, ctx context.Context) (bool, error)
}
