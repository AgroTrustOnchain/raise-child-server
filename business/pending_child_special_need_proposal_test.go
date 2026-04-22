package business

import (
	"context"
	"errors"
	"fmt"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"
	"raise-child/util"
	"strings"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPendingChildSpecialNeedProposalRepo struct {
	mock.Mock
}

func (m *mockPendingChildSpecialNeedProposalRepo) GetPendingChildSpecialNeedProposal(id string, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error) {
	args := m.Called(id, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PendingChildSpecialNeedProposal), args.Error(1)
}

func (m *mockPendingChildSpecialNeedProposalRepo) GetPendingChildSpecialNeedProposals(req request.GetPendingChildSpecialNeedProposalsRequest, ctx context.Context) ([]entities.PendingChildSpecialNeedProposal, int, error) {
	args := m.Called(req, ctx)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]entities.PendingChildSpecialNeedProposal), args.Int(1), args.Error(2)
}

func (m *mockPendingChildSpecialNeedProposalRepo) CreatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error {
	args := m.Called(proposal, ctx)
	return args.Error(0)
}

func (m *mockPendingChildSpecialNeedProposalRepo) UpdatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error {
	args := m.Called(proposal, ctx)
	return args.Error(0)
}

func ptrInt64(i int64) *int64 {
	return &i
}

var otherAddress = "0x" + strings.Repeat("1", 64)

func TestApprovePendingChildSpecialNeedProposal(t *testing.T) {
	t.Setenv(env.POOL_ID, "mainPool1")
	var mockClient = pkg.InitializeSuiMockApi()
	mockRepo := new(mockPendingChildSpecialNeedProposalRepo)
	service := initializePendingChildSpecialNeedProposalService(
		mockRepo,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	tcs := []struct {
		id                  string
		mockProposal        *entities.PendingChildSpecialNeedProposal
		mockProposalErr     error
		mockManage          models.SuiObjectResponse
		mockManageErr       error
		mockMainPool        models.SuiObjectResponse
		mockMainPoolErr     error
		mockLocalPools      []*models.SuiObjectResponse
		mockLocalPoolsErr   error
		mockTxBytes         string
		mockTxErr           error
		mockUpdateErr       error
		expectedErr         error
	}{
		{ // UTCID01: Normal Approve (Cache miss, triggers queries)
			id:              "p1",
			mockProposal:    &entities.PendingChildSpecialNeedProposal{Region: "reg1", ReviewStatus: request_pending_status},
			mockManage:      models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
			mockMainPool:    models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"local_pools": []interface{}{"lp1"}}}}}},
			mockLocalPools:  []*models.SuiObjectResponse{{Data: &models.SuiObjectData{ObjectId: "lp1", Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"region": "reg1"}}}}}},
			mockTxBytes:     "bytes",
		},
		{ // UTCID02: Fetch Repo Error
			id: "p1", mockProposalErr: errors.New("db error"), expectedErr: errors.New("db error"),
		},
		{ // UTCID03: Proposal Not Found
			id: "p1", mockProposal: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID04: Already Reviewed
			id: "p1", mockProposal: &entities.PendingChildSpecialNeedProposal{ReviewStatus: request_approved_status}, expectedErr: errors.New(noti.REQUEST_REVIEWED_MESSAGE),
		},
		{ // UTCID05: Admin Unauthorized
			id: "p1", mockProposal: &entities.PendingChildSpecialNeedProposal{ReviewStatus: request_pending_status},
			mockManage:  models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{otherAddress}}}}}},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID06: Main Pool Fetch Error
			id:              "p1",
			mockProposal:    &entities.PendingChildSpecialNeedProposal{Region: "reg1", ReviewStatus: request_pending_status},
			mockManage:      models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
			mockMainPoolErr: errors.New("onchain error"),
			expectedErr:     errors.New(noti.INTERNAL_ERR_MSG),
		},
		{ // UTCID07: Transaction Build Error
			id:              "p1",
			mockProposal:    &entities.PendingChildSpecialNeedProposal{Region: "reg1", ReviewStatus: request_pending_status},
			mockManage:      models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
			mockMainPool:    models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"local_pools": []interface{}{"lp1"}}}}}},
			mockLocalPools:  []*models.SuiObjectResponse{{Data: &models.SuiObjectData{ObjectId: "lp1", Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"region": "reg1"}}}}}},
			mockTxErr:       errors.New("tx error"), expectedErr: errors.New(noti.INTERNAL_ERR_MSG),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockRepo.ExpectedCalls = nil

			mockRepo.On("GetPendingChildSpecialNeedProposal", tc.id, mock.Anything).Return(tc.mockProposal, tc.mockProposalErr)

			if tc.mockProposal != nil && tc.mockProposal.ReviewStatus == request_pending_status {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.mockManage, tc.mockManageErr).Once()
				var isAdmin bool
				if tc.mockManageErr == nil {
					for _, ad := range tc.mockManage.Data.Content.SuiMoveObject.Fields["admin_ids"].([]interface{}) {
						if ad.(string) == sampleAddress {
							isAdmin = true
						}
					}
				}
				if isAdmin {
					mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.mockMainPool, tc.mockMainPoolErr)
					if tc.mockMainPoolErr == nil {
						mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return(tc.mockLocalPools, tc.mockLocalPoolsErr)
						if tc.mockLocalPoolsErr == nil {
							mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{TxBytes: tc.mockTxBytes}, tc.mockTxErr)
							if tc.mockTxErr == nil {
								mockRepo.On("UpdatePendingChildSpecialNeedProposal", mock.Anything, mock.Anything).Return(tc.mockUpdateErr)
							}
						}
					}
				}
			}

			_, err := service.ApprovePendingChildSpecialNeedProposal(tc.id, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetPendingChildSpecialNeedProposal(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockRepo := new(mockPendingChildSpecialNeedProposalRepo)
	service := initializePendingChildSpecialNeedProposalService(
		mockRepo,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	tcs := []struct {
		id              string
		mockProposal    *entities.PendingChildSpecialNeedProposal
		mockProposalErr error
		mockManage      models.SuiObjectResponse
		mockManageErr   error
		expectedErr     error
	}{
		{ // UTCID01: Actor access
			id:           "p1",
			mockProposal: &entities.PendingChildSpecialNeedProposal{ID: "p1", ActorAddress: sampleAddress},
		},
		{ // UTCID02: Admin access
			id:           "p1",
			mockProposal: &entities.PendingChildSpecialNeedProposal{ID: "p1", ActorAddress: otherAddress},
			mockManage:   models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
		},
		{ // UTCID03: Unauthorized external user
			id:           "p1",
			mockProposal: &entities.PendingChildSpecialNeedProposal{ID: "p1", ActorAddress: otherAddress},
			mockManage:   models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{otherAddress}}}}}},
			expectedErr:  errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID04: Fetch error
			id: "p1", mockProposalErr: errors.New("db error"), expectedErr: errors.New("db error"),
		},
		{ // UTCID05: Not found
			id: "p1", mockProposal: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockRepo.ExpectedCalls = nil

			mockRepo.On("GetPendingChildSpecialNeedProposal", tc.id, mock.Anything).Return(tc.mockProposal, tc.mockProposalErr)

			if tc.mockProposal != nil && tc.mockProposal.ActorAddress != sampleAddress {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.mockManage, tc.mockManageErr)
			}

			res, err := service.GetPendingChildSpecialNeedProposal(tc.id, ctx)
			if tc.expectedErr == nil {
				assert.Equal(t, tc.mockProposal, res)
			}
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetPendingChildSpecialNeedProposals(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockRepo := new(mockPendingChildSpecialNeedProposalRepo)
	service := initializePendingChildSpecialNeedProposalService(
		mockRepo,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	tcs := []struct {
		req             request.GetPendingChildSpecialNeedProposalsRequest
		mockManage      models.SuiObjectResponse
		mockManageErr   error
		mockDB          []entities.PendingChildSpecialNeedProposal
		mockPages       int
		mockDBErr       error
		expectedAmount  int
		expectedReturns bool
		expectedErr     error
	}{
		{ // UTCID01: Actor limits search to creator
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Creator: sampleAddress, MaxAmount: ptrInt64(10), MinAmount: ptrInt64(5)},
			mockDB:          []entities.PendingChildSpecialNeedProposal{{ID: "1"}}, expectedAmount: 1, expectedReturns: true,
		},
		{ // UTCID02: Admin fetching others
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Creator: otherAddress, Page: 0, PageSize: 0},
			mockManage:      models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
			mockDB:          []entities.PendingChildSpecialNeedProposal{{ID: "1"}}, expectedAmount: 1, expectedReturns: true,
		},
		{ // UTCID03: MaxAmount negative
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Creator: sampleAddress, MaxAmount: ptrInt64(-1)},
			expectedReturns: true, expectedAmount: 0,
		},
		{ // UTCID04: Min > Max
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Creator: sampleAddress, MinAmount: ptrInt64(10), MaxAmount: ptrInt64(5)},
			expectedReturns: true, expectedAmount: 0,
		},
		{ // UTCID05: Valid pagination limits default
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Creator: sampleAddress},
			mockDB:          []entities.PendingChildSpecialNeedProposal{{ID: "1"}}, expectedAmount: 1, expectedReturns: true,
		},
		{ // UTCID06: Invalid Reviewer String
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Reviewer: "badAddress"}, expectedReturns: false, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID07: Unauthorized View Request
			req:             request.GetPendingChildSpecialNeedProposalsRequest{Creator: otherAddress},
			mockManage:      models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{otherAddress}}}}}},
			expectedReturns: false, expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockRepo.ExpectedCalls = nil

			isValid := true
			if tc.req.Creator != "" && !util.IsValidSuiAddressStrict(tc.req.Creator) {
				isValid = false
			}
			if tc.req.Reviewer != "" && !util.IsValidSuiAddressStrict(tc.req.Reviewer) {
				isValid = false
			}

			if isValid && tc.req.Creator != sampleAddress {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.mockManage, tc.mockManageErr)
			}

			if tc.expectedReturns && tc.expectedAmount > 0 {
				mockRepo.On("GetPendingChildSpecialNeedProposals", mock.Anything, mock.Anything).Return(tc.mockDB, tc.mockPages, tc.mockDBErr)
			}

			res, err := service.GetPendingChildSpecialNeedProposals(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedReturns {
				assert.Equal(t, tc.expectedAmount, res.Amount)
			}
		})
	}
}

func TestRefusePendingChildSpecialNeedProposal(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockRepo := new(mockPendingChildSpecialNeedProposalRepo)
	service := initializePendingChildSpecialNeedProposalService(
		mockRepo,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	tcs := []struct {
		id              string
		mockProposal    *entities.PendingChildSpecialNeedProposal
		mockProposalErr error
		mockManage      models.SuiObjectResponse
		mockManageErr   error
		mockUpdateErr   error
		expectedErr     error
	}{
		{ // UTCID01: Normal Refuse
			id:           "p1",
			mockProposal: &entities.PendingChildSpecialNeedProposal{ReviewStatus: request_pending_status},
			mockManage:   models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
		},
		{ // UTCID02: Fetch Repo Error
			id: "p1", mockProposalErr: errors.New("db error"), expectedErr: errors.New("db error"),
		},
		{ // UTCID03: Proposal Not Found
			id: "p1", mockProposal: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID04: Already Reviewed
			id: "p1", mockProposal: &entities.PendingChildSpecialNeedProposal{ReviewStatus: request_approved_status}, expectedErr: errors.New(noti.REQUEST_REVIEWED_MESSAGE),
		},
		{ // UTCID05: Admin Unauthorized
			id:          "p1", mockProposal: &entities.PendingChildSpecialNeedProposal{ReviewStatus: request_pending_status},
			mockManage:  models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{otherAddress}}}}}},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockRepo.ExpectedCalls = nil

			mockRepo.On("GetPendingChildSpecialNeedProposal", tc.id, mock.Anything).Return(tc.mockProposal, tc.mockProposalErr)

			if tc.mockProposal != nil && tc.mockProposal.ReviewStatus == request_pending_status {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.mockManage, tc.mockManageErr).Once()
				var isAdmin bool
				if tc.mockManageErr == nil {
					for _, ad := range tc.mockManage.Data.Content.SuiMoveObject.Fields["admin_ids"].([]interface{}) {
						if ad.(string) == sampleAddress {
							isAdmin = true
						}
					}
				}
				if isAdmin {
					mockRepo.On("UpdatePendingChildSpecialNeedProposal", mock.Anything, mock.Anything).Return(tc.mockUpdateErr)
				}
			}

			err := service.RefusePendingChildSpecialNeedProposal(tc.id, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
