package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/util"
	on_chain "raise-child/util/on_chain"
	"strings"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

type sponsorService struct {
	clients   map[string]sui.ISuiAPI
	errLogger *log.Logger
}

func InitializeSponsorService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.ISponsorService {
	return &sponsorService{
		clients:   clients,
		errLogger: errLogger,
	}
}

func GenerateSponsorService() (business.ISponsorService, error) {
	return InitializeSponsorService(_networkAliases, util.GetLogConfig(shared.ERROR_LEVEL)), nil
}

const (
	sponsor_records_limit int = 10
)

// GetSponsor implements business.ISponsorService.
func (s *sponsorService) GetSponsor(id string, ctx context.Context) (response.SponsorResponse, error) {
	if !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return response.SponsorResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = s.clients[constant.SuiTestnet]

	var sponsorModule = on_chain.InitializeModuleSponsor()
	sponsors, err := on_chain.GetOnChainOwnedObjects[entities.Sponsor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: id,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), sponsorModule.GetModule(), sponsorModule.GetSponsorNftStruct()),
		ErrLogger:    s.errLogger,
	}, ctx)
	if err != nil {
		return response.SponsorResponse{}, err
	}

	var res = sponsors[0].ToSponsorResponse()
	var recordModule = on_chain.InitializeModuleRecord()
	txs, _ := on_chain.GetOnChainOwnedObjects[entities.Transaction](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: id,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), recordModule.GetModule(), recordModule.GetTransactionRecordStruct()),
		ErrLogger:    s.errLogger,
	}, ctx)

	if txs != nil && len(txs) > 0 {
		var contributions []response.TransactionResponse
		for _, tx := range txs {
			contributions = append(contributions, tx.ToTransactionResponse())
		}
		res.Contributions = contributions
	}

	// Set sponsor object ID to user wallet address
	res.ID = id

	return res, nil
}

// GetSponsors implements business.ISponsorService.
func (s *sponsorService) GetSponsors(req request.GetSponsorsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	var client = s.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	sponsors, err := on_chain.GetOnChainObjects[entities.Sponsor](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.SponsorNfts,
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if sponsors == nil {
		return response.PaginationDataResponse{}, nil
	}

	var page int = req.Page
	if page < 1 {
		page = 1
	}

	var keyword string = util.StanderizeString(req.Keyword)
	var filteredSponsors []entities.Sponsor
	for i := len(sponsors) - 1; i >= 0; i++ {
		var sponsor entities.Sponsor = sponsors[i]

		if req.Gender != "" {
			if req.Gender != sponsor.Gender {
				continue
			}
		}

		if keyword != "" {
			var firstName string = util.StanderizeString(sponsor.FirstName)
			var lastName string = util.StanderizeString(sponsor.LastName)
			var phoneNumber string = util.StanderizeString(sponsor.PhoneNumber)
			var email string = util.StanderizeString(sponsor.Email)
			if !strings.Contains(firstName, keyword) && !strings.Contains(lastName, keyword) && !strings.Contains(phoneNumber, keyword) && !strings.Contains(email, keyword) {
				continue
			}
		}

		filteredSponsors = append(filteredSponsors, sponsor)
	}

	var skippedRecords int = (page - 1) * sponsor_records_limit
	if len(filteredSponsors) <= skippedRecords {
		return response.PaginationDataResponse{}, err
	}

	var data []response.SponsorResponse
	for i := skippedRecords; i < len(filteredSponsors); i++ {
		data = append(data, filteredSponsors[i].ToSponsorResponse())
	}

	return response.PaginationDataResponse{
		Data:       data,
		Page:       page,
		TotalPages: int(math.Ceil(float64(len(filteredSponsors)) / float64(sponsor_records_limit))),
	}, nil
}
