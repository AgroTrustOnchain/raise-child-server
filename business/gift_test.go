package business

import (
	"context"
	"errors"
	"fmt"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
	"raise-child/util"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetGift(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var ctx = context.Background()
	var tcsInfo = []struct {
		id            string
		isCallGetGift bool
		giftRes       models.SuiObjectResponse
		expectedErr   error
	}{
		{ // Happy case
			id:            sampleAddress,
			isCallGetGift: true,
			giftRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{
							Fields: map[string]interface{}{
								"sender": "",
							},
						},
					},
				},
			},
		},
		{ // Invalid format id case
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
		{ // Not exist gift case
			id:            sampleAddress,
			isCallGetGift: true,
			giftRes: models.SuiObjectResponse{
				Data: &models.SuiObjectData{
					Content: &models.SuiParsedData{
						SuiMoveObject: models.SuiMoveObject{},
					},
				},
			},
			expectedErr: errors.New(noti.GENERIC_ERROR_WARN_MSG),
		},
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

func TestGetGiftsOfRegion(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()
	var service = initializeGiftService(
		nil,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)
}
