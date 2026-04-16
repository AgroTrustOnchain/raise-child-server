package business

import (
	"context"
	"fmt"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
	"raise-child/model/dtos/request"
	"raise-child/util"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDonors(t *testing.T) {
	var mockClient = pkg.InitializeSuiMockApi()

	var service = initializeDonorService(
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var fullJsons = getFullJsonDonorNfts()
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

	keywordData, keyword := getFoundJsonDonorNftsWithKeyWord()
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
	for i, tc := range tcsInfo {
		t.Run(fmt.Sprintf("Case-%d", i), func(t *testing.T) {
			mockClient.ExpectedCalls = nil

			mockClient.On("SuiGetObject", mock.Anything, mock.Anything).
				Return(models.SuiObjectResponse{
					Data: &models.SuiObjectData{
						Content: &models.SuiParsedData{
							SuiMoveObject: models.SuiMoveObject{
								Fields: sampleJsonManageObj,
							},
						},
					},
				}, nil)
			mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return(tc.suiJsonData, nil)
			res, err := service.GetDonors(request.GetDonorsRequest{
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
