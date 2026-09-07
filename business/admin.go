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
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"sort"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type adminService struct {
	profileRepo i_repository.IProfileRepository
	redisCache  cache.IRedisCache
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func initializeAdminService(profileRepo i_repository.IProfileRepository, clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IAdminService {
	return &adminService{
		profileRepo: profileRepo,
		redisCache:  cache.InitializeRedisCache(),
		clients:     clients,
		errLogger:   errLogger,
	}
}

func GenerateAdminService() (business.IAdminService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	//return InitializeAdminService(cnn, errLogger), nil

	return initializeAdminService(
		repository.InitializeProfileRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

const (
	admin_records_limit int = 10
)

// GetAdmin implements business.IAdminService.
func (a *adminService) GetAdmin(id string, ctx context.Context) (response.AdminNftResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.AdminNftResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	adminNft, err := on_chain.GetOnChainObject[entities.AdminNft](on_chain.GetOnChainObjectRequest{
		Client:    a.clients[constant.SuiTestnet],
		ObjectId:  id,
		ErrLogger: a.errLogger,
	}, ctx)

	return adminNft.ToAdminNftResponse(), err
}

// GetAdminByOwner implements business.IAdminService.
func (a *adminService) GetAdminByOwner(id string, ctx context.Context) (response.AdminNftResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.AdminNftResponse{}, genericErr
	}

	var client = a.clients[constant.SuiTestnet]
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: a.errLogger,
	}, ctx)
	if err != nil {
		return response.AdminNftResponse{}, err
	}

	if manage == nil {
		return response.AdminNftResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var adminNftId string
	for i, admin := range manage.AdminIds {
		if admin == id {
			adminNftId = manage.AdminNfts[i]
			break
		}
	}

	if adminNftId == "" {
		return response.AdminNftResponse{}, genericErr
	}

	adminNft, err := on_chain.GetOnChainObject[entities.AdminNft](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  adminNftId,
		ErrLogger: a.errLogger,
	}, ctx)

	return adminNft.ToAdminNftResponse(), err
}

// GetAdmins implements business.IAdminService.
func (a *adminService) GetAdmins(req request.GetAdminsRequest, ctx context.Context) (response.PaginationDataResponse, error) {
	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PaginationDataResponse

	var client = a.clients[constant.SuiTestnet]
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: a.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	admins, err := on_chain.GetOnChainObjects[entities.AdminNft](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: manageObj.AdminNfts,
		ErrLogger: a.errLogger,
	}, ctx)
	if err != nil {
		return response.PaginationDataResponse{}, err
	}

	if len(admins) == 0 {
		return response.PaginationDataResponse{
			Page:   req.Page,
			Amount: 0,
		}, nil
	}

	var filteredAdmins []entities.AdminNft
	for i := len(admins) - 1; i >= 0; i-- {
		var admin entities.AdminNft = admins[i]

		if req.Keyword != "" {
			var firstName string = util.StandardizeString(admin.FirstName)
			var lastName string = util.StandardizeString(admin.LastName)
			if !strings.Contains(firstName, req.Keyword) && !strings.Contains(lastName, req.Keyword) && !strings.Contains(admin.IdentityCode, req.Keyword) && !strings.Contains(admin.PhoneNumber, req.Keyword) && !strings.Contains(admin.Email, req.Keyword) { // Not matched
				continue
			}
		}

		if req.Gender != "" {
			if admin.Gender != req.Gender { // Not matched
				continue
			}
		}

		if req.YearOfBirth != nil {
			var dob time.Time = util.RawDateToTime(admin.DateOfBirth)
			if dob.Year() != *req.YearOfBirth { // Not matched
				continue
			}
		}

		filteredAdmins = append(filteredAdmins, admin)
	}

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredAdmins) <= skippedRecords {
		return response.PaginationDataResponse{}, nil
	}

	sort.Slice(filteredAdmins, func(i, j int) bool {
		if req.SortCriteria == "date_of_birth" {
			var dob1 time.Time = util.RawDateToTime(filteredAdmins[i].DateOfBirth)
			var dob2 time.Time = util.RawDateToTime(filteredAdmins[j].DateOfBirth)
			if req.SortOrder == "DESC" {
				return dob2.After(dob1)
			}

			return dob2.Before(dob1)
		}

		if req.SortOrder == "ASC" {
			return false
		}

		return true
	})

	var data []response.AdminNftResponse
	for i := skippedRecords; i < len(filteredAdmins); i++ {
		data = append(data, filteredAdmins[i].ToAdminNftResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	res = response.PaginationDataResponse{
		Data:       data,
		Amount:     len(data),
		Page:       req.Page,
		TotalPages: int(math.Ceil(float64(len(filteredAdmins)) / float64(req.PageSize))),
	}

	return res, nil
}

// UpdatePublisherInfo implements business.IAdminService.
func (a *adminService) UpdatePublisherInfo(req request.UpdatePublisherInfoRequest, ctx context.Context) (response.BuildTransactionResponse, error) {
	var profileId string = ctx.Value("sub").(string)
	var profile *entities.Profile
	var errRes error
	var isFoundAdmin bool = false
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	for i := 1; i <= 3; i++ {
		profile, errRes = a.profileRepo.GetProfileOfFirsts(i, ctx)
		if errRes != nil {
			return response.BuildTransactionResponse{}, errRes
		}

		if profileId == profile.ID {
			if profile.IdentityCode != nil {
				return response.BuildTransactionResponse{}, genericErr
			}

			isFoundAdmin = true
			break
		}
	}

	if !isFoundAdmin {
		return response.BuildTransactionResponse{}, genericErr
	}

	var dateOfBirth string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(dateOfBirth); dob.IsZero() {
		return response.BuildTransactionResponse{}, genericErr
	}

	var gender string = util.StandardizeGender(util.StandardizeString(req.Gender))
	if gender == "" {
		return response.BuildTransactionResponse{}, genericErr
	}

	var email string = strings.TrimSpace(req.Email)
	if !util.IsValidEmail(email) {
		return response.BuildTransactionResponse{}, genericErr
	}

	var client = a.clients[constant.SuiTestnet]
	var sender string = ctx.Value("address").(string)
	var manangeModule = on_chain.InitializeModuleManage()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.AdminNft](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manangeModule.GetModule(), manangeModule.GetAdminNftStruct()),
		ErrLogger:    a.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	editProfileCaps, err := on_chain.GetOnChainOwnedObjects[entities.Cap](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: sender,
		StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), manangeModule.GetModule(), manangeModule.GetUpdateAdminInfoCapStruct()),
		ErrLogger:    a.errLogger,
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	var identityCode string = strings.TrimSpace(req.IdentityCode)
	var firstName string = strings.TrimSpace(req.FirstName)
	var lastName string = strings.TrimSpace(req.LastName)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)

	// todo: validate identity code
	var module = on_chain.InitializeModuleManage()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    sender,
		Module:    module.GetModule(),
		Function:  module.GetFunctionUpdatePublisherNft(),
		ErrLogger: a.errLogger,
		Arguments: module.ToUpdatePublisherNftArguments(on_chain.UpdatePublisherNftArguments{
			AdminCap:           editProfileCaps[0].ID.ID,
			AdminNft:           nfts[0].ID.ID,
			IdentityCode:       identityCode,
			IdentityCardBlobID: strings.TrimSpace(req.IdentityCardBlobID),
			AvatarBlobID:       strings.TrimSpace(req.AvatarBlobID),
			FirstName:          firstName,
			LastName:           lastName,
			Gender:             gender,
			DateOfBirth:        dateOfBirth,
			PhoneNumber:        phoneNumber,
			Email:              email,
		}),
	}, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	profile.IdentityCode = &identityCode
	profile.FirstName = &firstName
	profile.LastName = &lastName
	profile.Gender = &gender
	profile.DateOfBirth = &dateOfBirth
	profile.PhoneNumber = &phoneNumber
	profile.Email = &email
	profile.UpdatedAt = time.Now()

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, a.profileRepo.UploadProfile(*profile, ctx)
}

func (a *adminService) getGetAdminsRedisKey(req request.GetAdminsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var gender string = "empty"
	if req.Gender != "" {
		gender = req.Gender
	}

	var yob string = "empty"
	if req.YearOfBirth != nil {
		yob = fmt.Sprintf("%d", *req.YearOfBirth)
	}

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("admin:kw:%s:g:%s:y:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, gender, yob, sortCriteria, req.SortOrder, req.PageSize, req.Page,
	)
}

var mockAdmins = []response.AdminNftResponse{
	{
		ID: "1", IdentityCode: "001092000123", IdentityCardBlobID: "ic-101", AvatarBlobID: "avt-101",
		FirstName: "An", LastName: "Nguyễn", Gender: "Female", DateOfBirth: time.Date(1992, 5, 15, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0901234567", Email: "an.nguyen@example.com", UploadedAt: time.Now(),
		Name: "Golden Dragon", Url: "https://nft.com",
	},
	{
		ID: "2", IdentityCode: "001092000456", IdentityCardBlobID: "ic-102", AvatarBlobID: "avt-102",
		FirstName: "Bình", LastName: "Trần", Gender: "Male", DateOfBirth: time.Date(1988, 11, 20, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0912345678", Email: "binh.tran@example.com", UploadedAt: time.Now(),
		Name: "Cyber Punk", Url: "https://nft.com",
	},
	{
		ID: "3", IdentityCode: "001092000789", IdentityCardBlobID: "ic-103", AvatarBlobID: "avt-103",
		FirstName: "Chi", LastName: "Lê", Gender: "Female", DateOfBirth: time.Date(1995, 2, 10, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0923456789", Email: "chi.le@example.com", UploadedAt: time.Now(),
		Name: "Ethereal", Url: "https://nft.com",
	},
	{
		ID: "4", IdentityCode: "001092000111", IdentityCardBlobID: "ic-104", AvatarBlobID: "avt-104",
		FirstName: "Dũng", LastName: "Phạm", Gender: "Male", DateOfBirth: time.Date(1990, 8, 30, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0934567890", Email: "dung.pham@example.com", UploadedAt: time.Now(),
		Name: "Mecha War", Url: "https://nft.com",
	},
	{
		ID: "5", IdentityCode: "001092000222", IdentityCardBlobID: "ic-105", AvatarBlobID: "avt-105",
		FirstName: "Giang", LastName: "Hoàng", Gender: "Female", DateOfBirth: time.Date(1993, 12, 12, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0945678901", Email: "giang.hoang@example.com", UploadedAt: time.Now(),
		Name: "Pixel Cat", Url: "https://nft.com",
	},
	{
		ID: "6", IdentityCode: "001092000333", IdentityCardBlobID: "ic-106", AvatarBlobID: "avt-106",
		FirstName: "Hải", LastName: "Vũ", Gender: "Male", DateOfBirth: time.Date(1985, 3, 25, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0956789012", Email: "hai.vu@example.com", UploadedAt: time.Now(),
		Name: "Ancient", Url: "https://nft.com",
	},
	{
		ID: "7", IdentityCode: "001092000444", IdentityCardBlobID: "ic-107", AvatarBlobID: "avt-107",
		FirstName: "Khánh", LastName: "Phan", Gender: "Male", DateOfBirth: time.Date(1998, 7, 19, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0967890123", Email: "khanh.phan@example.com", UploadedAt: time.Now(),
		Name: "Neon City", Url: "https://nft.com",
	},
	{
		ID: "8", IdentityCode: "001092000555", IdentityCardBlobID: "ic-108", AvatarBlobID: "avt-108",
		FirstName: "Lan", LastName: "Đặng", Gender: "Female", DateOfBirth: time.Date(1994, 9, 05, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0978901234", Email: "lan.dang@example.com", UploadedAt: time.Now(),
		Name: "Lotus Flower", Url: "https://nft.com",
	},
	{
		ID: "9", IdentityCode: "001092000666", IdentityCardBlobID: "ic-109", AvatarBlobID: "avt-109",
		FirstName: "Minh", LastName: "Bùi", Gender: "Male", DateOfBirth: time.Date(1991, 1, 14, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0989012345", Email: "minh.bui@example.com", UploadedAt: time.Now(),
		Name: "Ocean Wave", Url: "https://nft.com",
	},
	{
		ID: "10", IdentityCode: "001092000777", IdentityCardBlobID: "ic-110", AvatarBlobID: "avt-110",
		FirstName: "Nga", LastName: "Đỗ", Gender: "Female", DateOfBirth: time.Date(1996, 6, 21, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0990123456", Email: "nga.do@example.com", UploadedAt: time.Now(),
		Name: "Starry Night", Url: "https://nft.com0",
	},
	{
		ID: "11", IdentityCode: "001092000888", IdentityCardBlobID: "ic-111", AvatarBlobID: "avt-111",
		FirstName: "Phong", LastName: "Hồ", Gender: "Male", DateOfBirth: time.Date(1989, 4, 18, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0321234567", Email: "phong.ho@example.com", UploadedAt: time.Now(),
		Name: "Mountain", Url: "https://nft.com1",
	},
	{
		ID: "12", IdentityCode: "001092000999", IdentityCardBlobID: "ic-112", AvatarBlobID: "avt-112",
		FirstName: "Quỳnh", LastName: "Ngô", Gender: "Female", DateOfBirth: time.Date(1997, 10, 31, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0332345678", Email: "quynh.ngo@example.com", UploadedAt: time.Now(),
		Name: "Violet Aura", Url: "https://nft.com2",
	},
	{
		ID: "13", IdentityCode: "001092001122", IdentityCardBlobID: "ic-113", AvatarBlobID: "avt-113",
		FirstName: "Sơn", LastName: "Dương", Gender: "Male", DateOfBirth: time.Date(1990, 2, 14, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0343456789", Email: "son.duong@example.com", UploadedAt: time.Now(),
		Name: "Abstract", Url: "https://nft.com3",
	},
	{
		ID: "14", IdentityCode: "001092001133", IdentityCardBlobID: "ic-114", AvatarBlobID: "avt-114",
		FirstName: "Thảo", LastName: "Lý", Gender: "Female", DateOfBirth: time.Date(1993, 3, 03, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0354567890", Email: "thao.ly@example.com", UploadedAt: time.Now(),
		Name: "Green Forest", Url: "https://nft.com4",
	},
	{
		ID: "15", IdentityCode: "001092001144", IdentityCardBlobID: "ic-115", AvatarBlobID: "avt-115",
		FirstName: "Tuấn", LastName: "Lương", Gender: "Male", DateOfBirth: time.Date(1987, 12, 25, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0365678901", Email: "tuan.luong@example.com", UploadedAt: time.Now(),
		Name: "Iron Knight", Url: "https://nft.com5",
	},
	{
		ID: "16", IdentityCode: "001092001155", IdentityCardBlobID: "ic-116", AvatarBlobID: "avt-116",
		FirstName: "Vân", LastName: "Vương", Gender: "Female", DateOfBirth: time.Date(1999, 1, 01, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0376789012", Email: "van.vuong@example.com", UploadedAt: time.Now(),
		Name: "Cloud Nine", Url: "https://nft.com6",
	},
	{
		ID: "17", IdentityCode: "001092001166", IdentityCardBlobID: "ic-117", AvatarBlobID: "avt-117",
		FirstName: "Việt", LastName: "Tạ", Gender: "Male", DateOfBirth: time.Date(1992, 11, 11, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0387890123", Email: "viet.ta@example.com", UploadedAt: time.Now(),
		Name: "Red Phoenix", Url: "https://nft.com7",
	},
	{
		ID: "18", IdentityCode: "001092001177", IdentityCardBlobID: "ic-118", AvatarBlobID: "avt-118",
		FirstName: "Xuân", LastName: "Trịnh", Gender: "Female", DateOfBirth: time.Date(1991, 4, 30, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0398901234", Email: "xuan.trinh@example.com", UploadedAt: time.Now(),
		Name: "Spring", Url: "https://nft.com8",
	},
	{
		ID: "19", IdentityCode: "001092001188", IdentityCardBlobID: "ic-119", AvatarBlobID: "avt-119",
		FirstName: "Yên", LastName: "Mạc", Gender: "Female", DateOfBirth: time.Date(1995, 8, 15, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0869012345", Email: "yen.mac@example.com", UploadedAt: time.Now(),
		Name: "Silent Night", Url: "https://nft.com9",
	},
	{
		ID: "20", IdentityCode: "001092001199", IdentityCardBlobID: "ic-120", AvatarBlobID: "avt-120",
		FirstName: "Hùng", LastName: "Cao", Gender: "Male", DateOfBirth: time.Date(1986, 10, 10, 0, 0, 0, 0, time.UTC),
		PhoneNumber: "0880123456", Email: "hung.cao@example.com", UploadedAt: time.Now(),
		Name: "Deep Space", Url: "https://nft.com0",
	},
}
