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

func TestGetCenterRequests(t *testing.T) {
	var centerReqrepo = repository.InitializeCenterRequestMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeCenterRequestService(
		centerReqrepo,
		profileRepo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.Background()
	var tcsInfo = []struct {
		req               request.GetCenterRequests
		isDefaultPageSize bool
		expectedErr       error
	}{
		{
			expectedErr: errors.New(noti.INTERNALL_ERR_MSG),
		},
		{
			isDefaultPageSize: true,
		},
		{
			req: request.GetCenterRequests{
				PageSize: 7,
			},
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			centerReqrepo.ExpectedCalls = nil

			if tc.expectedErr != nil {
				centerReqrepo.On("GetRegistrationRequests", mock.Anything, mock.Anything).Return(nil, 0, tc.expectedErr)
				_, err := service.GetRequests(tc.req, ctx)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				var reqsNum int
				if tc.isDefaultPageSize {
					reqsNum = default_page_size
				} else {
					reqsNum = tc.req.PageSize
				}

				centerReqrepo.On("GetRegistrationRequests", mock.Anything, mock.Anything).Return(getSampleCenterRequests(reqsNum), 1, nil)
				res, _ := service.GetRequests(tc.req, ctx)
				assert.Equal(t, reqsNum, res.Amount)
			}
		})
	}
}

func TestGetWalletCenterRequests(t *testing.T) {
	var centerReqrepo = repository.InitializeCenterRequestMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeCenterRequestService(
		centerReqrepo,
		profileRepo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.Background()
	var tcsInfo = []struct {
		id          string
		isHaveReqs  bool
		expectedErr error
	}{
		{
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{
			id: sampleAddress,
		},
		{
			id:         sampleAddress,
			isHaveReqs: true,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			centerReqrepo.ExpectedCalls = nil

			var expectedRes []entities.CenterRequest
			if tc.isHaveReqs {
				expectedRes = getSampleCenterRequests(1)
			}

			centerReqrepo.On("GetWalletRegistrationRequests", mock.Anything, mock.Anything).Return(expectedRes, tc.expectedErr)
			res, err := service.GetWalletRequests(tc.id, ctx)
			assert.Equal(t, expectedRes, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetCenterRequest(t *testing.T) {
	var centerReqrepo = repository.InitializeCenterRequestMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeCenterRequestService(
		centerReqrepo,
		profileRepo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.Background()
	var tcsInfo = []struct {
		isExist     bool
		expectedErr error
	}{
		{
			isExist: true,
		},
		{
			isExist: false,
		},
		{
			expectedErr: errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			centerReqrepo.ExpectedCalls = nil

			var expectedRes *entities.CenterRequest
			if tc.isExist {
				expectedRes = &entities.CenterRequest{}
			}

			centerReqrepo.On("GetRequest", mock.Anything, mock.Anything).Return(expectedRes, tc.expectedErr)
			res, err := service.GetRequest("", ctx)
			assert.Equal(t, expectedRes, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

// func TestCreateCenterRequest(t *testing.T) {
// 	var centerReqrepo = repository.InitializeCenterRequestMockRepo()
// 	var profileRepo = repository.InitializeProfileMockRepo()
// 	var mockClient = pkg.InitializeSuiMockApi()
// 	var service = initializeCenterRequestService(
// 		centerReqrepo,
// 		profileRepo,
// 		map[string]sui.ISuiAPI{
// 			constant.SuiTestnet: mockClient,
// 		},
// 		util.GetLogConfig(shared.ERROR_LEVEL),
// 	)

// 	var tcsInfo = []struct {
// 		isExist     bool
// 		expectedErr error
// 		context     context.Context
// 	}{}
// }

func TestCreateCenterRequest(t *testing.T) {
	var centerReqrepo = repository.InitializeCenterRequestMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeCenterRequestService(
		centerReqrepo,
		profileRepo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", "sample_sub")

	var sampleStaffsJson []*models.SuiObjectResponse = []*models.SuiObjectResponse{
		&models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: sampleJsonLeaderNft1,
					},
				},
			},
		},
		&models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: sampleJsonLeaderNft1,
					},
				},
			},
		},
		&models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: sampleJsonLeaderNft1,
					},
				},
			},
		},
		&models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: sampleJsonLeaderNft1,
					},
				},
			},
		},
		&models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: sampleJsonLeaderNft1,
					},
				},
			},
		},
		&models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: sampleJsonLeaderNft1,
					},
				},
			},
		},
	}
	var tcsInfo = []struct {
		req                 request.CreateCenterRequest
		ownedNftsJson       models.PaginatedObjectsResponse
		manageJsonObj       models.SuiObjectResponse
		staffsJson          []*models.SuiObjectResponse
		isHaveToFetchManage bool
		isHaveToFetchStaffs bool
		expectedErr         error
		expectedRes         *entities.CenterRequest
	}{
		{ // Happy case
			req: request.CreateCenterRequest{
				Region: sampleLocalRegions[0],
			},
			ownedNftsJson: models.PaginatedObjectsResponse{
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
			manageJsonObj: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonManageObj,
						},
					},
				},
			},
			staffsJson:          sampleStaffsJson,
			isHaveToFetchManage: true,
			isHaveToFetchStaffs: true,
			expectedRes:         &entities.CenterRequest{},
		},
		{ // Not staff case
			ownedNftsJson: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Different region case
			ownedNftsJson: models.PaginatedObjectsResponse{
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
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
		},
		{ // Region-established case
			req: request.CreateCenterRequest{
				Region: sampleLocalRegions[1],
			},
			ownedNftsJson: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{
					models.SuiObjectResponse{
						Data: &models.SuiObjectData{
							Content: &models.SuiParsedData{
								SuiMoveObject: models.SuiMoveObject{
									Fields: sampleJsonLeaderNft2,
								},
							},
						},
					},
				},
			},
			manageJsonObj: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonManageObj,
						},
					},
				},
			},
			isHaveToFetchManage: true,
			expectedErr:         errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not enough staffs case
			req: request.CreateCenterRequest{
				Region: sampleLocalRegions[0],
			},
			ownedNftsJson: models.PaginatedObjectsResponse{
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
			manageJsonObj: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonManageObj,
						},
					},
				},
			},
			staffsJson:          sampleStaffsJson[0:4],
			isHaveToFetchManage: true,
			isHaveToFetchStaffs: true,
			expectedErr:         errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Insert fail case
			req: request.CreateCenterRequest{
				Region: sampleLocalRegions[0],
			},
			ownedNftsJson: models.PaginatedObjectsResponse{
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
			manageJsonObj: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: sampleJsonManageObj,
						},
					},
				},
			},
			isHaveToFetchManage: true,
			isHaveToFetchStaffs: true,
			staffsJson:          sampleStaffsJson,
			expectedErr:         errors.New(noti.INTERNALL_ERR_MSG),
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			centerReqrepo.ExpectedCalls = nil

			if tc.ownedNftsJson.Data != nil {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.ownedNftsJson, nil)
			}

			if tc.isHaveToFetchManage {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.manageJsonObj, nil)
			}

			if tc.isHaveToFetchStaffs {
				mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return(tc.staffsJson, nil)
			}

			centerReqrepo.On("CreateRegistrationRequest", mock.Anything, mock.Anything).Return(tc.expectedErr)

			res, err := service.CreateRequest(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
			if tc.expectedRes != nil {
				assert.Equal(t, true, res != nil)
			}
		})
	}
}

func TestVoteCenterRequest(t *testing.T) {
	var centerReqrepo = repository.InitializeCenterRequestMockRepo()
	var profileRepo = repository.InitializeProfileMockRepo()
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeCenterRequestService(
		centerReqrepo,
		profileRepo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var creatorAddress string = "creator address"
	var votedAddress string = "voted address"
	var votedCtx = context.WithValue(context.Background(), "address", votedAddress)
	votedCtx = context.WithValue(votedCtx, "sub", "sample_sub")
	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", "sample_sub")

	var sampleCenterReq *entities.CenterRequest = &entities.CenterRequest{
		CreatedBy: creatorAddress,
		Region:    sampleLocalRegions[0],
		Approvers: []string{votedAddress},
		Refusers:  []string{votedAddress},
		ClosedAt:  time.Now().Add(time.Minute * 5),
	}

	var ownedStaffNftsJson = models.PaginatedObjectsResponse{
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

	var tcsInfo = []struct {
		centerReq     *entities.CenterRequest
		voteReq       request.VoteRequest
		ownedNftsJson models.PaginatedObjectsResponse
		expectedErr   error
		context       context.Context
	}{
		{ // Happy case - approve
			centerReq: sampleCenterReq,
			voteReq: request.VoteRequest{
				IsVoteYes: true,
			},
			ownedNftsJson: ownedStaffNftsJson,
			context:       ctx,
		},
		{ // Not exist request case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			context:     ctx,
		},
		{ // Request closed case
			centerReq: &entities.CenterRequest{
				ClosedAt: sampleCenterReq.ClosedAt.Add(time.Hour * -3),
			},
			expectedErr: errors.New(noti.REQUEST_CLOSED_MESSAGE),
			context:     ctx,
		},
		{ // Creator-vote case
			centerReq:   sampleCenterReq,
			expectedErr: errors.New(noti.OWNER_VOTE_WARN_MSG),
			context:     context.WithValue(context.Background(), "address", creatorAddress),
		},
		{ // Already-vote case
			centerReq:   sampleCenterReq,
			expectedErr: errors.New(noti.ALREADY_VOTE_MESSAGE),
			context:     votedCtx,
		},
		{ // Not-staff case
			centerReq: sampleCenterReq,
			ownedNftsJson: models.PaginatedObjectsResponse{
				Data: []models.SuiObjectResponse{},
			},
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			context:     ctx,
		},
		{ // Not staff of region case
			centerReq: sampleCenterReq,
			ownedNftsJson: models.PaginatedObjectsResponse{
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
			expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			context:     ctx,
		},
		{ // Refuse vote with no reason case
			centerReq: sampleCenterReq,
			voteReq: request.VoteRequest{
				IsVoteYes: false,
			},
			ownedNftsJson: ownedStaffNftsJson,
			expectedErr:   errors.New(noti.FIELD_EMPTY_WARN_MSG),
			context:       ctx,
		},
		{ // Update-fail case
			centerReq: sampleCenterReq,
			voteReq: request.VoteRequest{
				IsVoteYes: true,
			},
			ownedNftsJson: ownedStaffNftsJson,
			expectedErr:   errors.New(noti.INTERNALL_ERR_MSG),
			context:       ctx,
		},
	}

	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			centerReqrepo.ExpectedCalls = nil

			centerReqrepo.On("GetRequest", mock.Anything, mock.Anything).Return(tc.centerReq, tc.expectedErr)
			if tc.ownedNftsJson.Data != nil {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(tc.ownedNftsJson, nil)
			}

			centerReqrepo.On("UpdateRegistrationRequest", mock.Anything, mock.Anything).Return(tc.expectedErr)

			assert.Equal(t, tc.expectedErr, service.VoteRequest("", tc.voteReq, tc.context))
		})
	}
}
