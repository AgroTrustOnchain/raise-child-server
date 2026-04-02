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
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
	var tcsInfo = []struct {
		needId               string
		req                  request.SupportMealNeadRequest
		profile              *entities.Profile
		getProfileErr        error
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
			isCallCreateDonation: true,
			createDonateRes:      errors.New(noti.INTERNALL_ERR_MSG),
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
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
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
		expectedErr          error
	}{
		{ // Happy case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
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
				NeedID: sampleAddress,
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
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			leaderNotiRepo.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil

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
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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
		expectedErr          error
	}{
		{ // Happy case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
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
				NeedID: sampleAddress,
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
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			leaderNotiRepo.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil

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
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeChildService(
		pendingChildSpecialNeedRepo,
		pendingWithdrawProposalRepo,
		offchainWithdrawRepo,
		offchainDonationRepo,
		paymentRepo,
		profileRepo,
		bankProfileRepo,
		leaderNotiRepo,
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

	// var passedWithdrawDateNoti = entities.LeaderNoti{
	// 	ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime.AddDate(0, -1, 0))},
	// 	Contents:                []string{"1", "2"},
	// }

	// var stillNotComeWithdrawDateNoti = entities.LeaderNoti{
	// 	ExpectedWithdrawPeriods: []string{util.TimeToRawDate(curTime.AddDate(0, 1, 0))},
	// 	Contents:                []string{"1", "2"},
	// }

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)
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
		expectedErr          error
	}{
		{ // Happy case
			req: request.CreateNormalNeedWithdrawProposalRequest{
				NeedID: sampleAddress,
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
				NeedID: sampleAddress,
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
			expectedErr:          errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			leaderNotiRepo.ExpectedCalls = nil
			pendingWithdrawProposalRepo.ExpectedCalls = nil

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

			res, err := service.CreateMealNeedWithdrawProposalV2(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.True(t, res != nil)
			}
		})
	}
}
