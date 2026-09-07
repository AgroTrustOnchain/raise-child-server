package business

import (
	"context"
	"raise-child/model/dtos/response"
)

type IImageService interface {
	GetPresignedUrl(ctx context.Context) (response.PresignedUrlResponse, error)
}