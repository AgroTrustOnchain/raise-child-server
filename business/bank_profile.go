package business

import (
	"context"
	"errors"
	"log"
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
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type bankProfileService struct {
	bankProfileRepo i_repository.IBankProfileRepository
	redisCache      cache.IRedisCache
	clients         map[string]sui.ISuiAPI
	errLogger       *log.Logger
}

func initializeBankProfileService(bankProfileRepo i_repository.IBankProfileRepository, clients map[string]sui.ISuiAPI, errLogger *log.Logger) business.IBankProfileService {
	return &bankProfileService{
		bankProfileRepo: bankProfileRepo,
		redisCache:      cache.InitializeRedisCache(),
		clients:         clients,
		errLogger:       errLogger,
	}
}

func GenerateBankProfileService() (business.IBankProfileService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeBankProfileService(repository.InitializeBankProfileRepository(cnn, errLogger), _networkAliases, errLogger), nil
}

// CreateBankProfile implements business.IBankProfileService.
func (b *bankProfileService) CreateBankProfile(req request.CreateBankProfileRequest, ctx context.Context) (*entities.BankProfile, error) {
	manage, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:   b.clients[constant.SuiTestnet],
		ObjectId: os.Getenv(env.MANAGE_OBJECT_ID),
	}, ctx)
	if err != nil {
		return nil, err
	}

	var rightAccessErr error = errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	var sender string = ctx.Value("address").(string)
	if !slices.Contains(manage.LocalLeaderIds, sender) {
		return nil, rightAccessErr
	}

	bp, err := b.bankProfileRepo.GetBankProfileByOwner(sender, ctx)
	if err != nil {
		return nil, err
	}

	if bp != nil {
		return nil, errors.New(noti.BANK_PROFILE_EXISTED_MESSAGE)
	}

	var payosClientId string = strings.TrimSpace(req.PayosClientID)
	if payosClientId != "" {
		payosClientId = util.Encrypt(payosClientId)
	}

	var payosApiKey string = strings.TrimSpace(req.PayosApiKey)
	if payosApiKey != "" {
		payosApiKey = util.Encrypt(payosApiKey)
	}

	var payosCheckSumKey string = strings.TrimSpace(req.PayosCheckSumKey)
	if payosCheckSumKey != "" {
		payosCheckSumKey = util.Encrypt(payosCheckSumKey)
	}

	var curTime time.Time = time.Now()
	var bankProfile = entities.BankProfile{
		ID:               util.GenerateId(),
		ProfileID:        ctx.Value("sub").(string),
		Owner:            sender,
		BankOrg:          strings.TrimSpace(req.BankOrg),
		BankCode:         strings.TrimSpace(req.BankCode),
		OwnerName:        strings.TrimSpace(req.OwnerName),
		PayosClientID:    payosClientId,
		PayosApiKey:      payosApiKey,
		PayosCheckSumKey: payosCheckSumKey,
		CreatedAt:        curTime,
		UpdatedAt:        curTime,
	}

	return &bankProfile, b.bankProfileRepo.CreateBankProfile(bankProfile, ctx)
}

// GetBankProfile implements business.IBankProfileService.
func (b *bankProfileService) GetBankProfile(id string, ctx context.Context) (response.BankProfileResponse, error) {
	res, err := b.bankProfileRepo.GetBankProfileById(id, ctx)
	if err != nil {
		return response.BankProfileResponse{}, err
	}

	if res == nil {
		return response.BankProfileResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	var address string = ctx.Value("address").(string)
	if res.ProfileID != ctx.Value("sub").(string) || res.Owner != address {
		manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
			Client:    b.clients[constant.SuiTestnet],
			ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
			ErrLogger: b.errLogger,
		}, ctx)
		if err != nil {
			return response.BankProfileResponse{}, err
		}

		if manageObj == nil {
			return response.BankProfileResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
		}

		if !slices.Contains(manageObj.AdminIds, address) {
			return response.BankProfileResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		}
	}

	return res.ToBankProfileResponse(), err
}

// GetBankProfileByOwner implements business.IBankProfileService.
func (b *bankProfileService) GetBankProfileByOwner(id string, ctx context.Context) (response.BankProfileResponse, error) {
	if !util.IsValidSuiAddressStrict(id) {
		return response.BankProfileResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	res, err := b.bankProfileRepo.GetBankProfileByOwner(id, ctx)
	if err != nil {
		return response.BankProfileResponse{}, err
	}

	if res == nil {
		return response.BankProfileResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return res.ToBankProfileResponse(), err
}

// UpdateBankProfile implements business.IBankProfileService.
func (b *bankProfileService) UpdateBankProfile(id string, req request.UpdateBankProfileRequest, ctx context.Context) (*entities.BankProfile, error) {
	bp, err := b.bankProfileRepo.GetBankProfileById(id, ctx)
	if err != nil {
		return nil, err
	}

	var sender string = ctx.Value("address").(string)
	if bp.Owner != sender {
		return nil, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	var bankOrg string = strings.TrimSpace(req.BankOrg)
	if bankOrg != "" {
		bp.BankOrg = bankOrg
	}

	var bankCode string = strings.TrimSpace(req.BankCode)
	if bankCode != "" {
		bp.BankCode = bankCode
	}

	var ownerName string = strings.TrimSpace(req.OwnerName)
	if ownerName != "" {
		bp.OwnerName = ownerName
	}

	var payosClientId string = strings.TrimSpace(req.PayosClientID)
	if payosClientId != "" {
		payosClientId = util.Encrypt(payosClientId)
	}

	var payosApiKey string = strings.TrimSpace(req.PayosApiKey)
	if payosApiKey != "" {
		payosApiKey = util.Encrypt(payosApiKey)
	}

	var payosCheckSumKey string = strings.TrimSpace(req.PayosCheckSumKey)
	if payosCheckSumKey != "" {
		payosCheckSumKey = util.Encrypt(payosCheckSumKey)
	}

	bp.PayosClientID = payosClientId
	bp.PayosApiKey = payosApiKey
	bp.PayosCheckSumKey = payosCheckSumKey
	bp.UpdatedAt = time.Now()

	return bp, b.bankProfileRepo.UpdateBankProfile(*bp, ctx)
}
