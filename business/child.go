package business

import (
	"context"
	"errors"
	"log"
	"math"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	internal_sui "raise-child/constants/on-chain/sui"
	"sort"
	"strings"
	"time"

	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"
	"slices"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type childService struct {
	clients   map[string]sui.ISuiAPI
	errLogger *log.Logger
}

func InitializeChildService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IChildService {
	return &childService{
		clients:   clients,
		errLogger: errLogger,
	}
}

func GenerateChildService() (business.IChildService, error) {
	return InitializeChildService(_networkAliases, util.GetLogConfig(shared.ERROR_LEVEL)), nil
}

const (
	child_records_limit int = 10
)

// GetChild implements business.IChildService.
func (c *childService) GetChild(id string, ctx context.Context) (response.ChildResponse, error) {
	if utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.ChildResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: c.errLogger,
	}, ctx)

	var res response.ChildResponse = child.ToChildResponse()
	if len(res.DynamicFields) > 0 {
		// Has dynamic fields
		if dynamicValues, _ := on_chain.GetDynamicFields(id, client, c.errLogger, ctx); dynamicValues != nil {
			res.DynamicValues = dynamicValues
		}
	}

	return res, err
}

// GetChilds implements business.IChildService.
func (c *childService) GetChildren(req request.GetChildrenRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var client = c.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	children, err := on_chain.GetOnChainObjects[entities.Child](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.ChildIds,
		ErrLogger: c.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if children == nil {
		return response.PaginationDataResponse{}, nil
	}

	var page int = req.Page
	if page < 1 {
		page = 1
	}

	var keyword string = util.StanderizeString(req.Keyword)
	var region string = util.StanderizeString(req.Region)
	var filteredChildren []entities.Child
	for i := len(children) - 1; i >= 0; i-- {
		var child entities.Child = children[i]

		if region != "" {
			if util.StanderizeString(child.Region) != region { // Not matched
				continue
			}
		}

		if req.Gender != "" {
			if child.Gender != req.Gender { // Not matched
				continue
			}
		}

		if req.YearOfBirth != nil {
			var dob time.Time = util.RawDateToTime(child.DateOfBirth)
			if dob.Year() != *req.YearOfBirth { // Not matched
				continue
			}
		}

		if keyword != "" {
			var firstName string = util.StanderizeString(child.FirstName)
			var lastName string = util.StanderizeString(child.LastName)
			if !strings.Contains(firstName, keyword) && !strings.Contains(lastName, keyword) && !strings.Contains(child.IdentityCode, keyword) { // Not matched
				continue
			}
		}

		filteredChildren = append(filteredChildren, child)
	}

	if req.SortOrder != "" {
		sort.Slice(filteredChildren, func(i, j int) bool {
			var name1 string = filteredChildren[i].LastName + " " + filteredChildren[i].FirstName
			var name2 string = filteredChildren[j].LastName + " " + filteredChildren[j].FirstName

			if req.SortOrder == "asc" {
				return name1 < name2
			}

			return name2 > name1
		})
	}

	var skippedRecords int = (page - 1) * child_records_limit
	if len(filteredChildren) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	var data []response.ChildResponse
	for i := skippedRecords; i < len(filteredChildren); i++ {
		data = append(data, filteredChildren[i].ToChildResponse())
	}

	return response.PaginationDataResponse{
		Data:       data,
		Page:       page,
		TotalPages: int(math.Ceil(float64(len(filteredChildren)) / float64(child_records_limit))),
	}, nil
}

// UploadChild implements business.IChildService.
func (c *childService) UploadChild(req request.UploadChildRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	// todo: validate if this child is existed or not
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var rawDate string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(rawDate); dob.IsZero() { // Invalid date
		return response.BuildTransactionResponse{}, genericErr
	}

	var gender string = util.StanderizeGender(req.Gender)
	if gender == "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	var module = on_chain.InitializeModuleChild()
	res, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    c.clients[constant.SuiTestnet],
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionAddChild(),
		ErrLogger: c.errLogger,
		Arguments: module.ToAddChildArguments(on_chain.AddChildArguments{
			IdentityCode: util.StanderizeString(req.IdentityCode),
			FirstName:    util.StanderizeString(req.FirstName),
			LastName:     util.StanderizeString(req.LastName),
			Gender:       gender,
			DateOfBirth:  rawDate,
			AvatarBlobId: util.StanderizeString(req.AvatarBlobId),
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: res,
	}, err

}

// AddNumberMetada implements business.IChildService.
func (c *childService) AddNumberMetada(id string, req request.AddChildNumberMetadaRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := getOnChainObject[entities.Child](client, id, c.errLogger, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var key string = util.StanderizeString(req.Key)
	if existed := slices.Contains(child.DynamicFields, key); existed { // Field existed
		return response.BuildTransactionResponse{}, errors.New(noti.METADATA_EXISTED_MESSAGE)
	}

	var module = on_chain.InitializeModuleChild()
	res, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionAddNumberMetadata(),
		ErrLogger: c.errLogger,
		Arguments: []interface{}{
			id,
			key,
			req.Value,
			internal_sui.CLOCK_OBJECT_ID,
		},
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: res,
	}, err
}

// AddStringMetada implements business.IChildService.
func (c *childService) AddStringMetada(id string, req request.AddChildStringMetadaRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.BuildTransactionResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = c.clients[constant.SuiTestnet]
	child, err := getOnChainObject[entities.Child](client, id, c.errLogger, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var key string = util.StanderizeString(req.Key)
	if existed := slices.Contains(child.DynamicFields, key); existed { // Field existed
		return response.BuildTransactionResponse{}, errors.New(noti.METADATA_EXISTED_MESSAGE)
	}

	var module = on_chain.InitializeModuleChild()
	res, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionAddStringMetadata(),
		ErrLogger: c.errLogger,
		Arguments: []interface{}{
			id,
			key,
			util.StanderizeString(req.Value),
			internal_sui.CLOCK_OBJECT_ID,
		},
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: res,
	}, err
}
