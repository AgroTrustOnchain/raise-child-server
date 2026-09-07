package business

import (
	"context"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/response"
	"raise-child/util"
	"raise-child/util/image/cloudinary"
	"time"
)

type imageService struct {
	errLogger *log.Logger
}

func initializeImageService(errLogger *log.Logger) business.IImageService {
	return &imageService{
		errLogger: errLogger,
	}
}

func GenerateImageService() business.IImageService {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)
	return initializeImageService(errLogger)
}

// PresignedUrl implements business.IImageService.
func (i *imageService) GetPresignedUrl(ctx context.Context) (response.PresignedUrlResponse, error) {
	timestamp := time.Now().Add(time.Hour).Unix()
	folder := os.Getenv(env.CLOUDINARY_STORAGE_FOLDER)
	signature := cloudinary.GenerateSignature(timestamp, folder)
	apiKey := os.Getenv(env.CLOUDINARY_API_KEY)
	uploadUrl := cloudinary.GetUploadUrl()

	return response.PresignedUrlResponse{
		Timestamp: timestamp,
		Folder:    folder,
		Signature: signature,
		ApiKey:    apiKey,
		UploadUrl: uploadUrl,
	}, nil
}
