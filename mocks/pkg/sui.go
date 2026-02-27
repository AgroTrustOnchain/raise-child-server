package pkg

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/mock"
)

type suiMockApi struct {
	mock.Mock
	sui.ISuiAPI
}

func InializeSuiMockApi() *suiMockApi {
	return &suiMockApi{}
}

func (s *suiMockApi) SuiGetObject(ctx context.Context, req models.SuiGetObjectRequest) (models.SuiObjectResponse, error) {
	args := s.Called(ctx, req)
	if args.Get(0) == nil {
		return models.SuiObjectResponse{}, args.Error(1)
	}
	return args.Get(0).(models.SuiObjectResponse), args.Error(1)

	// var mockData = s.Called(ctx, req)
	// var res1 models.SuiObjectResponse
	// if mockFunc, ok := mockData.Get(0).(func(context.Context, models.SuiGetObjectRequest) models.SuiObjectResponse); ok {
	// 	res1 = mockFunc(ctx, req)
	// } else {
	// 	res1 = mockData.Get(0).(models.SuiObjectResponse)
	// }

	// var res2 error
	// if mockFunc, ok := mockData.Get(1).(func(context.Context, models.SuiGetObjectRequest) error); ok {
	// 	res2 = mockFunc(ctx, req)
	// } else {
	// 	res2 = mockData.Get(1).(error)
	// }

	// return res1, res2
}

func (s *suiMockApi) SuiMultiGetObjects(ctx context.Context, req models.SuiMultiGetObjectsRequest) ([]*models.SuiObjectResponse, error) {
	// var mockData = s.Called(ctx, req)

	// var res1 []*models.SuiObjectResponse
	// if mockFunc, ok := mockData.Get(0).(func(context.Context, models.SuiMultiGetObjectsRequest) []*models.SuiObjectResponse); ok {
	// 	res1 = mockFunc(ctx, req)
	// } else {
	// 	res1 = mockData.Get(0).([]*models.SuiObjectResponse)
	// }

	// var res2 error
	// if mockFunc, ok := mockData.Get(1).(func(context.Context, models.SuiMultiGetObjectsRequest) error); ok {
	// 	res2 = mockFunc(ctx, req)
	// } else {
	// 	res2 = mockData.Get(1).(error)
	// }

	// return res1, res2

	args := s.Called(ctx, req)
	var err error
	if args.Get(1) != nil {
		err = args.Error(1)
	}
	if args.Get(0) == nil {
		return nil, err
	}
	return args.Get(0).([]*models.SuiObjectResponse), err
}

func (s *suiMockApi) SuiXGetOwnedObjects(ctx context.Context, req models.SuiXGetOwnedObjectsRequest) (models.PaginatedObjectsResponse, error) {
	// var mockData = s.Called(ctx, req)

	// var res1 []*models.SuiObjectResponse
	// if mockFunc, ok := mockData.Get(0).(func(context.Context, models.SuiMultiGetObjectsRequest) []*models.SuiObjectResponse); ok {
	// 	res1 = mockFunc(ctx, req)
	// } else {
	// 	res1 = mockData.Get(0).([]*models.SuiObjectResponse)
	// }

	// var res2 error
	// if mockFunc, ok := mockData.Get(1).(func(context.Context, models.SuiMultiGetObjectsRequest) error); ok {
	// 	res2 = mockFunc(ctx, req)
	// } else {
	// 	res2 = mockData.Get(1).(error)
	// }

	// return res1, res2

	args := s.Called(ctx, req)
	var err error
	if args.Get(1) != nil {
		err = args.Error(1)
	}

	if args.Get(0) == nil {
		return models.PaginatedObjectsResponse{}, err
	}

	return args.Get(0).(models.PaginatedObjectsResponse), err
}

func (s *suiMockApi) MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error) {
	args := s.Called(ctx, req)
	var err error
	if args.Get(1) != nil {
		err = args.Error(1)
	}

	if args.Get(0) == nil {
		return models.TxnMetaData{}, err
	}

	return args.Get(0).(models.TxnMetaData), err
}

func (s *suiMockApi) BatchTransaction(ctx context.Context, req models.BatchTransactionRequest) (models.BatchTransactionResponse, error) {
	args := s.Called(ctx, req)
	var err error
	if args.Get(1) != nil {
		err = args.Error(1)
	}

	if args.Get(0) == nil {
		return models.BatchTransactionResponse{}, err
	}

	return args.Get(0).(models.BatchTransactionResponse), err
}

func (s *suiMockApi) SuiExecuteTransactionBlock(ctx context.Context, req models.SuiExecuteTransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
	args := s.Called(ctx, req)
	var err error
	if args.Get(1) != nil {
		err = args.Error(1)
	}

	if args.Get(0) == nil {
		return models.SuiTransactionBlockResponse{}, err
	}

	return args.Get(0).(models.SuiTransactionBlockResponse), err
}
