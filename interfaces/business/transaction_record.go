package business

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
)

type ITransactionRecordService interface {
	GetTransactionRecords(req request.GetTransactionRecordsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetTransactionRecordsV2(req request.GetTransactionRecordsRequest, ctx context.Context) (response.GetTransactionRecordsResponse, error)
	GetTransactionRecord(id string, ctx context.Context) (response.TransactionResponse, error)
}
