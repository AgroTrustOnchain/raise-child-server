package business

import (
	"context"
	"errors"
	"fmt"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
	mock_repo "raise-child/mocks/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
	"raise-child/util"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfirmWithdrawProposal(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	service := initializePaymentService(
		nil, nil, nil, nil, nil, nil, nil, nil, nil,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	tcs := []struct {
		id          string
		expectedErr error
	}{
		{id: sampleAddress, expectedErr: nil},
		{id: "invalid", expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			_, err := service.(*paymentService).ConfirmWithdrawProposal(tc.id, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestDonate(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockProfileRepo := mock_repo.InitializeProfileMockRepo()
	service := initializePaymentService(
		nil, nil, nil, nil, nil, mockProfileRepo, nil, nil, nil,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	ctx := context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", "pro1")

	identityStr := "some-identity"

	tcs := []struct {
		req            request.DonateRequest
		mockProfile    *entities.Profile
		mockProfileErr error
		expectedErr    error
	}{
		{ // UTCID01: Normal PayOS Init (API triggers error naturally)
			req:         request.DonateRequest{PoolId: sampleAddress, Amount: 1000},
			mockProfile: &entities.Profile{ID: "pro1", IdentityCode: &identityStr},
			expectedErr: errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // UTCID02: Invalid Pool Address
			req:         request.DonateRequest{PoolId: "invalid", Amount: 1000},
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID03: Profile Base Err
			req:            request.DonateRequest{PoolId: sampleAddress, Amount: 1000},
			mockProfileErr: errors.New("db error"), expectedErr: errors.New("db error"),
		},
		{ // UTCID04: Empty Profile Found
			req:         request.DonateRequest{PoolId: sampleAddress, Amount: 1000},
			mockProfile: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID05: Identity Map Exluded Profiles
			req:         request.DonateRequest{PoolId: sampleAddress, Amount: 1000},
			mockProfile: &entities.Profile{ID: "pro1", IdentityCode: nil},
			expectedErr: errors.New(noti.NOT_UPLOADED_PROFILE_MESSAGE),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockProfileRepo.ExpectedCalls = nil

			if util.IsValidSuiAddressStrict(tc.req.PoolId) {
				mockProfileRepo.On("GetProfile", "pro1", mock.Anything).Return(tc.mockProfile, tc.mockProfileErr)
			}

			_, err := service.Donate(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestCallback(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockPayment := mock_repo.InitializePaymentMockRepo()

	service := initializePaymentService(
		nil, nil, nil, nil, mockPayment, nil, nil, nil, nil,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)
	ctx := context.Background()

	tcs := []struct {
		id                 string
		mockPayment        *entities.Payment
		mockPaymentErr     error
		expectedReturnsUrl bool
		expectedErr        error
	}{
		{ // UTCID01: Redirect Resolving API Err
			id:          "p1",
			mockPayment: &entities.Payment{ID: "p1", TransactionId: "123", ExpiredAt: time.Now().Add(10 * time.Minute)},
			expectedErr: errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // UTCID02: Payment Fault
			id: "p1", mockPaymentErr: errors.New("db error"), expectedErr: errors.New("db error"),
		},
		{ // UTCID03: Nil Target
			id: "p1", mockPayment: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID04: Valid Expired Time Passing
			id: "p1", mockPayment: &entities.Payment{ID: "p1", TransactionId: "123", ExpiredAt: time.Now().Add(-1 * time.Minute)},
			expectedReturnsUrl: true,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockPayment.ExpectedCalls = nil
			mockPayment.On("GetPaymentById", tc.id, mock.Anything).Return(tc.mockPayment, tc.mockPaymentErr)

			url, err := service.Callback(tc.id, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedReturnsUrl {
				assert.NotEmpty(t, url)
			}
		})
	}
}

func TestCallbackWithAuth(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockPayment := mock_repo.InitializePaymentMockRepo()

	service := initializePaymentService(
		nil, nil, nil, nil, mockPayment, nil, nil, nil, nil,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)
	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	tcs := []struct {
		id                 string
		mockPayment        *entities.Payment
		mockPaymentErr     error
		expectedReturnsUrl bool
		expectedErr        error
	}{
		{ // UTCID01: Actor Resolving Normal
			id:          "p1",
			mockPayment: &entities.Payment{ID: "p1", Actor: sampleAddress, TransactionId: "123", ExpiredAt: time.Now().Add(10 * time.Minute)},
			expectedErr: errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // UTCID02: Unauthenticated Cross Request
			id:          "p1",
			mockPayment: &entities.Payment{ID: "p1", Actor: "0x123"},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID03: Target Payment DB Errors
			id: "p1", mockPaymentErr: errors.New("db query error"), expectedErr: errors.New("db query error"),
		},
		{ // UTCID04: Target Expiration Caught By Redirection
			id:                 "p1",
			mockPayment:        &entities.Payment{ID: "p1", Actor: sampleAddress, TransactionId: "123", ExpiredAt: time.Now().Add(-1 * time.Minute)},
			expectedReturnsUrl: true,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockPayment.ExpectedCalls = nil
			mockPayment.On("GetPaymentById", tc.id, mock.Anything).Return(tc.mockPayment, tc.mockPaymentErr)

			url, err := service.CallbackWithAuth(tc.id, "imgBlobId", ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedReturnsUrl {
				assert.NotEmpty(t, url)
			}
		})
	}
}
