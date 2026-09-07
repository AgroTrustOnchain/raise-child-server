package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IStaffService interface {
	GetStaff(id string, ctx context.Context) (response.StaffNftResponse, error)
	GetStaffByOwnerWallet(id string, ctx context.Context) (response.StaffNftResponse, error)
	GetStaffs(req request.GetStaffsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetStaffsV2(req request.GetStaffsRequest, ctx context.Context) (response.PaginationDataResponse, error)
}
