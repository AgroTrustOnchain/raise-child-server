package business

import (
	"context"
	"errors"
	"fmt"
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
	var repo = repository.InitializeBankProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
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

			mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.suiPaginatedRes, nil)
			repo.On("GetBankProfileByOwner", mock.AnythingOfType("string"), mock.Anything).Return(bankProfile, nil)
			repo.On("CreateBankProfile", mock.Anything, mock.Anything).Return(nil)
			_, err := service.CreateBankProfile(request.CreateBankProfileRequest{}, tc.context)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetBankProfile(t *testing.T) {
	var repo = repository.InitializeBankProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
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
	var adminCtx = context.WithValue(context.Background(), "address", "1")
	adminCtx = context.WithValue(adminCtx, "sub", "sub")
	var tcsInfo = []struct {
		isOwner     bool
		isBankExist bool
		expectedErr error
		context     context.Context
	}{
		{
			isOwner:     true,
			isBankExist: true,
			context:     ctx,
		},
		{
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{
			isBankExist: true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			context:     invalidCtx,
		},
		{
			isBankExist: true,
			context:     adminCtx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			repo.ExpectedCalls = nil
			var bankProfile *entities.BankProfile = nil
			if tc.isBankExist {
				bankProfile = &sampleUploadedBankProfile
			}

			repo.On("GetBankProfileById", mock.AnythingOfType("string"), mock.Anything).Return(bankProfile, nil)
			if !tc.isOwner {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: sampleJsonManageObj,
							},
						},
					},
				}, nil)
			}

			_, err := service.GetBankProfile("", tc.context)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetBankProfileByOwner(t *testing.T) {
	var repo = repository.InitializeBankProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
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
	var adminCtx = context.WithValue(context.Background(), "address", "1")
	adminCtx = context.WithValue(adminCtx, "sub", "sub")
	var tcsInfo = []struct {
		id          string
		isOwner     bool
		isBankExist bool
		expectedErr error
		context     context.Context
	}{
		{
			id:          "",
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{
			id:          sampleAddress,
			isOwner:     true,
			isBankExist: true,
			context:     ctx,
		},
		{
			id:          sampleAddress,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			context:     ctx,
		},
		{
			id:          sampleAddress,
			isBankExist: true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			context:     invalidCtx,
		},
		{
			id:          sampleAddress,
			isBankExist: true,
			context:     adminCtx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			repo.ExpectedCalls = nil
			var bankProfile *entities.BankProfile = nil
			if tc.isBankExist {
				bankProfile = &sampleUploadedBankProfile
			}

			repo.On("GetBankProfileByOwner", mock.AnythingOfType("string"), mock.Anything).Return(bankProfile, nil)
			if !tc.isOwner {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: sampleJsonManageObj,
							},
						},
					},
				}, nil)
			}

			_, err := service.GetBankProfileByOwner(tc.id, tc.context)
			assert.Equal(t, tc.expectedErr, err)
		})
	}

}
