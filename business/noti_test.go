package business

import (
	"context"
	"errors"
	"fmt"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
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

type mockLeaderNotiRepo struct {
	mock.Mock
}

func (m *mockLeaderNotiRepo) GetNoti(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	args := m.Called(id, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LeaderNoti), args.Error(1)
}
func (m *mockLeaderNotiRepo) GetNotiByMealNeed(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	args := m.Called(id, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LeaderNoti), args.Error(1)
}
func (m *mockLeaderNotiRepo) GetNotiByNeed(id string, ctx context.Context) (*entities.LeaderNoti, error) {
	args := m.Called(id, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LeaderNoti), args.Error(1)
}
func (m *mockLeaderNotiRepo) GetCurrentLeaderNotis(req request.GetNotisRequest, leader string, ctx context.Context) ([]entities.LeaderNoti, error) {
	args := m.Called(req, leader, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.LeaderNoti), args.Error(1)
}
func (m *mockLeaderNotiRepo) CreateNoti(noti entities.LeaderNoti, ctx context.Context) error {
	args := m.Called(noti, ctx)
	return args.Error(0)
}
func (m *mockLeaderNotiRepo) UpdateNoti(noti entities.LeaderNoti, ctx context.Context) error {
	args := m.Called(noti, ctx)
	return args.Error(0)
}
func (m *mockLeaderNotiRepo) AssignLeader(leader, region string, ctx context.Context) error {
	args := m.Called(leader, region, ctx)
	return args.Error(0)
}

func TestGetCurrentWalletNotis(t *testing.T) {
	mockRepo := new(mockLeaderNotiRepo)
	var mockClient = pkg.InitializeSuiMockApi()

	// Initialize noti service cleanly so cache logic is respected internally
	// we just use the initialized function
	service := initializeNotiService(mockRepo, map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient}, util.GetLogConfig(shared.ERROR_LEVEL))

	ctx := context.Background()

	tcs := []struct {
		wallet         string
		req            request.GetNotisRequest
		mockOwnedNfts  []models.SuiObjectResponse
		mockOwnedErr   error
		mockDBNotis    []entities.LeaderNoti
		mockDBErr      error
		expectedAmount int
		expectedPage   int
		expectedErr    error
	}{
		{ // UTCID01: Normal
			wallet: sampleAddress, req: request.GetNotisRequest{Page: 2, PageSize: 5},
			mockOwnedNfts:  []models.SuiObjectResponse{{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"role": local_leader_role}}}}}},
			mockDBNotis:    []entities.LeaderNoti{{ID: "1"}, {ID: "2"}}, expectedAmount: 2, expectedPage: 2,
		},
		{ // UTCID02: Invalid Address Format
			wallet: "invalid", expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID03: Fetch On-chain Objects Fail
			wallet: sampleAddress, mockOwnedErr: errors.New("onchain error"), expectedErr: errors.New(noti.INTERNALL_ERR_MSG),
		},
		{ // UTCID04: Empty NFT Results
			wallet: sampleAddress, mockOwnedNfts: []models.SuiObjectResponse{}, expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID05: Profile Lacking Proper Role
			wallet: sampleAddress, mockOwnedNfts: []models.SuiObjectResponse{{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"role": "normal_staff"}}}}}}, expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID06: Pagination Normal Boundary
			wallet: sampleAddress, req: request.GetNotisRequest{Page: 0, PageSize: 0},
			mockOwnedNfts:  []models.SuiObjectResponse{{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"role": local_leader_role}}}}}},
			mockDBNotis:    []entities.LeaderNoti{{ID: "1"}}, expectedAmount: 1, expectedPage: 1,
		},
		{ // UTCID07: Nil DB Returns Response
			wallet: sampleAddress, req: request.GetNotisRequest{Page: 1},
			mockOwnedNfts:  []models.SuiObjectResponse{{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"role": local_leader_role}}}}}},
			mockDBNotis:    nil, expectedAmount: 0, expectedPage: 1,
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			mockRepo.ExpectedCalls = nil

			if util.IsValidSuiAddressStrict(tc.wallet) {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(models.PaginatedObjectsResponse{Data: tc.mockOwnedNfts}, tc.mockOwnedErr)

				if tc.mockOwnedErr == nil && len(tc.mockOwnedNfts) > 0 {
					var hasRole bool
					for _, item := range tc.mockOwnedNfts {
						if item.Data != nil && item.Data.Content != nil && item.Data.Content.SuiMoveObject.Fields["role"] == local_leader_role {
							hasRole = true
							break
						}
					}

					if hasRole {
						// Note: By default, test redis cache mock won't have the entry since no REDIS URL locally, or its Get() safely returns false!
						expectedReq := tc.req
						if expectedReq.Page < 1 {
							expectedReq.Page = 1
						}
						// `default_page_size` must be used here
						if expectedReq.PageSize < 1 {
							expectedReq.PageSize = default_page_size
						}
						
						mockRepo.On("GetCurrentLeaderNotis", expectedReq, tc.wallet, mock.Anything).Return(tc.mockDBNotis, tc.mockDBErr)
					}
				}
			}

			res, err := service.GetCurrentWalletNotis(tc.wallet, tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedErr == nil {
				assert.Equal(t, tc.expectedAmount, res.Amount)
				assert.Equal(t, tc.expectedPage, res.Page)
			}
		})
	}
}
