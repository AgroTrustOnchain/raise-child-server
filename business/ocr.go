package business

import (
	"context"
	"log"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/util"
	"raise-child/util/ai"
)

type ocrService struct {
	ai        ai.IAiClientProvider
	errLogger *log.Logger
}

func initializeOCRService(ai ai.IAiClientProvider, errLogger *log.Logger) business.IOCRService {
	return &ocrService{
		ai:        ai,
		errLogger: errLogger,
	}
}

func GenerateOCRService() (business.IOCRService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	return initializeOCRService(
		ai.GetAiProvider(),
		errLogger,
	), nil
}

// GetExtractChildInfo implements [business.IOCRService].
func (o *ocrService) GetExtractChildInfo(req request.ExtractChildUploadInfoRequest, ctx context.Context) (response.ExtractChildUploadInfoResponse, error) {
	aiResp, err := o.ai.ExtractChildInfo(
		req.ChildBirthCertificateURL,
		req.FirstGuardianIDCardURL,
		req.SecondGuardianIDCardURL,
		ctx,
	)
	if err != nil {
		o.errLogger.Printf("OCR service: failed to extract child info: %v", err)
		return response.ExtractChildUploadInfoResponse{}, err
	}

	return *aiResp, nil
}
