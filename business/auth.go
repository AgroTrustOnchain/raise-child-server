package business

import (
	"context"
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repo "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/cache"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"raise-child/util/security"
	"slices"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type authService struct {
	profileRepo i_repo.IProfileRepository
	redisCache  cache.IRedisCache
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func InitializeAuthService(db *sql.DB, errLogger *log.Logger, redisCache cache.IRedisCache) business.IAuthService {
	return &authService{
		profileRepo: repository.InitializeProfileRepository(db, errLogger),
		redisCache:  redisCache,
		clients:     _networkAliases,
		errLogger:   errLogger,
	}
}

func GenerateAuthService() (business.IAuthService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializeAuthService(cnn, errLogger, cache.InitializeRedisCache()), nil
}

// LoginV2 implements business.IAuthService.
func (a *authService) LoginV2(req request.LoginRequestV2, ctx context.Context) (response.LoginResponse, error) {
	var address string = strings.TrimSpace(req.Address)
	var client = a.clients[constant.SuiTestnet]

	var sub string = strings.TrimSpace(req.Sub)
	manageObj, err := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: a.errLogger,
	}, ctx)
	if err != nil {
		return response.LoginResponse{}, err
	}

	var roles []string
	if manageObj.ID.ID != "" {
		if slices.Contains(manageObj.AdminIds, address) {
			roles = append(roles, admin_role)
		}

		if slices.Contains(manageObj.LocalLeaderIds, address) {
			roles = append(roles, local_leader_role)
		}

		if slices.Contains(manageObj.VolunteerIds, address) {
			roles = append(roles, volunteer_role)
		}

		if slices.Contains(manageObj.DonorIds, address) {
			roles = append(roles, donor_role)
		} else {
			roles = append(roles, user_role)
		}
	}

	token, _, err := security.GenerateActionTokenV2(address, sub, roles, a.errLogger)
	if err != nil {
		return response.LoginResponse{}, err
	}

	profile, err := a.profileRepo.GetProfile(sub, ctx)
	if err != nil {
		return response.LoginResponse{}, err
	}

	if profile == nil {
		return response.LoginResponse{}, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	if profile.WalletAddress != nil {
		if *profile.WalletAddress != address {
			return response.LoginResponse{}, errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
		}
	}

	a.profileRepo.Login(sub, token, address, ctx)

	setLogin(sub, address)

	return response.LoginResponse{
		Token: token,
	}, nil
}

// LogoutV2 implements business.IAuthService.
func (a *authService) LogoutV2(ctx context.Context) error {
	var sub string = ctx.Value("sub").(string)
	if err := a.profileRepo.Logout(sub, ctx); err != nil {
		return err
	}

	logoutWallet(sub)
	return nil
}

// GetSalt implements business.IAuthService.
func (a *authService) GetSalt(id string, ctx context.Context) (response.GetSaltResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if id == "" {
		return response.GetSaltResponse{}, genericErr
	}

	profile, err := a.profileRepo.GetProfile(id, ctx)
	if err != nil {
		return response.GetSaltResponse{}, err
	}

	if profile == nil {
		var salt string = util.GenerateSalt()
		var curTime = time.Now()
		return response.GetSaltResponse{
				Salt: salt,
			}, a.profileRepo.CreateProfile(entities.Profile{
				ID:        id,
				Salt:      salt,
				CreatedAt: curTime,
				UpdatedAt: curTime,
			}, ctx)
	}

	return response.GetSaltResponse{
		Salt: profile.Salt,
	}, nil
}

// Login implements business.IAuthService.
func (a *authService) Login(req request.LoginRequest, ctx context.Context) (response.LoginResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var internalErr error = errors.New(noti.INTERNALL_ERR_MSG)
	var curTimeUnix int64 = time.Now().Unix()
	var info securityInfo = getSecurityInfo(req.Address)

	// Token still active
	if info.exp >= curTimeUnix {
		return response.LoginResponse{}, genericErr
	}

	var matchedNonce string = info.nonce
	if !strings.Contains(req.Message, matchedNonce) {
		return response.LoginResponse{}, genericErr
	}

	pubKeyBytes, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		a.errLogger.Println(err.Error())
		return response.LoginResponse{}, internalErr
	}

	sigBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		a.errLogger.Println(err.Error())
		return response.LoginResponse{}, internalErr
	}

	var msgBytes = []byte(req.Message)
	if !ed25519.Verify(ed25519.PublicKey(pubKeyBytes), msgBytes, sigBytes) {
		return response.LoginResponse{}, genericErr
	}

	var convertedAddress string = security.PublicKeyToSuiAddress(pubKeyBytes)
	if convertedAddress != req.Address {
		return response.LoginResponse{}, genericErr
	}

	token, exp, err := security.GenerateActionToken(req.Address, matchedNonce, "", a.errLogger)
	if err != nil {
		return response.LoginResponse{}, internalErr
	}
	setLoggedIn(req.Address, exp)

	return response.LoginResponse{Token: token}, nil
}

// Logout implements business.IAuthService.
func (a *authService) Logout(address string, ctx context.Context) error {
	var info securityInfo = getSecurityInfo(address)
	if info.isLoggedIn {
		removeLoggedIn(address)
		return nil
	}

	return errors.New(noti.GENERIC_ERROR_WARN_MSG)
}

// GetNonce implements business.IAuthService.
func (a *authService) GetNonce(address string, ctx context.Context) (response.GetNonceResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if address == "" {
		return response.GetNonceResponse{}, genericErr
	}

	var curTimeUnix int64 = time.Now().Unix()
	var securityInfo securityInfo = getSecurityInfo(address)
	var nonce string = util.GenerateNonce()

	// New login request
	if securityInfo.nonce == "" {
		setAddress(address, nonce, 0)
		return response.GetNonceResponse{Nonce: nonce}, nil
	}

	// Expired nonce, generate new one
	if securityInfo.exp < curTimeUnix {
		setNonce(address, nonce)
		return response.GetNonceResponse{Nonce: nonce}, nil
	}

	// Token still valid
	return response.GetNonceResponse{}, genericErr
}
