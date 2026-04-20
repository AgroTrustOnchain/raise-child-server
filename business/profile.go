package business

import (
	"context"
	"database/sql"
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
	"strconv"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type profileService struct {
	profileRepo i_repository.IProfileRepository
	redisCache  cache.IRedisCache
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func InitializeProfileService(db *sql.DB, errLogger *log.Logger) business.IProfileService {
	return &profileService{
		profileRepo: repository.InitializeProfileRepository(db, errLogger),
		redisCache:  cache.InitializeRedisCache(),
		clients:     _networkAliases,
		errLogger:   errLogger,
	}
}

func initializeProfileService(
	profileRepo i_repository.IProfileRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IProfileService {
	return &profileService{
		profileRepo: profileRepo,
		redisCache:  cache.InitializeRedisCache(),
		clients:     clients,
		errLogger:   errLogger,
	}
}

func GenerateProfileService() (business.IProfileService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeProfileService(repository.InitializeProfileRepository(cnn, errLogger), _networkAliases, errLogger), nil
}

// UploadProfile implements business.IProfileService.
func (p *profileService) UploadProfile(id string, req request.UploadProfileRequest, ctx context.Context) (response.PersonalProfileResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sub string = ctx.Value("sub").(string)
	if sub != id { // other edits his/her profile
		return response.PersonalProfileResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	profile, err := p.profileRepo.GetProfile(id, ctx)
	if err != nil {
		return response.PersonalProfileResponse{}, err
	}

	if profile == nil {
		return response.PersonalProfileResponse{}, genericErr
	}

	// Already upload profile
	if profile.IdentityCode != nil {
		return response.PersonalProfileResponse{}, genericErr
	}

	var gender string = util.StandardizeGender(util.StandardizeString(req.Gender))
	if gender == "" {
		return response.PersonalProfileResponse{}, errors.New(noti.UNDEFINED_GENDER_MESSAGE)
	}

	var dateOfBirth string = strings.TrimSpace(req.DateOfBirth)
	if dob := util.RawDateToTime(dateOfBirth); dob.IsZero() {
		return response.PersonalProfileResponse{}, errors.New(noti.INVALID_DATE_FORMAT_WARN_MSG)
	}

	var identityCode string = strings.TrimSpace(req.IdentityCode)
	var phoneNumber string = strings.TrimSpace(req.PhoneNumber)
	var email string = strings.TrimSpace(req.Email)
	if !util.IsValidEmail(email) {
		return response.PersonalProfileResponse{}, genericErr
	}

	isInfoExist, err := p.profileRepo.IsPersonalInfoExist(identityCode, phoneNumber, email, ctx)
	if err != nil {
		return response.PersonalProfileResponse{}, err
	}

	if isInfoExist {
		return response.PersonalProfileResponse{}, errors.New(noti.GENERIC_PERSONAL_INFO_REGISTERED_WANR_MSG)
	}

	var firstName string = strings.TrimSpace(req.FirstName)
	var lastName string = strings.TrimSpace(req.LastName)

	profile.IdentityCode = &identityCode
	profile.FirstName = &firstName
	profile.LastName = &lastName
	profile.Gender = &gender
	profile.DateOfBirth = &dateOfBirth
	profile.PhoneNumber = &phoneNumber
	profile.Email = &email
	profile.UpdatedAt = time.Now()

	return (*profile).ToPersonalProfile(), p.profileRepo.UploadProfile(*profile, ctx)
}

// GetWalletPersonalProfile implements business.IProfileService.
func (p *profileService) GetWalletPersonalProfile(id string, req request.GetTransactionRecordsRequest, ctx context.Context) (response.PersonalWalletProfileResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(id) {
		return response.PersonalWalletProfileResponse{}, genericErr
	}

	if req.PoolID != "" {
		if !util.IsValidSuiAddressStrict(req.PoolID) {
			return response.PersonalWalletProfileResponse{}, genericErr
		}
	}

	if req.MaxAmount != nil {
		if *req.MaxAmount < min_withdraw_proposal_amount_value {
			return response.PersonalWalletProfileResponse{}, nil
		}

		if req.MinAmount != nil {
			if *req.MaxAmount <= *req.MinAmount {
				return response.PersonalWalletProfileResponse{}, nil
			}
		}
	}

	req.SortOrder = util.StandardizeSortOrder(req.SortOrder)
	req.SortCriteria = util.StandardizeSortCriteria(req.SortCriteria)
	req.Keyword = util.StandardizeString(req.Keyword)
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = default_page_size
	}

	var res response.PersonalWalletProfileResponse
	var redisKey string = p.getGetWalletPersonalProfileRedisKey(id, req)
	if p.redisCache.Get(redisKey, &res, ctx) {
		return res, nil
	}

	var client = p.clients[constant.SuiTestnet]
	var packageId string = os.Getenv(env.PACKAGE_ID)
	var donorModule = on_chain.InitializeModuleDonor()
	nfts, err := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: id,
		StructType:   fmt.Sprintf("%s::%s::%s", packageId, donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		ErrLogger:    p.errLogger,
	}, ctx)
	if err != nil {
		return res, err
	}

	if nfts == nil || len(nfts) == 0 {
		return res, nil
	}

	var recordModule = on_chain.InitializeModuleRecord()
	txs, err := on_chain.GetOnChainOwnedObjects[entities.Transaction](on_chain.GetOnChainOwnedObjectsRequest{
		Client:       client,
		OwnerAddress: id,
		StructType:   fmt.Sprintf("%s::%s::%s", packageId, recordModule.GetModule(), recordModule.GetTransactionRecordStruct()),
		ErrLogger:    p.errLogger,
	}, ctx)
	if err != nil {
		return res, err
	}

	var poolName string
	if req.PoolID != "" {
		if req.PoolID == os.Getenv(env.POOL_ID) {
			poolName = "Main Pool"
		} else {
			pool, _ := on_chain.GetOnChainObject[entities.LocalPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  req.PoolID,
				ErrLogger: p.errLogger,
			}, ctx)

			if pool != nil {
				poolName = pool.Region
			}
		}
	}

	var filteredTxs []entities.Transaction
	for _, tx := range txs {
		if req.Keyword != "" {
			if !strings.Contains(tx.Message, req.Keyword) {
				continue
			}
		}

		if poolName != "" {
			if tx.PoolName != poolName {
				continue
			}
		}

		if req.ActionType != "" {
			if tx.ActionType != req.ActionType {
				continue
			}
		}

		amount, _ := strconv.ParseInt(tx.Amount, 10, 64)
		if req.MinAmount != nil {
			if amount < *req.MinAmount {
				continue
			}
		}

		if req.MaxAmount != nil {
			if amount > *req.MaxAmount {
				continue
			}
		}

		filteredTxs = append(filteredTxs, tx)
	}

	sort.Slice(filteredTxs, func(i, j int) bool {
		if req.SortCriteria == "amount" {
			amount1, _ := strconv.ParseInt(filteredTxs[i].Amount, 10, 64)
			amount2, _ := strconv.ParseInt(filteredTxs[j].Amount, 10, 64)
			if req.SortOrder == "DESC" {
				return amount2 > amount1
			}

			return amount2 < amount1
		}

		if req.SortOrder == "ASC" {
			return true
		}

		return false
	})

	var skippedRecords int = (req.Page - 1) * req.PageSize
	if len(filteredTxs) <= skippedRecords {
		return res, nil
	}

	var data []response.TransactionResponse
	for i := skippedRecords; i < len(filteredTxs); i++ {
		data = append(data, filteredTxs[i].ToTransactionResponse())
		if len(data) == req.PageSize {
			break
		}
	}

	totalDonation, _ := strconv.ParseInt(nfts[0].TotalDonation, 10, 64)
	res = response.PersonalWalletProfileResponse{
		WalletAddress:   id,
		FirstName:       nfts[0].FirstName,
		LastName:        nfts[0].LastName,
		TotalDonation:   totalDonation,
		SupportedChilds: nfts[0].SupportedChilds,
		TxRecords:       data,
		RecordAmount:    len(data),
		Page:            req.Page,
		TotalPages:      int(math.Ceil(float64(len(filteredTxs)) / float64(req.PageSize))),
	}

	p.redisCache.Set(redisKey, res, time.Minute*5, ctx)

	return res, nil
}

func (p *profileService) getGetWalletPersonalProfileRedisKey(wallet string, req request.GetTransactionRecordsRequest) string {
	var keyword string = "empty"
	if req.Keyword != "" {
		keyword = req.Keyword
	}

	var poolId string = "empty"
	if req.PoolID != "" {
		poolId = req.PoolID
	}

	var actionType string = "empty"
	if req.ActionType != "" {
		actionType = req.ActionType
	}

	var minAmount string = "empty"
	if req.MinAmount != nil {
		minAmount = fmt.Sprintf("%d", *req.MinAmount)
	}

	var maxAmount string = "empty"
	if req.MaxAmount != nil {
		maxAmount = fmt.Sprintf("%d", *req.MaxAmount)
	}

	var sortCriteria string = "empty"
	if req.SortCriteria != "" {
		sortCriteria = req.SortCriteria
	}

	return fmt.Sprintf("wallet_personal_profile:kw:%s:pool:%s:of:%s:type:%s:min:%s:max:%s:sc:%s:o:%s:s:%d:p:%d",
		keyword, poolId, wallet, actionType, minAmount, maxAmount, sortCriteria, req.SortOrder, req.PageSize, req.Page)
}
