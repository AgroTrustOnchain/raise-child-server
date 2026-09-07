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
	"raise-child/util/cache"
	on_chain "raise-child/util/on_chain"
	"strings"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type donorService struct {
	redisCache cache.IRedisCache
	clients    map[string]sui.ISuiAPI
	errLogger  *log.Logger
}

func InitializeDonorService(clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IDonorService {
	return &donorService{
		redisCache: cache.InitializeRedisCache(),
		clients:    clients,
		errLogger:  errLogger,
	}
}

func GenerateDonorService() (business.IDonorService, error) {
	return InitializeDonorService(_networkAliases, util.GetLogConfig(shared.ERROR_LEVEL)), nil
}

const (
	donor_records_limit int = 10
)

// GetDonor implements business.IDonorService.
func (s *donorService) GetDonor(id string, ctx context.Context) (response.DonorResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.DonorResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var client = s.clients[constant.SuiTestnet]
	donor, err := on_chain.GetOnChainObject[entities.Donor](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.DonorResponse{}, err
	}

	if donor == nil {
		return response.DonorResponse{}, nil
	}

	var recordModule = on_chain.InitializeModuleRecord()
	txs, _ := on_chain.GetOnChainOwnedObjects[entities.Transaction](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: donor.Owner,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), recordModule.GetModule(), recordModule.GetTransactionRecordStruct()),
		ErrLogger:    s.errLogger,
	}, ctx)

	var res response.DonorResponse = donor.ToDonorResponse()
	if len(txs) > 0 {
		var contributions []response.TransactionResponse
		for _, tx := range txs {
			contributions = append(contributions, tx.ToTransactionResponse())
		}
		res.Contributions = contributions
	}

	return res, nil
}

// GetDonors implements business.IDonorService.
func (s *donorService) GetDonors(req request.GetDonorsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse
	var client = s.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	donors, err := on_chain.GetOnChainObjects[entities.Donor](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.DonorNfts,
		ErrLogger: s.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if donors == nil {
		return response.PaginationDataResponse{}, nil
	}

	var filteredDonors []entities.Donor
	for i := len(donors) - 1; i >= 0; i-- {
		var donor entities.Donor = donors[i]

		if req.Gender != "" {
			if req.Gender != donor.Gender {
				continue
			}
		}

		if req.Keyword != "" {
			var firstName string = util.StandardizeString(donor.FirstName)
			var lastName string = util.StandardizeString(donor.LastName)
			var phoneNumber string = util.StandardizeString(donor.PhoneNumber)
			var email string = util.StandardizeString(donor.Email)
			if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(phoneNumber, req.Keyword) && !strings.Contains(email, req.Keyword) {
				continue
			}
		}

		filteredDonors = append(filteredDonors, donor)
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredDonors) <= skippedRecords {
		return response.PaginationDataResponse{}, err
	}

	var data []response.DonorResponse
	for i := skippedRecords; i < len(filteredDonors); i++ {
		data = append(data, filteredDonors[i].ToDonorResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	// ///////////////////
	// // MOCK DATA
	// var filteredDonors []response.DonorResponse
	// var donors = mockDonors
	// for i := len(donors) - 1; i >= 0; i++ {
	// 	var donor response.DonorResponse = donors[i]

	// 	if req.Gender != "" {
	// 		if req.Gender != donor.Gender {
	// 			continue
	// 		}
	// 	}

	// 	if req.Keyword != "" {
	// 		var firstName string = util.StandardizeString(donor.FirstName)
	// 		var lastName string = util.StandardizeString(donor.LastName)
	// 		var phoneNumber string = util.StandardizeString(donor.PhoneNumber)
	// 		var email string = util.StandardizeString(donor.Email)
	// 		if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(phoneNumber, req.Keyword) && !strings.Contains(email, req.Keyword) {
	// 			continue
	// 		}
	// 	}

	// 	filteredDonors = append(filteredDonors, donor)
	// }

	// var skippedRecords int = (req.Page - 1) * req.PageSize
	// if len(filteredDonors) <= skippedRecords {
	// 	return response.PaginationDataResponse{}, nil
	// }

	// var data []response.DonorResponse
	// for i := skippedRecords; i < len(filteredDonors); i++ {
	// 	data = append(data, filteredDonors[i])
	// 	if len(data) == req.PageSize {
	// 		break
	// 	}
	// }

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredDonors)) / float64(req.PageSize))),
	}

	return res, nil
}

func (s *donorService) getGetDonorsRedisKey(req request.GetDonorsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	return fmt.Sprintf("donor:kw:%s:g:%s:s:%d:p:%d", keyword, gender, req.PageSize, req.Page)
}

var mockDonors = []response.DonorResponse{
	{ID: "1", FirstName: "Minh", LastName: "Anh", Gender: "Male", PhoneNumber: "0901234567", Email: "minh.anh@example.com", Name: "Minh Anh", TotalDonation: 500000, Url: "https://api.example.com"},
	{ID: "2", FirstName: "Lan", LastName: "Huong", Gender: "Female", PhoneNumber: "0902345678", Email: "lan.huong@example.com", Name: "Lan Huong", TotalDonation: 1200000, Url: "https://api.example.com"},
	{ID: "3", FirstName: "Hoang", LastName: "Nam", Gender: "Male", PhoneNumber: "0903456789", Email: "hoang.nam@example.com", Name: "Hoang Nam", TotalDonation: 300000, Url: "https://api.example.com"},
	{ID: "4", FirstName: "Thu", LastName: "Thuy", Gender: "Female", PhoneNumber: "0904567890", Email: "thu.thuy@example.com", Name: "Thu Thuy", TotalDonation: 2500000, Url: "https://api.example.com"},
	{ID: "5", FirstName: "Tuan", LastName: "Tu", Gender: "Male", PhoneNumber: "0905678901", Email: "tuan.tu@example.com", Name: "Tuan Tu", TotalDonation: 750000, Url: "https://api.example.com"},
	{ID: "6", FirstName: "Ngoc", LastName: "Bich", Gender: "Female", PhoneNumber: "0906789012", Email: "ngoc.bich@example.com", Name: "Ngoc Bich", TotalDonation: 100000, Url: "https://api.example.com"},
	{ID: "7", FirstName: "Duy", LastName: "Manh", Gender: "Male", PhoneNumber: "0907890123", Email: "duy.manh@example.com", Name: "Duy Manh", TotalDonation: 4200000, Url: "https://api.example.com"},
	{ID: "8", FirstName: "Phuong", LastName: "Thao", Gender: "Female", PhoneNumber: "0908901234", Email: "phuong.thao@example.com", Name: "Phuong Thao", TotalDonation: 900000, Url: "https://api.example.com"},
	{ID: "9", FirstName: "Quoc", LastName: "Bao", Gender: "Male", PhoneNumber: "0909012345", Email: "quoc.bao@example.com", Name: "Quoc Bao", TotalDonation: 150000, Url: "https://api.example.com"},
	{ID: "10", FirstName: "Mai", LastName: "Chi", Gender: "Female", PhoneNumber: "0910123456", Email: "mai.chi@example.com", Name: "Mai Chi", TotalDonation: 3300000, Url: "https://api.example.com0"},
	{ID: "11", FirstName: "Huu", LastName: "Thang", Gender: "Male", PhoneNumber: "0911234567", Email: "huu.thang@example.com", Name: "Huu Thang", TotalDonation: 600000, Url: "https://api.example.com1"},
	{ID: "12", FirstName: "Kim", LastName: "Yen", Gender: "Female", PhoneNumber: "0912345678", Email: "kim.yen@example.com", Name: "Kim Yen", TotalDonation: 2100000, Url: "https://api.example.com2"},
	{ID: "13", FirstName: "Viet", LastName: "Bach", Gender: "Male", PhoneNumber: "0913456789", Email: "viet.bach@example.com", Name: "Viet Bach", TotalDonation: 450000, Url: "https://api.example.com3"},
	{ID: "14", FirstName: "Hanh", LastName: "Phuc", Gender: "Female", PhoneNumber: "0914567890", Email: "hanh.phuc@example.com", Name: "Hanh Phuc", TotalDonation: 880000, Url: "https://api.example.com4"},
	{ID: "15", FirstName: "Gia", LastName: "Huy", Gender: "Male", PhoneNumber: "0915678901", Email: "gia.huy@example.com", Name: "Gia Huy", TotalDonation: 1750000, Url: "https://api.example.com5"},
	{ID: "16", FirstName: "Thanh", LastName: "Tra", Gender: "Female", PhoneNumber: "0916789012", Email: "thanh.tra@example.com", Name: "Thanh Tra", TotalDonation: 550000, Url: "https://api.example.com6"},
	{ID: "17", FirstName: "Son", LastName: "Tung", Gender: "Male", PhoneNumber: "0917890123", Email: "son.tung@example.com", Name: "Son Tung", TotalDonation: 10000000, Url: "https://api.example.com7"},
	{ID: "18", FirstName: "Bao", LastName: "Ngoc", Gender: "Female", PhoneNumber: "0918901234", Email: "bao.ngoc@example.com", Name: "Bao Ngoc", TotalDonation: 125000, Url: "https://api.example.com8"},
	{ID: "19", FirstName: "Trong", LastName: "Nghia", Gender: "Male", PhoneNumber: "0919012345", Email: "trong.nghia@example.com", Name: "Trong Nghia", TotalDonation: 2200000, Url: "https://api.example.com9"},
	{ID: "20", FirstName: "Dieu", LastName: "Linh", Gender: "Female", PhoneNumber: "0920123456", Email: "dieu.linh@example.com", Name: "Dieu Linh", TotalDonation: 3100000, Url: "https://api.example.com0"},
}
