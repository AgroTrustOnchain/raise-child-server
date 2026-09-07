package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
)

type IProfileService interface {
	GetWalletPersonalProfile(id string, req request.GetTransactionRecordsRequest, ctx context.Context) (response.PersonalWalletProfileResponse, error)
	GetProfile(id string, ctx context.Context) (*entities.Profile, error)
	UploadProfile(id string, req request.UploadProfileRequest, ctx context.Context) (response.PersonalProfileResponse, error)
}
