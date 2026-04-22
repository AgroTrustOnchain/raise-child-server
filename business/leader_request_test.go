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
	"slices"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string {
	return &s
}

type mockLeaderRequestRepo struct {
	mock.Mock
}

func (m *mockLeaderRequestRepo) GetRequest(id string, ctx context.Context) (*entities.LocalLeaderRegistrationRequest, error) {
	args := m.Called(id, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.LocalLeaderRegistrationRequest), args.Error(1)
}

func (m *mockLeaderRequestRepo) GetRegistrationRequests(req request.GetNormalStaffRegistrationRequests, ctx context.Context) ([]entities.LocalLeaderRegistrationRequest, int, error) {
	args := m.Called(req, ctx)
	var res1 []entities.LocalLeaderRegistrationRequest
	if args.Get(0) != nil {
		res1 = args.Get(0).([]entities.LocalLeaderRegistrationRequest)
	}
	return res1, args.Int(1), args.Error(2)
}

func (m *mockLeaderRequestRepo) GetWalletRegistrationRequests(id string, ctx context.Context) ([]entities.LocalLeaderRegistrationRequest, error) {
	args := m.Called(id, ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.LocalLeaderRegistrationRequest), args.Error(1)
}

func (m *mockLeaderRequestRepo) CreateRegistrationRequest(req entities.LocalLeaderRegistrationRequest, ctx context.Context) error {
	args := m.Called(req, ctx)
	return args.Error(0)
}

func (m *mockLeaderRequestRepo) UpdateRegistrationRequest(req entities.LocalLeaderRegistrationRequest, ctx context.Context) error {
	args := m.Called(req, ctx)
	return args.Error(0)
}

func TestConfirmRequest(t *testing.T) {
	mockRepo := new(mockLeaderRequestRepo)
	var mockClient = pkg.InitializeSuiMockApi()
	service := &leaderRequestService{
		leaderRequestRepo: mockRepo,
		clients:           map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		errLogger:         util.GetLogConfig(shared.ERROR_LEVEL),
	}


	futureTime := time.Now().Add(time.Hour)
	pastTime := time.Now().Add(-time.Hour)

	tcs := []struct {
		id            string
		sender        string
		mockReq       *entities.LocalLeaderRegistrationRequest
		mockReqErr    error
		mockUpdateErr error
		mockOwnedCaps []entities.Cap
		mockOwnedErr  error
		mockTxBytes   string
		mockTxErr     error
		expectedErr   error
	}{
		{ // UTCID01: Normal Confirm Approved with Cap
			id: "req1", sender: sampleAddress,
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: sampleAddress, ClosedAt: pastTime, Approvers: []string{"a", "b"}, Refusers: []string{}, IsAvailableToConfirm: true}}},
			mockOwnedCaps: []entities.Cap{{ID: entities.ID{ID: "cap1"}}},
		},
		{ // UTCID02: Invalid format sender
			id: "req1", sender: "invalid", expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID03: Denied result
			id: "req1", sender: sampleAddress,
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: sampleAddress, ClosedAt: pastTime, Approvers: []string{}, Refusers: []string{"a", "b"}}}},
		},
		{ // UTCID04: Pending process
			id: "req1", sender: sampleAddress,
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: sampleAddress, ClosedAt: futureTime}}},
			expectedErr: errors.New(noti.STILL_PENDING_REQUEST_MESSAGE),
		},
		{ // UTCID05: Request belongs to another
			id: "req1", sender: sampleAddress,
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other"}}},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID06: Request not found
			id: "req1", sender: sampleAddress, mockReq: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil
			testCtx := context.WithValue(context.Background(), "address", tc.sender)

			if util.IsValidSuiAddressStrict(tc.sender) {
				mockRepo.On("GetRequest", tc.id, mock.Anything).Return(tc.mockReq, tc.mockReqErr)

				if tc.mockReq != nil && tc.mockReq.CreatedBy == tc.sender && tc.mockReq.ClosedAt.Before(time.Now()) {
					mockRepo.On("UpdateRegistrationRequest", mock.Anything, mock.Anything).Return(tc.mockUpdateErr)
					if tc.mockUpdateErr == nil && tc.mockReq.IsAvailableToConfirm {
						mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(models.PaginatedObjectsResponse{
							Data: []models.SuiObjectResponse{
								{
									Data: &models.SuiObjectData{
										Content: &models.SuiParsedData{
											SuiMoveObject: models.SuiMoveObject{
												Fields: map[string]interface{}{
													"id": map[string]interface{}{"id": "cap1"},
												},
											},
										},
									},
								},
							},
						}, tc.mockOwnedErr)

						if tc.mockOwnedErr == nil {
							mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{TxBytes: tc.mockTxBytes}, tc.mockTxErr)
						}
					}
				}
			}

			_, err := service.ConfirmRequest(tc.id, testCtx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestCreateRequest(t *testing.T) {
	mockRepo := new(mockLeaderRequestRepo)
	profileRepo := repository.InitializeProfileMockRepo()
	service := &leaderRequestService{
		leaderRequestRepo: mockRepo,
		profileRepo:       profileRepo,
	}

	ctx := context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)

	validProfile := &entities.Profile{IdentityCode: ptrString("123"), FirstName: ptrString("A"), LastName: ptrString("B")}

	tcs := []struct {
		req            request.CreateRegistrationRequest
		mockWalletReq  []entities.LocalLeaderRegistrationRequest
		mockWalletErr  error
		mockProfile    *entities.Profile
		mockProfileErr error
		mockCreateErr  error
		expectedReqNil bool
		expectedErr    error
	}{
		{ // UTCID01: Valid Create Request
			mockWalletReq:  []entities.LocalLeaderRegistrationRequest{},
			mockProfile:    validProfile, expectedReqNil: false, expectedErr: nil,
		},
		{ // UTCID02: Existing Pending Request
			mockWalletReq:  []entities.LocalLeaderRegistrationRequest{{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{Status: request_pending_status}}}},
			expectedReqNil: true, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID03: Existing Approved Request
			mockWalletReq:  []entities.LocalLeaderRegistrationRequest{{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{Status: request_approved_status}}}},
			expectedReqNil: true, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID04: Missing Profile
			mockWalletReq:  []entities.LocalLeaderRegistrationRequest{},
			mockProfile:    nil, expectedReqNil: true, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID05: Profile Fetch Error
			mockWalletReq:  []entities.LocalLeaderRegistrationRequest{},
			mockProfileErr: errors.New("db err"), expectedReqNil: true, expectedErr: errors.New("db err"),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			profileRepo.ExpectedCalls = nil

			mockRepo.On("GetWalletRegistrationRequests", sampleAddress, mock.Anything).Return(tc.mockWalletReq, tc.mockWalletErr)

			if tc.mockWalletErr == nil && (len(tc.mockWalletReq) == 0 || (tc.mockWalletReq[0].Status != request_pending_status && tc.mockWalletReq[0].Status != request_approved_status)) {
				profileRepo.On("GetProfile", sampleSub, mock.Anything).Return(tc.mockProfile, tc.mockProfileErr)
			}

			if tc.mockProfile != nil && tc.mockProfileErr == nil && len(tc.mockWalletReq) == 0 {
				mockRepo.On("CreateRegistrationRequest", mock.Anything, mock.Anything).Return(tc.mockCreateErr)
			}

			res, err := service.CreateRequest(tc.req, ctx)
			if tc.expectedReqNil {
				assert.Nil(t, res)
			} else {
				assert.NotNil(t, res)
			}
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetRequest(t *testing.T) {
	mockRepo := new(mockLeaderRequestRepo)
	service := &leaderRequestService{leaderRequestRepo: mockRepo}
	ctx := context.Background()

	tcs := []struct {
		id          string
		mockReq     *entities.LocalLeaderRegistrationRequest
		mockReqErr  error
		expectedErr error
	}{
		{id: "id", mockReq: &entities.LocalLeaderRegistrationRequest{}},
		{id: "", expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			if tc.id != "" {
				mockRepo.On("GetRequest", tc.id, mock.Anything).Return(tc.mockReq, tc.mockReqErr)
			}
			res, err := service.GetRequest(tc.id, ctx)
			if tc.id == "" {
				assert.Nil(t, res)
			} else {
				assert.Equal(t, tc.mockReq, res)
			}
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetRequests(t *testing.T) {
	mockRepo := new(mockLeaderRequestRepo)
	service := &leaderRequestService{leaderRequestRepo: mockRepo}
	ctx := context.Background()

	tcs := []struct {
		req       request.GetNormalStaffRegistrationRequests
		mockData  []entities.LocalLeaderRegistrationRequest
		mockPages int
		mockErr   error
	}{
		{req: request.GetNormalStaffRegistrationRequests{GetAdminRegistrationRequets: request.GetAdminRegistrationRequets{Page: 2}}, mockData: []entities.LocalLeaderRegistrationRequest{}, mockPages: 5, mockErr: nil},
		{req: request.GetNormalStaffRegistrationRequests{GetAdminRegistrationRequets: request.GetAdminRegistrationRequets{Page: 0}}, mockData: []entities.LocalLeaderRegistrationRequest{}, mockPages: 0, mockErr: nil},
		{req: request.GetNormalStaffRegistrationRequests{GetAdminRegistrationRequets: request.GetAdminRegistrationRequets{Page: 1}}, mockData: []entities.LocalLeaderRegistrationRequest{}, mockPages: 0, mockErr: errors.New("db error")},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			expectedPage := tc.req.Page
			if expectedPage < 1 {
				expectedPage = 1
			}

			repoReq := tc.req
			repoReq.Page = expectedPage

			mockRepo.On("GetRegistrationRequests", repoReq, mock.Anything).Return(tc.mockData, tc.mockPages, tc.mockErr)
			res, err := service.GetRequests(tc.req, ctx)

			assert.Equal(t, tc.mockErr, err)
			assert.Equal(t, expectedPage, res.Page)
			assert.Equal(t, tc.mockPages, res.TotalPages)
		})
	}
}

func TestGetWalletRequests(t *testing.T) {
	mockRepo := new(mockLeaderRequestRepo)
	service := &leaderRequestService{leaderRequestRepo: mockRepo}
	ctx := context.Background()

	tcs := []struct {
		id          string
		mockData    []entities.LocalLeaderRegistrationRequest
		expectedErr error
	}{
		{id: sampleAddress, mockData: []entities.LocalLeaderRegistrationRequest{{}}, expectedErr: nil},
		{id: "invalid", expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			if tc.id == sampleAddress {
				mockRepo.On("GetWalletRegistrationRequests", tc.id, mock.Anything).Return(tc.mockData, nil)
			}
			_, err := service.GetWalletRequests(tc.id, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestVoteRequest(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	mockRepo := new(mockLeaderRequestRepo)
	service := &leaderRequestService{
		leaderRequestRepo: mockRepo,
		clients:           map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		errLogger:         util.GetLogConfig(shared.ERROR_LEVEL),
	}
	ctx := context.WithValue(context.Background(), "address", sampleAddress)

	futureTime := time.Now().Add(time.Hour)
	pastTime := time.Now().Add(-time.Hour)

	tcs := []struct {
		id            string
		req           request.VoteRequest
		mockReq       *entities.LocalLeaderRegistrationRequest
		mockManage    models.SuiObjectResponse
		mockGetErr    error
		mockManageErr error
		expectedErr   error
	}{
		{ // UTCID01: Valid Vote Yes
			id: "req1", req: request.VoteRequest{IsVoteYes: true},
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other", ClosedAt: futureTime}}},
			mockManage: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
		},
		{ // UTCID02: Valid Vote No
			id: "req1", req: request.VoteRequest{IsVoteYes: false, RefuseReason: "nope"},
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other", ClosedAt: futureTime}}},
			mockManage: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
		},
		{ // UTCID03: Expired Voting Period
			id: "req1", mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other", ClosedAt: pastTime}}},
			expectedErr: errors.New(noti.REQUEST_CLOSED_MESSAGE),
		},
		{ // UTCID04: Owner Attempts Vote
			id: "req1", mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: sampleAddress, ClosedAt: futureTime}}},
			expectedErr: errors.New(noti.OWNER_VOTE_WARN_MSG),
		},
		{ // UTCID05: Redundant Voter Attempt
			id: "req1", mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other", ClosedAt: futureTime, Approvers: []string{sampleAddress}}}},
			expectedErr: errors.New(noti.ALREADY_VOTE_MESSAGE),
		},
		{ // UTCID06: Unprivileged Caller
			id: "req1", req: request.VoteRequest{IsVoteYes: true},
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other", ClosedAt: futureTime}}},
			mockManage: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{"other"}}}}}},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // UTCID07: Refusal Without Reason
			id: "req1", req: request.VoteRequest{IsVoteYes: false, RefuseReason: ""},
			mockReq: &entities.LocalLeaderRegistrationRequest{VolunteerRegistrationRequest: entities.VolunteerRegistrationRequest{AdminRegistrationRequest: entities.AdminRegistrationRequest{CreatedBy: "other", ClosedAt: futureTime}}},
			mockManage: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"admin_ids": []interface{}{sampleAddress}}}}}},
			expectedErr: errors.New(noti.FIELD_EMPTY_WARN_MSG),
		},
		{ // UTCID08: Request Not Found
			id: "req1", mockReq: nil, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockClient.ExpectedCalls = nil

			mockRepo.On("GetRequest", tc.id, mock.Anything).Return(tc.mockReq, tc.mockGetErr)

			if tc.mockReq != nil && !tc.mockReq.ClosedAt.Before(time.Now()) && tc.mockReq.CreatedBy != sampleAddress && !slices.Contains(tc.mockReq.Approvers, sampleAddress) && !slices.Contains(tc.mockReq.Refusers, sampleAddress) {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.mockManage, nil)
			}

			if tc.expectedErr == nil {
				mockRepo.On("UpdateRegistrationRequest", mock.Anything, mock.Anything).Return(nil)
			}

			err := service.VoteRequest(tc.id, tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
