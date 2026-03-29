package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IProfileService interface {
	GetWalletPersonalProfile(id string, req request.GetTransactionRecordsRequest, ctx context.Context) (response.PersonalWalletProfileResponse, error)
	UploadProfile(id string, req request.UploadProfileRequest, ctx context.Context) (response.PersonalProfileResponse, error)
}
