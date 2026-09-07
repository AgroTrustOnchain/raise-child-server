package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type IOCRService interface {
	GetExtractChildInfo(request.ExtractChildUploadInfoRequest, context.Context) (response.ExtractChildUploadInfoResponse, error)
}
