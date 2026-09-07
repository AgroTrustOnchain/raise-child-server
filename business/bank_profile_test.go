package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
	"raise-child/mocks/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
	"raise-child/util"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBankProfile(t *testing.T) {
	var repo = repository.InializeBankProfileMockRepo()
	var mockClient = pkg.InializeSuiMockApi()
	var service = initializeBankProfileService(
		repo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var invalidCtx = context.WithValue(context.Background(), "address", sampleInvalidAddress)
	invalidCtx = context.WithValue(invalidCtx, "sub", "sub")
	var tcsInfo = []struct {
		suiPaginatedRes models.PaginatedObjectsResponse
		isAlreadyUpload bool
		expectedErr     error
		context         context.Context
	}{
		{ // Not valid address case
			context:     invalidCtx,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Happy case
			suiPaginatedRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonLeaderNft1,
								},
							},
						},
					},
				},
			},
			context: ctx,
		},
		{ // Already upload bank profile
			suiPaginatedRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonLeaderNft1,
								},
							},
						},
					},
				},
			},
			isAlreadyUpload: true,
			expectedErr:     errors.New(noti.BANK_PROFILE_EXISTED_MESSAGE),
			context:         ctx,
		},
		{ // Not leader case
			suiPaginatedRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonVolunteerNft1,
								},
							},
						},
					},
				},
			},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			context:     ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			repo.ExpectedCalls = nil

			var bankProfile *entities.BankProfile = nil
			if tc.isAlreadyUpload {
				bankProfile = &sampleUploadedBankProfile
			}

			log.Printf("\nAddress in case %d: %s", i, tc.context.Value("address").(string))
			log.Printf("\nSub in case %d: %s", i, tc.context.Value("sub").(string))

			mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.suiPaginatedRes, nil)
			repo.On("GetBankProfileByOwner", mock.AnythingOfType("string"), mock.Anything).Return(bankProfile, nil)
			repo.On("CreateBankProfile", mock.Anything, mock.Anything).Return(nil)
			_, err := service.CreateBankProfile(request.CreateBankProfileRequest{}, tc.context)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
