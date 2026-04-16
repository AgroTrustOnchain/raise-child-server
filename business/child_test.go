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
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetChildren(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var fullJsons = getFullJsonChildren()
	var exepectedFullJsonData []*models.SuiObjectResponse
	for _, json := range fullJsons {
		exepectedFullJsonData = append(exepectedFullJsonData, &models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: json,
					},
				},
			},
		})
	}

	keywordData, keyword := getFoundJsonChildrenWithKeyWord()
	var exepectedDataWithKw []*models.SuiObjectResponse
	for _, json := range keywordData {
		exepectedDataWithKw = append(exepectedDataWithKw, &models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: json,
					},
				},
			},
		})
	}

	var ctx = context.Background()
	var tcsInfo = []struct {
		suiJsonData []*models.SuiObjectResponse
		keyword     string
		page        int
		isEmpty     bool
	}{
		{
			suiJsonData: exepectedFullJsonData,
		},
		{
			suiJsonData: exepectedFullJsonData,
			page:        2,
			isEmpty:     true,
		},
		{
			suiJsonData: exepectedDataWithKw,
			keyword:     keyword,
		},
	}

	var manageJsonObj = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonManageObj,
				},
			},
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(manageJsonObj, nil)
			mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return(tc.suiJsonData, nil)
			res, err := service.GetChildren(request.GetChildrenRequest{
				Keyword: tc.keyword,
				Page:    tc.page,
			}, ctx)

			assert.NoError(t, err)
			if tc.isEmpty {
				assert.Empty(t, res.Data)
			} else {
				assert.Equal(t, len(tc.suiJsonData), res.Amount)
			}
		})
	}
}

func TestGetChild(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.Background()
	var tcsInfo = []struct {
		id          string
		suiRes      models.SuiObjectResponse
		expectedErr error
	}{
		{ // Happy case
			id: sampleAddress,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonChild1,
						},
					},
				},
			},
		},
		{ // Invalid id case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			if tc.suiRes.Data != nil {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.suiRes, nil)
			}
		})

		_, err := service.GetChild(tc.id, ctx)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestSupportBooksNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var tcsInfo = []struct {
		needId               string
		profile              *entities.Profile
		getProfileErr        error
		suiRes               models.SuiObjectResponse
		isCallCreateDonation bool
		createDonateRes      error
		isCallCreatePayment  bool
		createPaymentRes     error
		expectedErr          error
	}{
		// Can't execute case 1 with load .env
		// { // Happy case
		// 	needId:  sampleAddress,
		// 	profile: &sampleProfileObj,
		// 	suiRes: models.SuiObjectResponse{
		// 		Data: &models.SuiObjectData{
		// 			Content: &models.SuiParsedData{
		// 				SuiMoveObject: models.SuiMoveObject{
		// 					Fields: sampleNotSupportedBooksNeedJson,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	isCallCreateDonation: true,
		// 	isCallCreatePayment:  true,
		// },
		{ // Error get profile case
			getProfileErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:   errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ //
			expectedErr: errors.New(noti.PROFILE_EMPTY_MESSAGE),
		},
		{ // Invalid need id case
			profile:     &sampleProfileObj,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Need supported case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleSupportedBooksNeedJson,
						},
					},
				},
			},
			expectedErr: errors.New(noti.NEED_SUPPORTED_MESSAGE),
		},
		{ // Create donation fail case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleNotSupportedBooksNeedJson,
						},
					},
				},
			},
			isCallCreateDonation: true,
			createDonateRes:      errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Create payment fail case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleNotSupportedBooksNeedJson,
						},
					},
				},
			},
			isCallCreateDonation: true,
			isCallCreatePayment:  true,
			createPaymentRes:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			profileRepo.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			offchainDonationRepo.ExpectedCalls = nil
			paymentRepo.ExpectedCalls = nil

			profileRepo.On("GetProfile", mock.Anything, mock.Anything).Return(tc.profile, tc.getProfileErr)
			if tc.suiRes.Data != nil {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.suiRes, nil)
			}

			if tc.isCallCreateDonation {
				offchainDonationRepo.On("CreateDonation", mock.Anything, mock.Anything).Return(tc.createDonateRes)
			}

			if tc.isCallCreatePayment {
				paymentRepo.On("CreatePayment", mock.Anything, mock.Anything).Return(tc.createPaymentRes)
			}

			_, err := service.SupportBooksNeed(tc.needId, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSupportHealthInsuranceNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var aiClient = pkg.InitializeAiMockClient()
	var walrusProvider = pkg.InitializeWalrusMockProvider()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		aiClient,
		walrusProvider,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var tcsInfo = []struct {
		needId               string
		profile              *entities.Profile
		getProfileErr        error
		suiRes               models.SuiObjectResponse
		isCallCreateDonation bool
		createDonateRes      error
		isCallCreatePayment  bool
		createPaymentRes     error
		expectedErr          error
	}{
		// Can't execute case 1 with load .env
		// { // Happy case
		// 	needId:  sampleAddress,
		// 	profile: &sampleProfileObj,
		// 	suiRes: models.SuiObjectResponse{
		// 		Data: &models.SuiObjectData{
		// 			Content: &models.SuiParsedData{
		// 				SuiMoveObject: models.SuiMoveObject{
		// 					Fields: sampleNotSupportedBooksNeedJson,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	isCallCreateDonation: true,
		// 	isCallCreatePayment:  true,
		// },
		{ // Error get profile case
			getProfileErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:   errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ //
			expectedErr: errors.New(noti.PROFILE_EMPTY_MESSAGE),
		},
		{ // Invalid need id case
			profile:     &sampleProfileObj,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Need supported case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleSupportedBooksNeedJson,
						},
					},
				},
			},
			expectedErr: errors.New(noti.NEED_SUPPORTED_MESSAGE),
		},
		{ // Create donation fail case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleNotSupportedBooksNeedJson,
						},
					},
				},
			},
			isCallCreateDonation: true,
			createDonateRes:      errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Create payment fail case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleNotSupportedBooksNeedJson,
						},
					},
				},
			},
			isCallCreateDonation: true,
			isCallCreatePayment:  true,
			createPaymentRes:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			profileRepo.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			offchainDonationRepo.ExpectedCalls = nil
			paymentRepo.ExpectedCalls = nil

			profileRepo.On("GetProfile", mock.Anything, mock.Anything).Return(tc.profile, tc.getProfileErr)
			if tc.suiRes.Data != nil {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.suiRes, nil)
			}

			if tc.isCallCreateDonation {
				offchainDonationRepo.On("CreateDonation", mock.Anything, mock.Anything).Return(tc.createDonateRes)
			}

			if tc.isCallCreatePayment {
				paymentRepo.On("CreatePayment", mock.Anything, mock.Anything).Return(tc.createPaymentRes)
			}

			_, err := service.SupportHealthInsuranceNeed(tc.needId, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSupportMealNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var mealSupportDurationRepo = repository.InitializeMealSupportDurationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		mealSupportDurationRepo,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var tcsInfo = []struct {
		needId                          string
		req                             request.SupportMealNeadRequest
		profile                         *entities.Profile
		getProfileErr                   error
		suiRes                          models.SuiObjectResponse
		isCallCreateMealSupportDuration bool
		createMealSupportDurationRes    error
		isCallCreateDonation            bool
		createDonateRes                 error
		isCallCreatePayment             bool
		createPaymentRes                error
		expectedErr                     error
	}{
		{ // Error get profile case
			getProfileErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:   errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ //
			expectedErr: errors.New(noti.PROFILE_EMPTY_MESSAGE),
		},
		{ // Invalid need id case
			profile:     &sampleProfileObj,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Support out of range case
			needId: sampleAddress,
			req: request.SupportMealNeadRequest{
				Months: 2,
			},
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleMealNeedJson,
						},
					},
				},
			},
			expectedErr: errors.New(noti.MEAL_NEED_SUPPORT_DURATION_OUT_RANGE_MESSAGE),
		},
		{ // Create meal support duration fail case
			needId: sampleAddress,
			req: request.SupportMealNeadRequest{
				Months: 1,
			},
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleMealNeedJson,
						},
					},
				},
			},
			isCallCreateMealSupportDuration: true,
			createMealSupportDurationRes:    errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:                     errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Create donation fail case
			needId: sampleAddress,
			req: request.SupportMealNeadRequest{
				Months: 1,
			},
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleMealNeedJson,
						},
					},
				},
			},
			isCallCreateMealSupportDuration: true,
			isCallCreateDonation:            true,
			createDonateRes:                 errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:                     errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Create payment fail case
			needId: sampleAddress,
			req: request.SupportMealNeadRequest{
				Months: 1,
			},
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleMealNeedJson,
						},
					},
				},
			},
			isCallCreateMealSupportDuration: true,
			isCallCreateDonation:            true,
			isCallCreatePayment:             true,
			createPaymentRes:                errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:                     errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			profileRepo.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			offchainDonationRepo.ExpectedCalls = nil
			mealSupportDurationRepo.ExpectedCalls = nil
			paymentRepo.ExpectedCalls = nil

			profileRepo.On("GetProfile", mock.Anything, mock.Anything).Return(tc.profile, tc.getProfileErr)
			if tc.suiRes.Data != nil {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.suiRes, nil)
			}

			if tc.isCallCreateMealSupportDuration {
				mealSupportDurationRepo.On("CreateMealSupportDuration", mock.Anything, mock.Anything).Return(tc.createMealSupportDurationRes)
			}

			if tc.isCallCreateDonation {
				offchainDonationRepo.On("CreateDonation", mock.Anything, mock.Anything).Return(tc.createDonateRes)
			}

			if tc.isCallCreatePayment {
				paymentRepo.On("CreatePayment", mock.Anything, mock.Anything).Return(tc.createPaymentRes)
			}

			_, err := service.SupportMealNeed(tc.needId, tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSupportSpecialNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var tcsInfo = []struct {
		needId               string
		profile              *entities.Profile
		getProfileErr        error
		req                  request.SupportSpecialNeedRequest
		suiRes               models.SuiObjectResponse
		isCallCreateDonation bool
		createDonateRes      error
		isCallCreatePayment  bool
		createPaymentRes     error
		expectedErr          error
	}{
		{ // Error get profile case
			getProfileErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:   errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ //
			expectedErr: errors.New(noti.PROFILE_EMPTY_MESSAGE),
		},
		{ // Invalid need id case
			profile:     &sampleProfileObj,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Support out of target
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			req: request.SupportSpecialNeedRequest{
				Amount: 150000,
			},
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleSpecialNeedCampaignJson,
						},
					},
				},
			},
			expectedErr: errors.New(noti.SUPPORT_SURPASS_CAMPAIGN_TARGET_MESSAGE),
		},
		{ // Create donation fail case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			req: request.SupportSpecialNeedRequest{
				Amount: 100000,
			},
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleSpecialNeedCampaignJson,
						},
					},
				},
			},
			isCallCreateDonation: true,
			createDonateRes:      errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Create payment fail case
			needId:  sampleAddress,
			profile: &sampleProfileObj,
			req: request.SupportSpecialNeedRequest{
				Amount: 100000,
			},
			suiRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleSpecialNeedCampaignJson,
						},
					},
				},
			},
			isCallCreateDonation: true,
			isCallCreatePayment:  true,
			createPaymentRes:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			profileRepo.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			offchainDonationRepo.ExpectedCalls = nil
			paymentRepo.ExpectedCalls = nil

			profileRepo.On("GetProfile", mock.Anything, mock.Anything).Return(tc.profile, tc.getProfileErr)
			if tc.suiRes.Data != nil {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.suiRes, nil)
			}

			if tc.isCallCreateDonation {
				offchainDonationRepo.On("CreateDonation", mock.Anything, mock.Anything).Return(tc.createDonateRes)
			}

			if tc.isCallCreatePayment {
				paymentRepo.On("CreatePayment", mock.Anything, mock.Anything).Return(tc.createPaymentRes)
			}

			_, err := service.SupportSpecialNeed(tc.needId, tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestCreateBooksNeedWithdrawProposalV2(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var aiClient = pkg.InitializeAiMockClient()
	var walrusProvider = pkg.InitializeWalrusMockProvider()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		aiClient,
		walrusProvider,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validStaffsRes = models.PaginatedObjectsResponse{
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
	}

	var sameRegionChildRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var validNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleStillNotWithdrawnBooksNeedJson,
				},
			},
		},
	}

	var curTime time.Time = time.Now()
	var notiObj = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime)},
		Contents:                []string{"1", "2"},
	}

	var passedWithdrawDateNoti = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime.AddDate(0, -1, 0))},
		Contents:                []string{"1", "2"},
	}

	var stillNotComeWithdrawDateNoti = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime.AddDate(0, 1, 0))},
		Contents:                []string{"1", "2"},
	}

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var sampleBlobId string = "blobID"
	var tcsInfo = []struct {
		req                  request.CreateNormalNeedWithdrawProposalRequest
		needRes              models.SuiObjectResponse
		isGetNeed            bool
		childRes             models.SuiObjectResponse
		isGetChild           bool
		staffNftsRes         models.PaginatedObjectsResponse
		isStaffChecked       bool
		noti                 *entities.LeaderNoti
		isGetNoti            bool
		getNotiErr           error
		isGetProposal        bool
		isProposalProposed   bool
		getIsProposedErr     error
		isGetPools           bool
		isCallCreateProposal bool
		createProposalRes    error
		isCallAiValidate     bool
		expectedErr          error
	}{
		{ // Happy case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID:      sampleAddress,
				ProofBlobID: &sampleBlobId,
			},
			needRes:              validNeedRes,
			isGetNeed:            true,
			childRes:             sameRegionChildRes,
			isGetChild:           true,
			staffNftsRes:         validStaffsRes,
			isStaffChecked:       true,
			noti:                 &notiObj,
			isGetNoti:            true,
			isGetProposal:        true,
			isGetPools:           true,
			isCallCreateProposal: true,
			isCallAiValidate:     true,
		},
		{ // Invalid need case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetNeed:   true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Already withdrawn case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleWithdrawnBooksNeedJson,
						},
					},
				},
			},
			isGetNeed:   true,
			expectedErr: errors.New(noti.NEED_WITHDRAWN_MESSAGE),
		},
		{ // Not exist child case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:   validNeedRes,
			isGetNeed: true,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetChild:  true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not staff case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:    validNeedRes,
			isGetNeed:  true,
			childRes:   sameRegionChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isStaffChecked: true,
			expectedErr:    errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Different region leader case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:   validNeedRes,
			isGetNeed: true,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonChild10,
						},
					},
				},
			},
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			expectedErr:    errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Get noti fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			isGetNoti:      true,
			getNotiErr:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:    errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Still not withdraw date case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			noti:           &stillNotComeWithdrawDateNoti,
			isGetNoti:      true,
			expectedErr:    errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE),
		},
		{ // Passed withdraw date case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			noti:           &passedWithdrawDateNoti,
			isGetNoti:      true,
			expectedErr:    errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE),
		},
		{ // Check proposal fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:          validNeedRes,
			isGetNeed:        true,
			childRes:         sameRegionChildRes,
			isGetChild:       true,
			staffNftsRes:     validStaffsRes,
			isStaffChecked:   true,
			noti:             &notiObj,
			isGetNoti:        true,
			isGetProposal:    true,
			getIsProposedErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:      errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Proposal proposed case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:            validNeedRes,
			isGetNeed:          true,
			childRes:           sameRegionChildRes,
			isGetChild:         true,
			staffNftsRes:       validStaffsRes,
			isStaffChecked:     true,
			noti:               &notiObj,
			isGetNoti:          true,
			isGetProposal:      true,
			isProposalProposed: true,
			expectedErr:        errors.New(noti.STILL_PENDING_REQUEST_MESSAGE),
		},
		{ // Create pending proposal fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID:      sampleAddress,
				ProofBlobID: &sampleBlobId,
			},
			needRes:              validNeedRes,
			isGetNeed:            true,
			childRes:             sameRegionChildRes,
			isGetChild:           true,
			staffNftsRes:         validStaffsRes,
			isStaffChecked:       true,
			noti:                 &notiObj,
			isGetNoti:            true,
			isGetProposal:        true,
			isGetPools:           true,
			isCallCreateProposal: true,
			isCallAiValidate:     true,
			createProposalRes:    errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			leaderNotiRepo.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil
			walrusProvider.ExpectedCalls = nil
			aiClient.ExpectedCalls = nil

			if tc.isGetNeed {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.needRes, nil).Once()
			}

			if tc.isGetChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}

			if tc.isStaffChecked {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.staffNftsRes, nil).Once()
			}

			if tc.isGetNoti {
				leaderNotiRepo.On("GetNotiByNeed", mock.Anything, mock.Anything).Return(tc.noti, tc.getNotiErr)
			}

			if tc.isGetProposal {
				pendingWithdrawProposalRepo.On("IsPendingWithdrawProposalProposedWithSpecificInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.isProposalProposed, tc.getIsProposedErr)
			}

			if tc.isCallCreateProposal {
				pendingWithdrawProposalRepo.On("CreatePendingWithdrawProposal", mock.Anything, mock.Anything).Return(tc.createProposalRes)
			}

			if tc.isGetPools {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: samplePoolObj,
							},
						},
					},
				}, nil).Once()

				mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return([]*models.SuiObjectResponse{
					&models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: samplePoolObj,
								},
							},
						},
					},
				}, nil).Once()
			}

			if tc.isCallAiValidate {
				walrusProvider.On("FetchBytesImage", mock.Anything).Return([]byte{}, nil)
				aiClient.On("ValidateWithdrawProposal", mock.Anything, mock.Anything).Return("")
			}

			res, err := service.CreateBooksNeedWithdrawProposalV2(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.True(t, res != nil)
			}
		})
	}
}

func TestCreateHealthInsuranceNeedWithdrawProposalV2(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var aiClient = pkg.InitializeAiMockClient()
	var walrusProvider = pkg.InitializeWalrusMockProvider()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		aiClient,
		walrusProvider,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validStaffsRes = models.PaginatedObjectsResponse{
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
	}

	var sameRegionChildRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var validNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleStillNotWithdrawnBooksNeedJson,
				},
			},
		},
	}

	var curTime time.Time = time.Now()
	var notiObj = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime)},
		Contents:                []string{"1", "2"},
	}

	var passedWithdrawDateNoti = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime.AddDate(0, -1, 0))},
		Contents:                []string{"1", "2"},
	}

	var stillNotComeWithdrawDateNoti = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime.AddDate(0, 1, 0))},
		Contents:                []string{"1", "2"},
	}

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var sampleBlobId string = "blobID"
	var tcsInfo = []struct {
		req                  request.CreateNormalNeedWithdrawProposalRequest
		needRes              models.SuiObjectResponse
		isGetNeed            bool
		childRes             models.SuiObjectResponse
		isGetChild           bool
		staffNftsRes         models.PaginatedObjectsResponse
		isStaffChecked       bool
		noti                 *entities.LeaderNoti
		isGetNoti            bool
		getNotiErr           error
		isGetProposal        bool
		isProposalProposed   bool
		getIsProposedErr     error
		isGetPools           bool
		isCallCreateProposal bool
		createProposalRes    error
		isCallAiValidate     bool
		expectedErr          error
	}{
		{ // Happy case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID:      sampleAddress,
				ProofBlobID: &sampleBlobId,
			},
			needRes:              validNeedRes,
			isGetNeed:            true,
			childRes:             sameRegionChildRes,
			isGetChild:           true,
			staffNftsRes:         validStaffsRes,
			isStaffChecked:       true,
			noti:                 &notiObj,
			isGetNoti:            true,
			isGetProposal:        true,
			isGetPools:           true,
			isCallCreateProposal: true,
			isCallAiValidate:     true,
		},
		{ // Invalid need case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetNeed:   true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Already withdrawn case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleWithdrawnBooksNeedJson,
						},
					},
				},
			},
			isGetNeed:   true,
			expectedErr: errors.New(noti.NEED_WITHDRAWN_MESSAGE),
		},
		{ // Not exist child case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:   validNeedRes,
			isGetNeed: true,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetChild:  true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not staff case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:    validNeedRes,
			isGetNeed:  true,
			childRes:   sameRegionChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isStaffChecked: true,
			expectedErr:    errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Different region leader case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:   validNeedRes,
			isGetNeed: true,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonChild10,
						},
					},
				},
			},
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			expectedErr:    errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Get noti fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			isGetNoti:      true,
			getNotiErr:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:    errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Still not withdraw date case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			noti:           &stillNotComeWithdrawDateNoti,
			isGetNoti:      true,
			expectedErr:    errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE),
		},
		{ // Passed withdraw date case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			noti:           &passedWithdrawDateNoti,
			isGetNoti:      true,
			expectedErr:    errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE),
		},
		{ // Check proposal fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:          validNeedRes,
			isGetNeed:        true,
			childRes:         sameRegionChildRes,
			isGetChild:       true,
			staffNftsRes:     validStaffsRes,
			isStaffChecked:   true,
			noti:             &notiObj,
			isGetNoti:        true,
			isGetProposal:    true,
			getIsProposedErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:      errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Proposal proposed case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:            validNeedRes,
			isGetNeed:          true,
			childRes:           sameRegionChildRes,
			isGetChild:         true,
			staffNftsRes:       validStaffsRes,
			isStaffChecked:     true,
			noti:               &notiObj,
			isGetNoti:          true,
			isGetProposal:      true,
			isProposalProposed: true,
			expectedErr:        errors.New(noti.STILL_PENDING_REQUEST_MESSAGE),
		},
		{ // Create pending proposal fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID:      sampleAddress,
				ProofBlobID: &sampleBlobId,
			},
			needRes:              validNeedRes,
			isGetNeed:            true,
			childRes:             sameRegionChildRes,
			isGetChild:           true,
			staffNftsRes:         validStaffsRes,
			isStaffChecked:       true,
			noti:                 &notiObj,
			isGetNoti:            true,
			isGetProposal:        true,
			isGetPools:           true,
			isCallCreateProposal: true,
			createProposalRes:    errors.New(noti.INTERNALL_ERR_MSG),
			isCallAiValidate:     true,
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			leaderNotiRepo.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil
			walrusProvider.ExpectedCalls = nil
			aiClient.ExpectedCalls = nil

			if tc.isGetNeed {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.needRes, nil).Once()
			}

			if tc.isGetChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}

			if tc.isStaffChecked {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.staffNftsRes, nil).Once()
			}

			if tc.isGetNoti {
				leaderNotiRepo.On("GetNotiByNeed", mock.Anything, mock.Anything).Return(tc.noti, tc.getNotiErr)
			}

			if tc.isGetProposal {
				pendingWithdrawProposalRepo.On("IsPendingWithdrawProposalProposedWithSpecificInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.isProposalProposed, tc.getIsProposedErr)
			}

			if tc.isCallCreateProposal {
				pendingWithdrawProposalRepo.On("CreatePendingWithdrawProposal", mock.Anything, mock.Anything).Return(tc.createProposalRes)
			}

			if tc.isGetPools {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: samplePoolObj,
							},
						},
					},
				}, nil).Once()

				mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return([]*models.SuiObjectResponse{
					&models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: samplePoolObj,
								},
							},
						},
					},
				}, nil).Once()
			}

			if tc.isCallAiValidate {
				walrusProvider.On("FetchBytesImage", mock.Anything).Return([]byte{}, nil)
				aiClient.On("ValidateWithdrawProposal", mock.Anything, mock.Anything).Return("")
			}

			res, err := service.CreateHealthInsuranceNeedWithdrawProposalV2(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.True(t, res != nil)
			}
		})
	}
}

func TestCreateMealNeedWithdrawProposalV2(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var aiClient = pkg.InitializeAiMockClient()
	var walrusProvider = pkg.InitializeWalrusMockProvider()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		aiClient,
		walrusProvider,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validStaffsRes = models.PaginatedObjectsResponse{
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
	}

	var sameRegionChildRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var curTime time.Time = time.Now()
	var validEndPeriod time.Time
	var monthsAdded int = 1
	var upperBoundPeriod time.Time = util.ToEndOfDate(util.RawDateToTime(fmt.Sprintf("15/01/%d", (curTime.Year() + 1))))
	for {
		var tmpTime time.Time = curTime.AddDate(0, monthsAdded, 0)
		if tmpTime.After(upperBoundPeriod) {
			validEndPeriod = tmpTime.AddDate(0, -1, 0)
			break
		}

		monthsAdded += 1

	}

	var endMonth int = int(validEndPeriod.Month())
	if endMonth == 1 {
		endMonth = 13
	}

	var totalSupportedMonths int = endMonth - int(curTime.Month())

	var validNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"durations": []map[string]interface{}{
							{
								"fields": map[string]interface{}{
									"start_period": util.TimeToRawDate(curTime),
									"end_period":   util.TimeToRawDate(validEndPeriod),
								},
							},
						},
						"total_supported_months": fmt.Sprintf("%d", totalSupportedMonths),
					},
				},
			},
		},
	}

	var stillNotWithdrawDateNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"durations": []map[string]interface{}{
							{
								"fields": map[string]interface{}{
									"start_period": util.TimeToRawDate(curTime.AddDate(0, 0, -7)),
									"end_period":   util.TimeToRawDate(validEndPeriod.AddDate(0, 0, -7)),
								},
							},
						},
						"total_supported_months": fmt.Sprintf("%d", totalSupportedMonths),
					},
				},
			},
		},
	}

	var withdrawnNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"total_supported_months": "1",
						"withdraws_for_need":     []string{"1"},
					},
				},
			},
		},
	}

	var notiObj = entities.LeaderNoti{
		ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime)},
		Contents:                make([]string, totalSupportedMonths),
	}

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var sampleBlobId string = "blobID"
	var tcsInfo = []struct {
		req                  request.CreateNormalNeedWithdrawProposalRequest
		needRes              models.SuiObjectResponse
		isGetNeed            bool
		childRes             models.SuiObjectResponse
		isGetChild           bool
		staffNftsRes         models.PaginatedObjectsResponse
		isStaffChecked       bool
		noti                 *entities.LeaderNoti
		isGetNoti            bool
		getNotiErr           error
		isGetProposal        bool
		isProposalProposed   bool
		getIsProposedErr     error
		isGetPools           bool
		isCallCreateProposal bool
		createProposalRes    error
		isCallAiValidate     bool
		expectedErr          error
	}{
		{ // Happy case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID:      sampleAddress,
				ProofBlobID: &sampleBlobId,
			},
			needRes:              validNeedRes,
			isGetNeed:            true,
			childRes:             sameRegionChildRes,
			isGetChild:           true,
			staffNftsRes:         validStaffsRes,
			isStaffChecked:       true,
			noti:                 &notiObj,
			isGetNoti:            true,
			isGetProposal:        true,
			isGetPools:           true,
			isCallCreateProposal: true,
			isCallAiValidate:     true,
		},
		{ // Invalid need case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist need case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetNeed:   true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Already withdrawn case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:     withdrawnNeedRes,
			isGetNeed:   true,
			expectedErr: errors.New(noti.NEED_WITHDRAWN_MESSAGE),
		},
		{ // Not date to withdraw case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:     stillNotWithdrawDateNeedRes,
			isGetNeed:   true,
			expectedErr: errors.New(noti.NOT_WITHDRAW_EXPECTED_DATE_MESSAGE),
		},
		{ // Not exist child case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:   validNeedRes,
			isGetNeed: true,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetChild:  true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not staff case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:    validNeedRes,
			isGetNeed:  true,
			childRes:   sameRegionChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isStaffChecked: true,
			expectedErr:    errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Different region leader case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:   validNeedRes,
			isGetNeed: true,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonChild10,
						},
					},
				},
			},
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			expectedErr:    errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Get noti fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:        validNeedRes,
			isGetNeed:      true,
			childRes:       sameRegionChildRes,
			isGetChild:     true,
			staffNftsRes:   validStaffsRes,
			isStaffChecked: true,
			isGetNoti:      true,
			getNotiErr:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:    errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Check proposal fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:          validNeedRes,
			isGetNeed:        true,
			childRes:         sameRegionChildRes,
			isGetChild:       true,
			staffNftsRes:     validStaffsRes,
			isStaffChecked:   true,
			noti:             &notiObj,
			isGetNoti:        true,
			isGetProposal:    true,
			getIsProposedErr: errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:      errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // Proposal proposed case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
			},
			needRes:            validNeedRes,
			isGetNeed:          true,
			childRes:           sameRegionChildRes,
			isGetChild:         true,
			staffNftsRes:       validStaffsRes,
			isStaffChecked:     true,
			noti:               &notiObj,
			isGetNoti:          true,
			isGetProposal:      true,
			isProposalProposed: true,
			expectedErr:        errors.New(noti.STILL_PENDING_REQUEST_MESSAGE),
		},
		{ // Create pending proposal fail case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID:      sampleAddress,
				ProofBlobID: &sampleBlobId,
			},
			needRes:              validNeedRes,
			isGetNeed:            true,
			childRes:             sameRegionChildRes,
			isGetChild:           true,
			staffNftsRes:         validStaffsRes,
			isStaffChecked:       true,
			noti:                 &notiObj,
			isGetNoti:            true,
			isGetProposal:        true,
			isGetPools:           true,
			isCallCreateProposal: true,
			createProposalRes:    errors.New(noti.INTERNALL_ERR_MSG),
			isCallAiValidate:     true,
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			leaderNotiRepo.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil
			walrusProvider.ExpectedCalls = nil
			aiClient.ExpectedCalls = nil

			if tc.isGetNeed {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.needRes, nil).Once()
			}

			if tc.isGetChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}

			if tc.isStaffChecked {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.staffNftsRes, nil).Once()
			}

			if tc.isGetNoti {
				leaderNotiRepo.On("GetNotiByNeed", mock.Anything, mock.Anything).Return(tc.noti, tc.getNotiErr)
			}

			if tc.isGetProposal {
				pendingWithdrawProposalRepo.On("IsPendingWithdrawProposalProposedWithSpecificInfo", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.isProposalProposed, tc.getIsProposedErr)
			}

			if tc.isCallCreateProposal {
				pendingWithdrawProposalRepo.On("CreatePendingWithdrawProposal", mock.Anything, mock.Anything).Return(tc.createProposalRes)
			}

			if tc.isGetPools {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: samplePoolObj,
							},
						},
					},
				}, nil).Once()

				mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return([]*models.SuiObjectResponse{
					&models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: samplePoolObj,
								},
							},
						},
					},
				}, nil).Once()
			}

			if tc.isCallAiValidate {
				walrusProvider.On("FetchBytesImage", mock.Anything).Return([]byte{}, nil)
				aiClient.On("ValidateWithdrawProposal", mock.Anything, mock.Anything).Return("")
			}

			res, err := service.CreateMealNeedWithdrawProposalV2(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.True(t, res != nil)
			}
		})
	}
}

func TestCreateSpecialNeedWithdrawProposalV2(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var aiClient = pkg.InitializeAiMockClient()
	var walrusProvider = pkg.InitializeWalrusMockProvider()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		aiClient,
		walrusProvider,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var childRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var validCampaignRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleSpecialNeedCampaignJson,
				},
			},
		},
	}

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var sampleBlobId string = "blobID"
	var tcsInfo = []struct {
		req                  request.CreateSpecialNeedWithdrawProposalRequest
		camapginRes          models.SuiObjectResponse
		isGetCampaign        bool
		isGetPoolsWithChild  bool
		isCallCreateProposal bool
		createProposalRes    error
		expectedErr          error
		isCallAiValidate     bool
		ctx                  context.Context
	}{
		{ // Happy case
			req: request.CreateSpecialNeedWithdrawProposalRequest{
				CampaignID:  sampleAddress,
				Amount:      100000,
				ProofBlobID: &sampleBlobId,
			},
			camapginRes:          validCampaignRes,
			isGetCampaign:        true,
			isGetPoolsWithChild:  true,
			isCallCreateProposal: true,
			isCallAiValidate:     true,
			ctx:                  ctx,
		},
		{ // Invalid campaign case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist campaign case
			req: request.CreateSpecialNeedWithdrawProposalRequest{
				CampaignID: sampleAddress,
			},
			camapginRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetCampaign: true,
			expectedErr:   errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:           ctx,
		},
		{ // Not creator case
			req: request.CreateSpecialNeedWithdrawProposalRequest{
				CampaignID: sampleAddress,
			},
			camapginRes:   validCampaignRes,
			isGetCampaign: true,
			expectedErr:   errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:           context.WithValue(context.Background(), "address", "1"),
		},
		{ // Out of budget case
			req: request.CreateSpecialNeedWithdrawProposalRequest{
				CampaignID: sampleAddress,
				Amount:     150000,
			},
			camapginRes:   validCampaignRes,
			isGetCampaign: true,
			expectedErr:   errors.New(noti.CURRENT_BUDGET_NOT_ENOUGH_MESSAGE),
			ctx:           ctx,
		},
		{ // Create pending proposal fail case
			req: request.CreateSpecialNeedWithdrawProposalRequest{
				CampaignID:  sampleAddress,
				Amount:      100000,
				ProofBlobID: &sampleBlobId,
			},
			camapginRes:          validCampaignRes,
			isGetCampaign:        true,
			isGetPoolsWithChild:  true,
			isCallCreateProposal: true,
			createProposalRes:    errors.New(noti.INTERNALL_ERR_MSG),
			isCallAiValidate:     true,
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
			ctx:                  ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil
			walrusProvider.ExpectedCalls = nil
			aiClient.ExpectedCalls = nil

			if tc.isGetCampaign {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.camapginRes, nil).Once()
			}

			if tc.isGetPoolsWithChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: samplePoolObj,
							},
						},
					},
				}, nil).Once()

				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(childRes, nil).Once()

				mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return([]*models.SuiObjectResponse{
					&models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: samplePoolObj,
								},
							},
						},
					},
				}, nil).Once()
			}

			if tc.isCallCreateProposal {
				pendingWithdrawProposalRepo.On("CreatePendingWithdrawProposal", mock.Anything, mock.Anything).Return(tc.createProposalRes)
			}

			if tc.isCallAiValidate {
				walrusProvider.On("FetchBytesImage", mock.Anything).Return([]byte{}, nil)
				aiClient.On("ValidateWithdrawProposal", mock.Anything, mock.Anything).Return("")
			}

			res, err := service.CreateSpecialNeedWithdrawProposalV2(tc.req, tc.ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.True(t, res != nil)
			}
		})
	}
}

func TestConfirmSpecialNeedProposal(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validProposalRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleSpecialNeedProposalJson,
				},
			},
		},
	}

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var tcsInfo = []struct {
		id                  string
		proposalRes         models.SuiObjectResponse
		isGetProposal       bool
		isGetSpecialNeedDao bool
		isBuildTx           bool
		buildTxRes          error
		expectedErr         error
		ctx                 context.Context
	}{
		{ // Happy case
			id:                  sampleAddress,
			proposalRes:         validProposalRes,
			isGetProposal:       true,
			isGetSpecialNeedDao: true,
			isBuildTx:           true,
			ctx:                 ctx,
		},
		{ // Invalid proposal case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist proposal case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetProposal: true,
			expectedErr:   errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:           ctx,
		},
		{ // Not creator case
			id:            sampleAddress,
			proposalRes:   validProposalRes,
			isGetProposal: true,
			expectedErr:   errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:           context.WithValue(context.Background(), "address", "1"),
		},
		{ // Proposal confirmed case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"creator":    sampleAddress,
								"is_confirm": true,
							},
						},
					},
				},
			},
			isGetProposal: true,
			expectedErr:   errors.New(noti.SPECIAL_NEED_PROPOSAL_CONFIRMED_MESSAGE),
			ctx:           ctx,
		},
		{ // Proposal pending case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"creator":    sampleAddress,
								"is_confirm": false,
								"closed_at":  "4102444800000", // Year: 2100
							},
						},
					},
				},
			},
			isGetProposal: true,
			expectedErr:   errors.New(noti.STILL_PENDING_REQUEST_MESSAGE),
			ctx:           ctx,
		},
		{ // Proposal fail condition case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"child":           sampleAddress,
								"target":          "1000000",
								"total_donated":   "900000",
								"withdraw_amount": "800000",
								"creator":         sampleAddress,
								"approvers":       []string{"", "", "", "", "", "", ""},
								"refusers":        []string{""},
								"approve_weight":  "500000",
								"refuse_weight":   "50000",
								"is_confirm":      false,
								"closed_at":       "1262304000000", // year 2010
							},
						},
					},
				},
			},
			isGetProposal:       true,
			isGetSpecialNeedDao: true,
			expectedErr:         errors.New(noti.PROPOSAL_FAIL_CONDITION_TO_CONFIRM_MESSAGE),
			ctx:                 ctx,
		},
		{ // Build tx fail case
			id:                  sampleAddress,
			proposalRes:         validProposalRes,
			isGetProposal:       true,
			isGetSpecialNeedDao: true,
			isBuildTx:           true,
			buildTxRes:          errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:         errors.New(noti.INTERNALL_ERR_MSG),
			ctx:                 ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			if tc.isGetProposal {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.proposalRes, nil).Once()
			}

			if tc.isGetSpecialNeedDao {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: sampleDaoObjJson,
							},
						},
					},
				}, nil).Once()
			}

			if tc.isBuildTx {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{}, tc.buildTxRes)
			}

			_, err := service.ConfirmSpecialNeedProposal(tc.id, tc.ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestVoteSpecialNeedProposal(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validProposalRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"child":           sampleAddress,
						"target":          "1000000",
						"total_donated":   "900000",
						"withdraw_amount": "800000",
						"creator":         sampleAddress,
						"approvers":       []string{"", "", "", "", "", "", "", "", "", ""},
						"refusers":        []string{""},
						"approve_weight":  "500000",
						"refuse_weight":   "50000",
						"is_confirm":      false,
						"closed_at":       "4102444800000", // year 2100
					},
				},
			},
		},
	}

	var nftsRes = models.PaginatedObjectsResponse{
		Data: []models.SuiObjectResponse{
			models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"id": map[string]string{
									"id": "",
								},
							},
						},
					},
				},
			},
		},
	}

	var ctx = context.WithValue(context.Background(), "address", "address")
	var tcsInfo = []struct {
		id            string
		proposalRes   models.SuiObjectResponse
		isGetProposal bool
		nftsRes       models.PaginatedObjectsResponse
		isGetNfts     bool
		isBuildTx     bool
		buildTxRes    error
		expectedErr   error
		ctx           context.Context
	}{
		{ // Happy case
			id:            sampleAddress,
			proposalRes:   validProposalRes,
			isGetProposal: true,
			nftsRes:       nftsRes,
			isGetNfts:     true,
			isBuildTx:     true,
			ctx:           ctx,
		},
		{ // Invalid proposal case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist proposal case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetProposal: true,
			expectedErr:   errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:           ctx,
		},
		{ // Creator case
			id:            sampleAddress,
			proposalRes:   validProposalRes,
			isGetProposal: true,
			expectedErr:   errors.New(noti.OWNER_VOTE_WARN_MSG),
			ctx:           context.WithValue(context.Background(), "address", sampleAddress),
		},
		{ // Proposal closed case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"creator":    sampleAddress,
								"is_confirm": false,
								"closed_at":  "1262304000000", // Year: 2010
							},
						},
					},
				},
			},
			isGetProposal: true,
			expectedErr:   errors.New(noti.REQUEST_CLOSED_MESSAGE),
			ctx:           ctx,
		},
		{ // Voted case
			id: sampleAddress,
			proposalRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"child":           sampleAddress,
								"target":          "1000000",
								"total_donated":   "900000",
								"withdraw_amount": "800000",
								"creator":         sampleAddress,
								"approvers":       []string{"", "", "", "", "", "", ""},
								"refusers":        []string{"address"},
								"approve_weight":  "500000",
								"refuse_weight":   "50000",
								"is_confirm":      false,
								"closed_at":       "4102444800000", // year 2100
							},
						},
					},
				},
			},
			isGetProposal: true,
			expectedErr:   errors.New(noti.ALREADY_VOTE_MESSAGE),
			ctx:           ctx,
		},
		{ // Not donor case
			id:            sampleAddress,
			proposalRes:   validProposalRes,
			isGetProposal: true,
			nftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.HAVE_TO_DONATE_TO_VOTE),
			ctx:         ctx,
		},
		{ // Build tx fail case
			id:            sampleAddress,
			proposalRes:   validProposalRes,
			isGetProposal: true,
			nftsRes:       nftsRes,
			isGetNfts:     true,
			isBuildTx:     true,
			buildTxRes:    errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:   errors.New(noti.INTERNALL_ERR_MSG),
			ctx:           ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			if tc.isGetProposal {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.proposalRes, nil)
			}

			if tc.isGetNfts {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.nftsRes, nil)
			}

			if tc.isBuildTx {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{}, tc.buildTxRes)
			}

			_, err := service.VoteSpecialNeedProposal(tc.id, request.VoteRequest{}, tc.ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUpdateChildBooksNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validChildRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var nftsRes = models.PaginatedObjectsResponse{
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
	}

	var needRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleNotSupportedBooksNeedJson,
				},
			},
		},
	}

	var notUpdatedNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"is_updated": false,
					},
				},
			},
		},
	}

	var value int64 = 10000
	var req = request.UpdateChildNeedRequest{
		ChildID: sampleAddress,
		NeedID:  sampleAddress,
		Value:   &value,
	}

	var zeroValue int64
	var zeroValueReq = request.UpdateChildNeedRequest{
		ChildID: sampleAddress,
		NeedID:  sampleAddress,
		Value:   &zeroValue,
	}

	var standerizeDate = func(src string) string {
		if len(src) == 2 {
			return src
		}

		return "0" + src
	}

	var curTime time.Time = time.Now()
	var previousDate time.Time = curTime.AddDate(0, 0, -1)
	var nextDate time.Time = curTime.AddDate(0, 0, 1)
	var validEditDatesRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"start_date": fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", previousDate.Day())), standerizeDate(fmt.Sprintf("%d", previousDate.Month()))),
						"end_date":   fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextDate.Day())), standerizeDate(fmt.Sprintf("%d", nextDate.Month()))),
					},
				},
			},
		},
	}

	var nextWeekDate time.Time = curTime.AddDate(0, 0, 7)
	var passedEditDatesRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"start_date": fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextDate.Day())), standerizeDate(fmt.Sprintf("%d", nextDate.Month()))),
						"end_date":   fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextWeekDate.Day())), standerizeDate(fmt.Sprintf("%d", nextWeekDate.Month()))),
					},
				},
			},
		},
	}

	var ctx = context.WithValue(context.Background(), "address", "address")
	var tcsInfo = []struct {
		req            request.UpdateChildNeedRequest
		childRes       models.SuiObjectResponse
		isGetChild     bool
		staffNftsRes   models.PaginatedObjectsResponse
		isGetNfts      bool
		needRes        models.SuiObjectResponse
		isGetNeed      bool
		editDatesRes   models.SuiObjectResponse
		isGetEditDates bool
		isBuildTx      bool
		buildTxRes     error
		expectedErr    error
		ctx            context.Context
	}{
		{ // Happy case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   validEditDatesRes,
			isGetEditDates: true,
			isBuildTx:      true,
			ctx:            ctx,
		},
		{ // Empty value case
			req: request.UpdateChildNeedRequest{
				ChildID: sampleAddress,
				NeedID:  sampleAddress,
			},
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			ctx:          ctx,
		},
		{ // Invalid need and child case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not exist child case
			req: req,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetChild:  true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not staff case
			req:        req,
			childRes:   validChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not leader case
			req:        req,
			childRes:   validChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonVolunteerNft10,
								},
							},
						},
					},
				},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:         ctx,
		},
		{ // Zero value case
			req:          zeroValueReq,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			expectedErr:  errors.New(noti.NEED_VALUE_INVALID_WARN_MSG),
			ctx:          ctx,
		},
		{ // Not updated need build tx fail case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes:      notUpdatedNeedRes,
			isGetNeed:    true,
			isBuildTx:    true,
			buildTxRes:   errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:  errors.New(noti.INTERNALL_ERR_MSG),
			ctx:          ctx,
		},
		{ // Not updated need build tx success case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes:      notUpdatedNeedRes,
			isGetNeed:    true,
			isBuildTx:    true,
			ctx:          ctx,
		},
		{ // Updated need year changes case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"supported_years": []string{"0"},
								"year_changes":    []string{fmt.Sprintf("%d", curTime.Year())},
								"year":            "1",
								"value":           "10000",
								"child":           sampleAddress,
								"is_updated":      true,
							},
						},
					},
				},
			},
			isGetNeed:      true,
			editDatesRes:   passedEditDatesRes,
			isGetEditDates: true,
			expectedErr:    errors.New(noti.CHILD_NEED_UPDATED_MESSAGE),
			ctx:            ctx,
		},
		{ // Updated need not edit date case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   passedEditDatesRes,
			isGetEditDates: true,
			expectedErr:    errors.New(noti.NOTE_UPDATE_CHILD_NEED_DATE_MESSAGE),
			ctx:            ctx,
		},
		{ // Updated need build tx fail case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   validEditDatesRes,
			isGetEditDates: true,
			isBuildTx:      true,
			buildTxRes:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:    errors.New(noti.INTERNALL_ERR_MSG),
			ctx:            ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			if tc.isGetChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}

			if tc.isGetNfts {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.staffNftsRes, nil)
			}

			if tc.isGetNeed {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.needRes, nil).Once()
			}

			if tc.isGetEditDates {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.editDatesRes, nil).Once()
			}

			if tc.isBuildTx {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{}, tc.buildTxRes)
			}

			_, err := service.UpdateBooksNeed(tc.req, tc.ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUpdateChildHealthInsuranceNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validChildRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var nftsRes = models.PaginatedObjectsResponse{
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
	}

	var needRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleNotSupportedBooksNeedJson,
				},
			},
		},
	}

	var notUpdatedNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"is_updated": false,
					},
				},
			},
		},
	}

	var value int64 = 10000
	var req = request.UpdateChildNeedRequest{
		ChildID: sampleAddress,
		NeedID:  sampleAddress,
		Value:   &value,
	}

	var zeroValue int64
	var zeroValueReq = request.UpdateChildNeedRequest{
		ChildID: sampleAddress,
		NeedID:  sampleAddress,
		Value:   &zeroValue,
	}

	var standerizeDate = func(src string) string {
		if len(src) == 2 {
			return src
		}

		return "0" + src
	}

	var curTime time.Time = time.Now()
	var previousDate time.Time = curTime.AddDate(0, 0, -1)
	var nextDate time.Time = curTime.AddDate(0, 0, 1)
	var validEditDatesRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"start_date": fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", previousDate.Day())), standerizeDate(fmt.Sprintf("%d", previousDate.Month()))),
						"end_date":   fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextDate.Day())), standerizeDate(fmt.Sprintf("%d", nextDate.Month()))),
					},
				},
			},
		},
	}

	var nextWeekDate time.Time = curTime.AddDate(0, 0, 7)
	var passedEditDatesRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"start_date": fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextDate.Day())), standerizeDate(fmt.Sprintf("%d", nextDate.Month()))),
						"end_date":   fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextWeekDate.Day())), standerizeDate(fmt.Sprintf("%d", nextWeekDate.Month()))),
					},
				},
			},
		},
	}

	var ctx = context.WithValue(context.Background(), "address", "address")
	var tcsInfo = []struct {
		req            request.UpdateChildNeedRequest
		childRes       models.SuiObjectResponse
		isGetChild     bool
		staffNftsRes   models.PaginatedObjectsResponse
		isGetNfts      bool
		needRes        models.SuiObjectResponse
		isGetNeed      bool
		editDatesRes   models.SuiObjectResponse
		isGetEditDates bool
		isBuildTx      bool
		buildTxRes     error
		expectedErr    error
		ctx            context.Context
	}{
		{ // Happy case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   validEditDatesRes,
			isGetEditDates: true,
			isBuildTx:      true,
			ctx:            ctx,
		},
		{ // Empty value case
			req: request.UpdateChildNeedRequest{
				ChildID: sampleAddress,
				NeedID:  sampleAddress,
			},
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			ctx:          ctx,
		},
		{ // Invalid need and child case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not exist child case
			req: req,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetChild:  true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not staff case
			req:        req,
			childRes:   validChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not leader case
			req:        req,
			childRes:   validChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonVolunteerNft10,
								},
							},
						},
					},
				},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:         ctx,
		},
		{ // Zero value case
			req:          zeroValueReq,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			expectedErr:  errors.New(noti.NEED_VALUE_INVALID_WARN_MSG),
			ctx:          ctx,
		},
		{ // Not updated need build tx fail case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes:      notUpdatedNeedRes,
			isGetNeed:    true,
			isBuildTx:    true,
			buildTxRes:   errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:  errors.New(noti.INTERNALL_ERR_MSG),
			ctx:          ctx,
		},
		{ // Not updated need build tx success case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes:      notUpdatedNeedRes,
			isGetNeed:    true,
			isBuildTx:    true,
			ctx:          ctx,
		},
		{ // Updated need year changes case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"supported_years": []string{"0"},
								"year_changes":    []string{fmt.Sprintf("%d", curTime.Year())},
								"year":            "1",
								"value":           "10000",
								"child":           sampleAddress,
								"is_updated":      true,
							},
						},
					},
				},
			},
			isGetNeed:      true,
			editDatesRes:   passedEditDatesRes,
			isGetEditDates: true,
			expectedErr:    errors.New(noti.CHILD_NEED_UPDATED_MESSAGE),
			ctx:            ctx,
		},
		{ // Updated need not edit date case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   passedEditDatesRes,
			isGetEditDates: true,
			expectedErr:    errors.New(noti.NOTE_UPDATE_CHILD_NEED_DATE_MESSAGE),
			ctx:            ctx,
		},
		{ // Updated need build tx fail case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   validEditDatesRes,
			isGetEditDates: true,
			isBuildTx:      true,
			buildTxRes:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:    errors.New(noti.INTERNALL_ERR_MSG),
			ctx:            ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			if tc.isGetChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}

			if tc.isGetNfts {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.staffNftsRes, nil)
			}

			if tc.isGetNeed {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.needRes, nil).Once()
			}

			if tc.isGetEditDates {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.editDatesRes, nil).Once()
			}

			if tc.isBuildTx {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{}, tc.buildTxRes)
			}

			_, err := service.UpdateHealthInsuranceNeed(tc.req, tc.ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestUpdateChildMealNeed(t *testing.T) {
	var pendingChildSpecialNeedRepo = repository.InitializeChildPendingSpecialProposalMockRepo()
	var pendingWithdrawProposalRepo = repository.InitializePendingWithdrawProposalMockRepo()
	var offchainWithdrawRepo = repository.InitializeOffChainWithdrawProposalMockRepo()
	var offchainDonationRepo = repository.InitializeOffChainDonationMockRepo()
	var paymentRepo = repository.InitializePaymentMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var bankProfileRepo = repository.InitializeBankProfileMockRepo()
	var leaderNotiRepo = repository.InitializeLeaderNotiMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		nil,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		nil,
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var validChildRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleJsonChild1,
				},
			},
		},
	}

	var nftsRes = models.PaginatedObjectsResponse{
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
	}

	var needRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: sampleNotSupportedBooksNeedJson,
				},
			},
		},
	}

	var notUpdatedNeedRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"is_updated": false,
					},
				},
			},
		},
	}

	var value int64 = 10000
	var req = request.UpdateChildNeedRequest{
		ChildID: sampleAddress,
		NeedID:  sampleAddress,
		Value:   &value,
	}

	var zeroValue int64
	var zeroValueReq = request.UpdateChildNeedRequest{
		ChildID: sampleAddress,
		NeedID:  sampleAddress,
		Value:   &zeroValue,
	}

	var standerizeDate = func(src string) string {
		if len(src) == 2 {
			return src
		}

		return "0" + src
	}

	var curTime time.Time = time.Now()
	var previousDate time.Time = curTime.AddDate(0, 0, -1)
	var nextDate time.Time = curTime.AddDate(0, 0, 1)
	var validEditDatesRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"start_date": fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", previousDate.Day())), standerizeDate(fmt.Sprintf("%d", previousDate.Month()))),
						"end_date":   fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextDate.Day())), standerizeDate(fmt.Sprintf("%d", nextDate.Month()))),
					},
				},
			},
		},
	}

	var nextWeekDate time.Time = curTime.AddDate(0, 0, 7)
	var passedEditDatesRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"start_date": fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextDate.Day())), standerizeDate(fmt.Sprintf("%d", nextDate.Month()))),
						"end_date":   fmt.Sprintf("%s/%s", standerizeDate(fmt.Sprintf("%d", nextWeekDate.Day())), standerizeDate(fmt.Sprintf("%d", nextWeekDate.Month()))),
					},
				},
			},
		},
	}

	var ctx = context.WithValue(context.Background(), "address", "address")
	var tcsInfo = []struct {
		req            request.UpdateChildNeedRequest
		childRes       models.SuiObjectResponse
		isGetChild     bool
		staffNftsRes   models.PaginatedObjectsResponse
		isGetNfts      bool
		needRes        models.SuiObjectResponse
		isGetNeed      bool
		editDatesRes   models.SuiObjectResponse
		isGetEditDates bool
		isBuildTx      bool
		buildTxRes     error
		expectedErr    error
		ctx            context.Context
	}{
		{ // Happy case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   validEditDatesRes,
			isGetEditDates: true,
			isBuildTx:      true,
			ctx:            ctx,
		},
		{ // Empty value case
			req: request.UpdateChildNeedRequest{
				ChildID: sampleAddress,
				NeedID:  sampleAddress,
			},
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			ctx:          ctx,
		},
		{ // Invalid need and child case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not exist child case
			req: req,
			childRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			isGetChild:  true,
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not staff case
			req:        req,
			childRes:   validChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:         ctx,
		},
		{ // Not leader case
			req:        req,
			childRes:   validChildRes,
			isGetChild: true,
			staffNftsRes: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonVolunteerNft10,
								},
							},
						},
					},
				},
			},
			isGetNfts:   true,
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			ctx:         ctx,
		},
		{ // Zero value case
			req:          zeroValueReq,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			expectedErr:  errors.New(noti.NEED_VALUE_INVALID_WARN_MSG),
			ctx:          ctx,
		},
		{ // Not updated need build tx fail case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes:      notUpdatedNeedRes,
			isGetNeed:    true,
			isBuildTx:    true,
			buildTxRes:   errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:  errors.New(noti.INTERNALL_ERR_MSG),
			ctx:          ctx,
		},
		{ // Not updated need build tx success case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes:      notUpdatedNeedRes,
			isGetNeed:    true,
			isBuildTx:    true,
			ctx:          ctx,
		},
		{ // Updated need year changes case
			req:          req,
			childRes:     validChildRes,
			isGetChild:   true,
			staffNftsRes: nftsRes,
			isGetNfts:    true,
			needRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"supported_years": []string{"0"},
								"year_changes":    []string{fmt.Sprintf("%d", curTime.Year())},
								"year":            fmt.Sprintf("%d", curTime.Year()),
								"value":           "10000",
								"child":           sampleAddress,
								"is_updated":      true,
							},
						},
					},
				},
			},
			isGetNeed:      true,
			editDatesRes:   passedEditDatesRes,
			isGetEditDates: true,
			expectedErr:    errors.New(noti.CHILD_NEED_UPDATED_MESSAGE),
			ctx:            ctx,
		},
		{ // Updated need not edit date case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   passedEditDatesRes,
			isGetEditDates: true,
			expectedErr:    errors.New(noti.NOTE_UPDATE_CHILD_NEED_DATE_MESSAGE),
			ctx:            ctx,
		},
		{ // Updated need build tx fail case
			req:            req,
			childRes:       validChildRes,
			isGetChild:     true,
			staffNftsRes:   nftsRes,
			isGetNfts:      true,
			needRes:        needRes,
			isGetNeed:      true,
			editDatesRes:   validEditDatesRes,
			isGetEditDates: true,
			isBuildTx:      true,
			buildTxRes:     errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:    errors.New(noti.INTERNALL_ERR_MSG),
			ctx:            ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			if tc.isGetChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}

			if tc.isGetNfts {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.staffNftsRes, nil)
			}

			if tc.isGetNeed {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.needRes, nil).Once()
			}

			if tc.isGetEditDates {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.editDatesRes, nil).Once()
			}

			if tc.isBuildTx {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{}, tc.buildTxRes)
			}

			_, err := service.UpdateMealNeed(tc.req, tc.ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
