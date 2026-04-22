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

func ptr(s string) *string {
	return &s
}

var (
	sampleGiftRes = models.SuiObjectResponse{
		Data: &models.SuiObjectData{
			Content: &models.SuiParsedData{
				SuiMoveObject: models.SuiMoveObject{
					Fields: map[string]interface{}{
						"id":        map[string]string{"id": sampleAddress},
						"sender":    sampleAddress,
						"recipient": sampleAddress,
						"status":    "Created",
					},
				},
			},
		},
	}
)

func TestGetGift(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(
		nil,
		map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)
	var ctx = context.Background()
	var tcsInfo = []struct {
		id            string
		isCallGetGift bool
		giftRes       models.SuiObjectResponse
		expectedErr   error
	}{
		{id: sampleAddress, isCallGetGift: true, giftRes: sampleGiftRes},
		{expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
		{id: sampleAddress, isCallGetGift: true, giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{}}}}, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
	}
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			if tc.isCallGetGift {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.giftRes, nil)
			}
			_, err := service.GetGift(tc.id, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestCancelGift(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(nil, map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient}, util.GetLogConfig(shared.ERROR_LEVEL))
	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	tcsInfo := []struct {
		id           string
		giftRes      models.SuiObjectResponse
		mockGet      bool
		mockMoveCall bool
		expectedErr  error
	}{
		{ // UTCID01: Normal Cancellation
			id: sampleAddress, giftRes: sampleGiftRes, mockGet: true, mockMoveCall: true,
		},
		{ // UTCID02: Invalid ID
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID03: Gift Not Found
			id: sampleAddress, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			mockGet: true, giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{}}}},
		},
		{ // UTCID04: Gift Already Delivered
			id: sampleAddress, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			mockGet: true, giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"status": "Delivered"}}}}},
		},
		{ // UTCID05: Gift Already Canceled
			id: sampleAddress, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
			mockGet: true, giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"status": "Canceled"}}}}},
		},
		{ // UTCID06: Unauthorized Sender
			id: sampleAddress, expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			mockGet: true, giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"status": "Created", "sender": "other_address"}}}}},
		},
	}
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			if tc.mockGet {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.giftRes, nil)
			}
			if tc.mockMoveCall {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{TxBytes: "bytes"}, nil)
			}
			_, err := service.CancelGift(tc.id, request.CancelGiftRequest{}, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestConfirmReceiveGift(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(nil, map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient}, util.GetLogConfig(shared.ERROR_LEVEL))
	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	tcsInfo := []struct {
		id           string
		giftRes      models.SuiObjectResponse
		rolesRes     []models.SuiObjectResponse
		childRes     models.SuiObjectResponse
		mockGet      bool
		mockOwned    bool
		mockChild    bool
		mockMoveCall bool
		expectedErr  error
	}{
		{ // UTCID01: Child Gift Confirmation Normal
			id: sampleAddress,
			giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"sender": "other", "status": "Created", "is_for_child": true, "recipient": sampleAddress}}}}},
			rolesRes: []models.SuiObjectResponse{{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"id": map[string]string{"id": sampleAddress}, "region": "North"}}}}}},
			childRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"id": map[string]string{"id": sampleAddress}, "region": "North"}}}}},
			mockGet: true, mockOwned: true, mockChild: true, mockMoveCall: true,
		},
		{ // UTCID03: Invalid ID
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID06: Sender Confirms Own Gift
			id: sampleAddress, expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			mockGet: true, giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"sender": sampleAddress, "status": "Created"}}}}},
		},
		{ // UTCID07: No Staff Role
			id: sampleAddress, expectedErr: errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG),
			mockGet: true, mockOwned: true, rolesRes: nil,
			giftRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"sender": "other", "status": "Created"}}}}},
		},
	}
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			if tc.mockGet {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.giftRes, nil).Once()
			}
			if tc.mockOwned {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(models.PaginatedObjectsResponse{Data: tc.rolesRes}, nil)
			}
			if tc.mockChild {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil)
			}
			if tc.mockMoveCall {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{TxBytes: "bytes"}, nil)
			}
			_, err := service.ConfirmReceiveGift(tc.id, request.ConfirmReceiveGiftRequest{}, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestCreateGift(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var profileRepo = repository.InitializeProfileMockRepo()
	var service = initializeGiftService(profileRepo, map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient}, util.GetLogConfig(shared.ERROR_LEVEL))
	var ctx = context.WithValue(context.Background(), "address", sampleAddress)
	ctx = context.WithValue(ctx, "sub", sampleSub)

	validProfile := &entities.Profile{IdentityCode: ptr("123"), FirstName: ptr("A"), LastName: ptr("B"), Gender: ptr("Male"), PhoneNumber: ptr("000"), Email: ptr("a@b.c")}
	tcsInfo := []struct {
		req           request.CreateGiftRequest
		profile       *entities.Profile
		childRes      models.SuiObjectResponse
		manageRes     models.SuiObjectResponse
		mockProfile   bool
		mockOwned     bool
		mockRecipient bool
		mockManage    bool
		mockMoveCall  bool
		expectedErr   error
	}{
		{ // UTCID01: Normal Create For Child
			req: request.CreateGiftRequest{Recipient: sampleAddress}, profile: validProfile,
			childRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: sampleJsonChild1}}}},
			manageRes: models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: sampleJsonManageObj}}}},
			mockProfile: true, mockOwned: true, mockRecipient: true, mockManage: true, mockMoveCall: true,
		},
		{ // UTCID03: Invalid Recipient ID
			req: request.CreateGiftRequest{Recipient: "invalid"}, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // UTCID04: Donor Profile Without KYC
			req: request.CreateGiftRequest{Recipient: sampleAddress}, profile: &entities.Profile{},
			mockProfile: true, expectedErr: errors.New(noti.PROFILE_EMPTY_MESSAGE),
		},
	}
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			profileRepo.ExpectedCalls = nil
			if tc.mockProfile {
				profileRepo.On("GetProfile", mock.Anything, mock.Anything).Return(tc.profile, nil)
			}
			if tc.mockOwned {
				mockClient.On("SuiXGetOwnedObjects", mock.Anything, mock.Anything).Return(models.PaginatedObjectsResponse{Data: []models.SuiObjectResponse{}}, nil)
			}
			if tc.mockRecipient {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.childRes, nil).Once()
			}
			if tc.mockManage {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(tc.manageRes, nil)
			}
			if tc.mockMoveCall {
				mockClient.On("MoveCall", mock.Anything, mock.Anything).Return(models.TxnMetaData{TxBytes: "bytes"}, nil)
			}
			_, err := service.CreateGift(tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetGiftsOfChild(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(nil, map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient}, util.GetLogConfig(shared.ERROR_LEVEL))
	var ctx = context.Background()
	tcsInfo := []struct {
		id          string
		req         request.GetGiftsRequest
		mockGetAuth bool
		mockMulti   bool
		expectedErr error
	}{
		{id: sampleAddress, req: request.GetGiftsRequest{}, mockGetAuth: true},
		{id: "invalid", expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
	}
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			if tc.mockGetAuth {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"gifts": []interface{}{sampleAddress}}}}}}, nil)
				mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return([]*models.SuiObjectResponse{}, nil)
			}
			_, err := service.GetGiftsOfChild(tc.id, tc.req, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestGetGiftsOfRegion(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(nil, map[string]sui.ISuiAPI{constant.SuiTestnet: mockClient}, util.GetLogConfig(shared.ERROR_LEVEL))
	var ctx = context.Background()
	tcsInfo := []struct {
		region      string
		mockManage  bool
		expectedErr error
	}{
		{region: "invalid region", mockManage: true, expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG)},
		{region: sampleLocalRegions[0], mockManage: true},
	}
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil
			if tc.mockManage {
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: sampleJsonManageObj}}}}, nil).Once()
				mockClient.On("SuiGetObject", mock.Anything, mock.Anything).Return(models.SuiObjectResponse{Data: &models.SuiObjectData{Content: &models.SuiParsedData{SuiMoveObject: models.SuiMoveObject{Fields: map[string]interface{}{"all_gifts": []interface{}{}}}}}}, nil)
			}
			_, err := service.GetGiftsOfRegion(tc.region, request.GetGiftsRequest{}, ctx)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
