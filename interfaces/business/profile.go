package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IProfileService interface {
	UploadProfile(id string, req request.UploadProfileRequest, ctx context.Context) (response.PersonalProfileResponse, error)
}
