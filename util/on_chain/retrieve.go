package onchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"raise-child/constants/noti"
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strings"
	"sync"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/machinebox/graphql"
)

const (
	default_chunk_size int = 50
)

type GetOnChainObjectRequest struct {
	Client    sui.ISuiAPI
	ObjectId  string
	ErrLogger *log.Logger
}

type GetOnChainObjectsRequest struct {
	Client    sui.ISuiAPI
	ObjectIds []string
	ErrLogger *log.Logger
}

type GetOnChainOwnedObjectsRequest struct {
	Client       sui.ISuiAPI
	OwnerAddress string
	StructType   string
	ErrLogger    *log.Logger
}

type GetOnChainSpecificTypeObjectsRequest struct {
	Client     *graphql.Client
	StructType string
	ErrLogger  *log.Logger
}

func GetOnChainObject[T any](req GetOnChainObjectRequest, ctx context.Context) (*T, error) {
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	var retrieveReq = models.SuiGetObjectRequest{
		ObjectId: req.ObjectId,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	}

	res, err := req.Client.SuiGetObject(ctx, retrieveReq)
	if err != nil {
		req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
		return nil, internalErr
	}

	var jsonData = res.Data.Content.Fields
	if jsonData == nil {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
		return nil, internalErr
	}

	var object = util.JsonStringToObject[T](string(jsonBytes))

	return &object, nil
}

func GetOnChainObjects[T any](req GetOnChainObjectsRequest, ctx context.Context) ([]T, error) {
	// VER 1
	// var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	// var retrieveReq = models.SuiMultiGetObjectsRequest{
	// 	ObjectIds: req.ObjectIds,
	// 	Options: models.SuiObjectDataOptions{
	// 		ShowContent: true,
	// 	},
	// }

	// res, err := req.Client.SuiMultiGetObjects(ctx, retrieveReq)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "no result") {
	// 		return nil, nil
	// 	}

	// 	req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
	// 	return nil, internalErr
	// }

	// var objects []T
	// for i := 0; i < len(res); i++ {
	// 	var object = res[i]
	// 	if object.Error != nil {
	// 		handleGetEmptyOnChainObject(object.Error.Code, req.ObjectIds[i], req.ErrLogger)
	// 		continue
	// 	}

	// 	jsonBytes, err := json.Marshal(object.Data.Content.Fields)
	// 	if err != nil {
	// 		req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
	// 		continue
	// 	}

	// 	objects = append(objects, util.JsonStringToObject[T](string(jsonBytes)))
	// }

	// return objects, nil

	// VER 2
	var defaultWorkersNum int = 5
	var idChunks = make(chan []string, len(req.ObjectIds)/default_chunk_size+1)
	var suiRes = make(chan []*models.SuiObjectResponse, len(req.ObjectIds)/default_chunk_size+1)
	var wg sync.WaitGroup

	for i := 1; i < defaultWorkersNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for batch := range idChunks {
				res, err := req.Client.SuiMultiGetObjects(ctx, models.SuiMultiGetObjectsRequest{
					ObjectIds: batch,
					Options: models.SuiObjectDataOptions{
						ShowContent: true,
					},
				})

				if err != nil {
					if strings.Contains(err.Error(), "no result") {
						break
					}

					req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
					continue
				}

				suiRes <- res
			}
		}()
	}

	for i := 0; i < len(req.ObjectIds); i += default_chunk_size {
		var end int = i + default_chunk_size
		if end > len(req.ObjectIds) {
			end = len(req.ObjectIds)
		}

		idChunks <- req.ObjectIds[i:end]
	}
	close(idChunks)

	go func() {
		wg.Wait()
		close(suiRes)
	}()

	var objects []T
	if suiRes == nil {
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	for batchRes := range suiRes {
		for _, object := range batchRes {
			jsonBytes, err := json.Marshal(object.Data.Content.Fields)
			if err != nil {
				req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
				continue
			}

			objects = append(objects, util.JsonStringToObject[T](string(jsonBytes)))
		}
	}

	return objects, nil
}

func GetOnChainOwnedObjects[T any](req GetOnChainOwnedObjectsRequest, ctx context.Context) ([]T, error) {
	// VER 1
	// var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	// var filter = map[string]interface{}{
	// 	"StructType": req.StructType,
	// }
	// var retrieveReq = models.SuiXGetOwnedObjectsRequest{
	// 	Address: req.OwnerAddress,
	// 	Query: models.SuiObjectResponseQuery{
	// 		Filter: filter,
	// 		Options: models.SuiObjectDataOptions{
	// 			ShowContent: true,
	// 		},
	// 	},
	// }

	// res, err := req.Client.SuiXGetOwnedObjects(ctx, retrieveReq)
	// if err != nil {
	// 	req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
	// 	return nil, internalErr
	// }

	// var objects []T
	// for i := 0; i < len(res.Data); i++ {
	// 	var object = res.Data[i]
	// 	if object.Error != nil {
	// 		handleGetEmptyOnChainObject(object.Error.Code, object.Data.ObjectId, req.ErrLogger)
	// 		continue
	// 	}

	// 	jsonBytes, err := json.Marshal(object.Data.Content.Fields)
	// 	if err != nil {
	// 		req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
	// 		continue
	// 	}

	// 	objects = append(objects, util.JsonStringToObject[T](string(jsonBytes)))
	// }

	// return objects, nil

	// VER 2
	var objects []T
	var cursor interface{} = nil

	for {
		var retrieveReq = models.SuiXGetOwnedObjectsRequest{
			Address: req.OwnerAddress,
			Query: models.SuiObjectResponseQuery{
				Filter: map[string]interface{}{
					"StructType": req.StructType,
				},
				Options: models.SuiObjectDataOptions{
					ShowContent: true,
				},
			},
			Cursor: cursor,
			Limit:  uint64(default_chunk_size),
		}

		res, err := req.Client.SuiXGetOwnedObjects(ctx, retrieveReq)
		if err != nil {
			if strings.Contains(err.Error(), "no result") {
				return nil, nil
			}

			req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
			return nil, errors.New(noti.INTERNALL_ERR_MSG)
		}

		for _, object := range res.Data {
			if object.Error != nil {
				handleGetEmptyOnChainObject(object.Error.Code, object.Data.ObjectId, req.ErrLogger)
				continue
			}

			jsonBytes, err := json.Marshal(object.Data.Content.Fields)
			if err != nil {
				continue
			}
			objects = append(objects, util.JsonStringToObject[T](string(jsonBytes)))
		}

		if !res.HasNextPage || res.NextCursor == "" {
			break
		}

		cursor = res.NextCursor
	}

	return objects, nil
}

func GetDynamicFields(id string, client sui.ISuiAPI, errLogger *log.Logger, ctx context.Context) (map[string]interface{}, error) {
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	fields, err := client.SuiXGetDynamicField(ctx, models.SuiXGetDynamicFieldRequest{
		ObjectId: id,
	})

	if err != nil {
		errLogger.Println(noti.ApproRETRIEVE_DYNAMIC_FIELDS_ERR_MSGvers + err.Error())
		return nil, internalErr
	}

	var res map[string]interface{}
	for _, fieldInfo := range fields.Data {
		fieldValue, err := client.SuiXGetDynamicFieldObject(ctx, models.SuiXGetDynamicFieldObjectRequest{
			ObjectId: fieldInfo.ObjectId,
			DynamicFieldName: models.DynamicFieldObjectName{
				Type:  fieldInfo.ObjectType,
				Value: fieldInfo.Name.Value,
			},
		})

		if err != nil {
			errLogger.Println(noti.ApproRETRIEVE_DYNAMIC_FIELDS_ERR_MSGvers + err.Error())
			return nil, internalErr
		}

		var field = fieldValue.Data.Content.SuiMoveObject.Fields
		for k, v := range field {
			res[k] = v
		}
	}

	return res, nil
}

// func GetOnChainEvents[T any](req GetOnChainOwnedObjectsRequest, ctx context.Context) {
// 	res, err := req.Client.SuiXQueryEvents(ctx, models.SuiXQueryEventsRequest{})
// }

func GetOnChainSpecificTypeObjects[T any](req GetOnChainSpecificTypeObjectsRequest, ctx context.Context) ([]T, error) {
	var retrieveReq = graphql.NewRequest(`
		query ($type: String!) {
			objects(filter: { type: $type }) {
				nodes {
					address
					asMoveObject {
						contents {
							json
						}
					}
				}
			}
		}
	`)
	retrieveReq.Var("type", req.StructType)

	var res response.SuiGraphQlObjectResponse
	if err := req.Client.Run(ctx, retrieveReq, &res); err != nil {
		req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var objects []T
	for _, node := range res.Objects.Nodes {
		var jsonObject = node.AsMoveObject.Contents.Json

		jsonBytes, err := json.Marshal(jsonObject)
		if err != nil {
			req.ErrLogger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + err.Error())
			continue
		}

		objects = append(objects, util.JsonStringToObject[T](string(jsonBytes)))
	}

	return objects, nil
}

func handleGetEmptyOnChainObject(err, id string, logger *log.Logger) {
	var msg string
	switch err {
	case "deleted":
		msg = fmt.Sprintf("Object %s was deleted from the network.", id)
		break
	case "notExists":
		msg = fmt.Sprintf("Object %s does not exist.", id)
		break
	}

	if msg != "" {
		logger.Println(noti.RETRIEVE_ON_CHAIN_DATA_ERR_MSG + msg)
	}
}
