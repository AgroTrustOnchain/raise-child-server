package business

import (
	"context"
	"raise-child/constants/shared"
	"raise-child/mocks/pkg"
	"raise-child/mocks/repository"
	"raise-child/model/dtos/request"
	"raise-child/util"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAdmins(t *testing.T) {
	var repo = repository.InializeProfileMockRepo()
	var mockClient = pkg.InializeSuiMockApi()
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

	var service = initializeAdminService(
		repo,
		map[string]sui.ISuiAPI{
			constant.SuiTestnet: mockClient,
		},
		util.GetLogConfig(shared.ERROR_LEVEL),
	)

	var fullJsons = getFullJsonAdminNfts()
	var expectedFullJsonData []*models.SuiObjectResponse
	for _, json := range fullJsons {
		expectedFullJsonData = append(expectedFullJsonData, &models.SuiObjectResponse{
			Data: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: json,
					},
				},
			},
		})
	}

	keywordData, keyword := getFoundJsonAdminNftsWithKeyWord()
	var expectedDataWithKw []*models.SuiObjectResponse
	for _, json := range keywordData {
		expectedDataWithKw = append(expectedDataWithKw, &models.SuiObjectResponse{
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
			suiJsonData: expectedFullJsonData,
		},
		{
			suiJsonData: expectedFullJsonData,
			page:        2,
			isEmpty:     true,
		},
		{
			suiJsonData: expectedDataWithKw,
			keyword:     keyword,
		},
	}

	for _, tc := range tcsInfo {
		mockClient.On("SuiMultiGetObjects", mock.Anything, mock.Anything).Return(tc.suiJsonData, nil)
		res, err := service.GetAdmins(request.GetAdminsRequest{
			Keyword: tc.keyword,
			Page:    tc.page,
		}, ctx)

		assert.NoError(t, err)
		if tc.isEmpty {
			assert.Empty(t, res.Data)
		} else {
			assert.Equal(t, len(tc.suiJsonData), res.Amount)
		}
	}
}
