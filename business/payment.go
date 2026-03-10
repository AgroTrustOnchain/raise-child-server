package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/env/payment"
	"raise-child/constants/noti"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/model/dtos/response"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/payOSHQ/payos-lib-golang"
)

const (
	payment_pending_status string = "Pending"
	payment_success_status string = "Success"
	payment_cancel_status  string = "Canceled"
)

type paymentService struct {
	paymentRepo  i_repository.IPaymentRepository
	profileRepo  i_repository.IProfileRepository
	donationRepo i_repository.IOffChainDonationRepository
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository
	clients      map[string]sui.ISuiAPI
	errLogger    *log.Logger
}

func InitializePaymentService(db *sql.DB, errLogger *log.Logger) business.IPaymentService {
	return &paymentService{
		paymentRepo:  repository.InitializePaymentRepository(db, errLogger),
		profileRepo:  repository.InitializeProfileRepository(db, errLogger),
		donationRepo: repository.InitializeOffChainDonationRepository(db, errLogger),
		withdrawRepo: repository.InitializeOffChainWithdrawProposalRepository(db, errLogger),
		clients:      _networkAliases,
		errLogger:    errLogger,
	}
}

func GeneratePaymentService() (business.IPaymentService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return InitializePaymentService(cnn, errLogger), nil
}

// ConfirmWithdrawProposal implements business.IPaymentService.
func (p *paymentService) ConfirmWithdrawProposal(id string, ctx context.Context) (map[string]interface{}, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) || !utils.IsValidSuiAddress(models.SuiAddress(id)) {
		return nil, genericErr
	}

	return nil, nil
}

// // CallbackTx implements business.IPaymentService.
// func (p *paymentService) CallbackTx(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
// 	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

// 	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if payment == nil {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	if !payment.IsDonateTx {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	// Expired
// 	if payment.ExpiredAt.Before(time.Now()) {
// 		return response.BuildTransactionResponse{}, errors.New(noti.PAYMENT_EXPIRED_MESSAGE)
// 	}

// 	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
// 	if err != nil {
// 		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
// 		return response.BuildTransactionResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
// 	}

// 	switch data.Status {
// 	case shared.PAYOS_PAID_STATUS:
// 		payment.Status = shared.PAYMENT_SUCCESS_STATUS
// 	case shared.PAYOS_CANCELLED_STATUS:
// 		payment.Status = shared.PAYMENT_CANCELED_STATUS
// 		if data.CancellationReason != nil {
// 			payment.CancelReason = *data.CancellationReason
// 		}
// 	default:
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	payment.UpdatedAt = time.Now()
// 	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if data.Status == shared.PAYOS_CANCELLED_STATUS {
// 		return response.BuildTransactionResponse{}, errors.New(noti.PAYMENT_CANCEL_MESSAGE)
// 	}

// 	profile, err := p.profileRepo.GetProfile(payment.Sub, ctx)
// 	if err != nil {
// 		return response.BuildTransactionResponse{}, err
// 	}

// 	if profile == nil {
// 		return response.BuildTransactionResponse{}, genericErr
// 	}

// 	var module = on_chain.InitializeModulePool()
// 	var function string
// 	var args []interface{}
// 	if payment.Target != os.Getenv(env.POOL_ID) { // Local Pool
// 		function = module.GetFunctionDonateToLocalPool()
// 		args = module.ToDonateToLocalPoolArguments(on_chain.DonateToLocalPoolArguments{
// 			LocalPoolId: payment.Target,
// 			DonateToPoolArguments: on_chain.DonateToPoolArguments{
// 				Amount:      payment.Amount,
// 				FirstName:   profile.FirstName,
// 				LastName:    profile.LastName,
// 				Gender:      profile.Gender,
// 				PhoneNumber: profile.PhoneNumber,
// 				Email:       profile.Email,
// 				Message:     payment.Message,
// 			},
// 		})
// 	} else { // Main pool
// 		function = module.GetFunctionDonateToPool()
// 		args = module.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
// 			Amount:      payment.Amount,
// 			FirstName:   profile.FirstName,
// 			LastName:    profile.LastName,
// 			Gender:      profile.Gender,
// 			PhoneNumber: profile.PhoneNumber,
// 			Email:       profile.Email,
// 			Message:     payment.Message,
// 		})
// 	}

// 	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
// 		Client:    p.clients[constant.SuiTestnet],
// 		Sender:    payment.Actor,
// 		Module:    module.GetModule(),
// 		Function:  function,
// 		ErrLogger: p.errLogger,
// 		Arguments: args,
// 	}, ctx)

// 	return response.BuildTransactionResponse{
// 		TxBytes: txBytes,
// 	}, err
// }

// Donate implements business.IPaymentService.
func (p *paymentService) Donate(req request.DonateRequest, ctx context.Context) (response.UrlAPIResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !util.IsValidSuiAddressStrict(req.PoolId) {
		return response.UrlAPIResponse{}, genericErr
	}
	profile, err := p.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.UrlAPIResponse{}, err
	}

	if profile == nil {
		return response.UrlAPIResponse{}, genericErr
	}

	if profile.IdentityCode == "" {
		return response.UrlAPIResponse{}, errors.New(noti.NOT_UPLOADED_PROFILE_MESSAGE)
	}

	var paymentId string = util.GenerateId()
	var orderCode int = util.GenerateNumber()
	var callbackUrl string = os.Getenv(payment.PAYMENT_CALLBACK_URL) + paymentId
	var description string = req.Message
	if description == "" {
		description = fmt.Sprint(orderCode)
	}

	data, err := payos.CreatePaymentLink(payos.CheckoutRequestType{
		OrderCode:   int64(orderCode),
		Amount:      int(req.Amount),
		Description: description,
		ReturnUrl:   callbackUrl,
		CancelUrl:   callbackUrl,
	})

	if err != nil {
		p.errLogger.Println("Err: ", err.Error())
		return response.UrlAPIResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var donationId string = util.GenerateId()
	var curTime time.Time = time.Now()
	if err := p.donationRepo.CreateDonation(entities.OffChainDonation{
		ID:        donationId,
		Purpose:   string(entities.DONATE_PURPOSE),
		Target:    req.PoolId,
		CreatedAt: curTime,
	}, ctx); err != nil {
		return response.UrlAPIResponse{}, err
	}

	return response.UrlAPIResponse{
			Url: data.CheckoutUrl,
		}, p.paymentRepo.CreatePayment(entities.Payment{
			ID:            paymentId,
			Actor:         ctx.Value("address").(string),
			ProfileID:     profile.ID,
			DonationID:    &donationId,
			IsDonateTx:    true,
			TransactionId: fmt.Sprint(orderCode),
			Amount:        req.Amount,
			Currency:      shared.VIETNAMDONG_CURRENCY,
			Status:        payment_pending_status,
			Method:        shared.PAYMENT_PAYOS_METHOD,
			Message:       description,
			ExpiredAt:     time.Unix(int64(*data.ExpiredAt), 0),
			CreatedAt:     curTime,
			UpdatedAt:     curTime,
		}, ctx)
}

// Callback implements business.IPaymentService.
func (p *paymentService) Callback(id string, ctx context.Context) (string, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return "", err
	}

	if payment == nil {
		return "", genericErr
	}

	// Expired
	if payment.ExpiredAt.Before(time.Now()) {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    shared.PAYAMENT_EXPIRED_STATUS,
			Message:   noti.PAYMENT_EXPIRED_MESSAGE,
		}), nil
	}

	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
	if err != nil {
		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	switch data.Status {
	case shared.PAYOS_PAID_STATUS:
		payment.Status = shared.PAYMENT_SUCCESS_STATUS
	case shared.PAYOS_CANCELLED_STATUS:
		payment.Status = shared.PAYMENT_CANCELED_STATUS
		if data.CancellationReason != nil {
			payment.CancelReason = *data.CancellationReason
		}
	default:
		return "", genericErr
	}

	payment.UpdatedAt = time.Now()
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return "", err
	}

	if data.Status == shared.PAYMENT_CANCELED_STATUS {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    data.Status,
			Message:   noti.PAYMENT_CANCEL_MESSAGE,
		}), nil
	}

	profile, err := p.profileRepo.GetProfile(payment.ProfileID, ctx)
	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", genericErr
	}

	var client = p.clients[constant.SuiTestnet]
	var module string
	var function string
	var args []interface{}
	if payment.IsDonateTx {
		detail, err := p.donationRepo.GetDonation(*payment.DonationID, ctx)
		if err != nil {
			return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				OrderCode: payment.TransactionId,
				Status:    data.Status,
				Message:   err.Error(),
			}), nil
		}

		var donorModule = on_chain.InitializeModuleDonor()
		var nftId string
		nfts, err := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: payment.Actor,
			StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		}, ctx)
		if err != nil {
			return "", err
		}

		if nfts != nil && len(nfts) > 0 {
			nftId = nfts[0].ID.ID
		} else {
			nftId = os.Getenv(env.PUBLISHER_NFT_ID)
		}

		if detail.Purpose == string(entities.DONATE_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			module = poolModule.GetModule()
			if detail.Target != os.Getenv(env.POOL_ID) { // Donate to local pool
				function = poolModule.GetFunctionDonateToLocalPool()
				args = poolModule.ToDonateToLocalPoolArguments(on_chain.DonateToLocalPoolArguments{
					LocalPoolId: detail.Target,
					DonateToPoolArguments: on_chain.DonateToPoolArguments{
						DonorID:     nftId,
						Amount:      payment.Amount,
						FirstName:   profile.FirstName,
						LastName:    profile.LastName,
						Gender:      profile.Gender,
						PhoneNumber: profile.PhoneNumber,
						Email:       profile.Email,
						Message:     payment.Message,
					},
				})
			} else {
				function = poolModule.GetFunctionDonateToPool()
				args = poolModule.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
					DonorID:     nftId,
					Amount:      payment.Amount,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
					Email:       profile.Email,
					Message:     payment.Message,
				})
			}
		} else {
			var childModule = on_chain.InitializeModuleChild()
			module = childModule.GetModule()
			var targetId, childId string

			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				targetId = need.ID.ID
				childId = need.ChildID
			case string(entities.MEAL_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				targetId = need.ID.ID
				childId = need.ChildID
			case string(entities.SPECIAL_NEED_PURPOSE):
				campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				targetId = campaign.ID.ID
				childId = campaign.ChildID
			}

			child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  childId,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  os.Getenv(env.POOL_ID),
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.LocalPools,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			var localPoolId string
			for _, localPool := range localPools {
				if localPool.Region == child.Region {
					localPoolId = localPool.ID.ID
					break
				}
			}

			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				function = childModule.GetFunctionSupportChildBooksNeed()
				args = childModule.ToSupportChildBooksNeedArguments(on_chain.SupportChildBooksNeedArguments{
					NeedID:      targetId,
					LocalPool:   localPoolId,
					ChildID:     childId,
					DonorNft:    nftId,
					Amount:      payment.Amount,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
					Email:       profile.Email,
					Message:     payment.Message,
				})
			case string(entities.MEAL_NEED_PURPOSE):
				function = childModule.GetFunctionSupportChildMealNeed()
				args = childModule.ToSupportChildMealNeedArguments(on_chain.SupportChildMealNeedArguments{
					StartPeriod: "",
					EndPeriod:   "",
					SupportChildBooksNeedArguments: on_chain.SupportChildBooksNeedArguments{
						NeedID:      targetId,
						LocalPool:   localPoolId,
						ChildID:     childId,
						DonorNft:    nftId,
						Amount:      payment.Amount,
						FirstName:   profile.FirstName,
						LastName:    profile.LastName,
						Gender:      profile.Gender,
						PhoneNumber: profile.PhoneNumber,
						Email:       profile.Email,
						Message:     payment.Message,
					},
				})
			case string(entities.SPECIAL_NEED_PURPOSE):
				function = childModule.GetFunctionSupportChildSpecialNeedCampaign()
				args = childModule.ToSupportChildSpeicalNeedArguments(on_chain.SupportChildSpeicalNeedArguments{
					CampaignID:  targetId,
					LocalPool:   localPoolId,
					ChildID:     childId,
					DonorNft:    nftId,
					Amount:      payment.Amount,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
					Email:       profile.Email,
					Message:     payment.Message,
				})
			}
		}
	} else {
		detail, err := p.withdrawRepo.GetOffChainWithdrawProposal(*payment.ProposalID, ctx)
		if err != nil {
			return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				OrderCode: payment.TransactionId,
				Status:    data.Status,
				Message:   err.Error(),
			}), nil
		}

		if detail.Purpose == string(entities.WITHDRAW_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			var localPoolId string
			if detail.Target != os.Getenv(env.POOL_ID) {
				localPoolId = detail.Target
			} else {
				localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
			}

			module = poolModule.GetModule()
			function = poolModule.GetFunctionWithdrawFromPool()
			args = poolModule.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
				LocalPoolId:        localPoolId,
				WithdrawProposalId: detail.ProposalID,
			})
		} else {
			proposal, err := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  detail.ProposalID,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  os.Getenv(env.POOL_ID),
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.LocalPools,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			var localPoolId string
			for _, localPool := range localPools {
				if localPool.Region == proposal.PoolName {
					localPoolId = localPool.ID.ID
					break
				}
			}

			var childModule = on_chain.InitializeModuleChild()
			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				function = childModule.GetFunctionWithdrawFromBooksNeedProposal()
			case string(entities.MEAL_NEED_PURPOSE):
				function = childModule.GetFunctionWithdrawFromMealNeedProposal()
			case string(entities.SPECIAL_NEED_PURPOSE):
				function = childModule.GetFunctionWithdrawFromSpecialNeedCampaign()
			}

			args = childModule.ToWithdrawFromNeedArguments(on_chain.WithdrawFromNeedArguments{
				LocalPool:  localPoolId,
				TargetID:   detail.Target,
				ProposalID: detail.ProposalID,
			})
		}
	}

	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    payment.Actor,
		Module:    module,
		Function:  function,
		ErrLogger: p.errLogger,
		Arguments: args,
	}, ctx)

	var req = util.GenerateRedirectParamRequest{
		OrderCode: payment.TransactionId,
		Status:    data.Status,
		TxBytes:   txBytes,
	}

	if err != nil {
		req.Message = err.Error()
	} else {
		req.Message = "Success"
	}

	return util.GeneratePaymentRedirectUrl(req), nil
}

// CallbackWithAuth implements business.IPaymentService.
func (p *paymentService) CallbackWithAuth(id string, capturedImgBlobId string, ctx context.Context) (string, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return "", err
	}

	if payment == nil {
		return "", genericErr
	}

	var sender string = ctx.Value("address").(string)
	if !utils.IsValidSuiAddress(models.SuiAddress(sender)) || payment.Actor != sender {
		return "", errors.New(noti.GENERIC_RIGHT_ACCESS_WARN_MSG)
	}

	// Expired
	if payment.ExpiredAt.Before(time.Now()) {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    shared.PAYAMENT_EXPIRED_STATUS,
			Message:   noti.PAYMENT_EXPIRED_MESSAGE,
		}), nil
	}

	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
	if err != nil {
		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	switch data.Status {
	case shared.PAYOS_PAID_STATUS:
		payment.Status = shared.PAYMENT_SUCCESS_STATUS
	case shared.PAYOS_CANCELLED_STATUS:
		payment.Status = shared.PAYMENT_CANCELED_STATUS
		if data.CancellationReason != nil {
			payment.CancelReason = *data.CancellationReason
		}
	default:
		return "", genericErr
	}

	payment.UpdatedAt = time.Now()
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return "", err
	}

	if data.Status == shared.PAYMENT_CANCELED_STATUS {
		return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
			OrderCode: payment.TransactionId,
			Status:    data.Status,
			Message:   noti.PAYMENT_CANCEL_MESSAGE,
		}), nil
	}

	profile, err := p.profileRepo.GetProfile(payment.ProfileID, ctx)
	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", genericErr
	}

	var client = p.clients[constant.SuiTestnet]
	var module string
	var function string
	var args []interface{}
	if payment.IsDonateTx {
		detail, err := p.donationRepo.GetDonation(*payment.DonationID, ctx)
		if err != nil {
			return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				OrderCode: payment.TransactionId,
				Status:    data.Status,
				Message:   err.Error(),
			}), nil
		}

		var donorModule = on_chain.InitializeModuleDonor()
		var nftId string
		nfts, err := on_chain.GetOnChainOwnedObjects[entities.Donor](on_chain.GetOnChainOwnedObjectsRequest{
			Client:       client,
			OwnerAddress: payment.Actor,
			StructType:   fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), donorModule.GetModule(), donorModule.GetDonorNftStruct()),
		}, ctx)
		if err != nil {
			return "", err
		}

		if nfts != nil && len(nfts) > 0 {
			nftId = nfts[0].ID.ID
		} else {
			nftId = os.Getenv(env.PUBLISHER_NFT_ID)
		}

		if detail.Purpose == string(entities.DONATE_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			module = poolModule.GetModule()
			if detail.Target != os.Getenv(env.POOL_ID) { // Donate to local pool
				function = poolModule.GetFunctionDonateToLocalPool()
				args = poolModule.ToDonateToLocalPoolArguments(on_chain.DonateToLocalPoolArguments{
					LocalPoolId: detail.Target,
					DonateToPoolArguments: on_chain.DonateToPoolArguments{
						DonorID:     nftId,
						Amount:      payment.Amount,
						FirstName:   profile.FirstName,
						LastName:    profile.LastName,
						Gender:      profile.Gender,
						PhoneNumber: profile.PhoneNumber,
						Email:       profile.Email,
						Message:     payment.Message,
					},
				})
			} else {
				function = poolModule.GetFunctionDonateToPool()
				args = poolModule.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
					DonorID:     nftId,
					Amount:      payment.Amount,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
					Email:       profile.Email,
					Message:     payment.Message,
				})
			}
		} else {
			var childModule = on_chain.InitializeModuleChild()
			module = childModule.GetModule()
			var targetId, childId string

			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				targetId = need.ID.ID
				childId = need.ChildID
			case string(entities.MEAL_NEED_PURPOSE):
				need, err := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				targetId = need.ID.ID
				childId = need.ChildID
			case string(entities.SPECIAL_NEED_PURPOSE):
				campaign, err := on_chain.GetOnChainObject[entities.SpecialNeedCampaign](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  detail.Target,
					ErrLogger: p.errLogger,
				}, ctx)
				if err != nil {
					return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
						OrderCode: payment.TransactionId,
						Status:    data.Status,
						Message:   err.Error(),
					}), nil
				}

				targetId = campaign.ID.ID
				childId = campaign.ChildID
			}

			child, err := on_chain.GetOnChainObject[entities.Child](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  childId,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  os.Getenv(env.POOL_ID),
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.LocalPools,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			var localPoolId string
			for _, localPool := range localPools {
				if localPool.Region == child.Region {
					localPoolId = localPool.ID.ID
					break
				}
			}

			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				args = childModule.ToSupportChildBooksNeedArguments(on_chain.SupportChildBooksNeedArguments{
					NeedID:      targetId,
					LocalPool:   localPoolId,
					ChildID:     childId,
					DonorNft:    nftId,
					Amount:      payment.Amount,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
					Email:       profile.Email,
					Message:     payment.Message,
				})
			case string(entities.MEAL_NEED_PURPOSE):
				args = childModule.ToSupportChildMealNeedArguments(on_chain.SupportChildMealNeedArguments{
					StartPeriod: "",
					EndPeriod:   "",
					SupportChildBooksNeedArguments: on_chain.SupportChildBooksNeedArguments{
						NeedID:      targetId,
						LocalPool:   localPoolId,
						ChildID:     childId,
						DonorNft:    nftId,
						Amount:      payment.Amount,
						FirstName:   profile.FirstName,
						LastName:    profile.LastName,
						Gender:      profile.Gender,
						PhoneNumber: profile.PhoneNumber,
						Email:       profile.Email,
						Message:     payment.Message,
					},
				})
			case string(entities.SPECIAL_NEED_PURPOSE):
				args = childModule.ToSupportChildSpeicalNeedArguments(on_chain.SupportChildSpeicalNeedArguments{
					CampaignID:  targetId,
					LocalPool:   localPoolId,
					ChildID:     childId,
					DonorNft:    nftId,
					Amount:      payment.Amount,
					FirstName:   profile.FirstName,
					LastName:    profile.LastName,
					Gender:      profile.Gender,
					PhoneNumber: profile.PhoneNumber,
					Email:       profile.Email,
					Message:     payment.Message,
				})
			}
		}
	} else {
		detail, err := p.withdrawRepo.GetOffChainWithdrawProposal(*payment.ProposalID, ctx)
		if err != nil {
			return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
				OrderCode: payment.TransactionId,
				Status:    data.Status,
				Message:   err.Error(),
			}), nil
		}

		if detail.Purpose == string(entities.WITHDRAW_PURPOSE) {
			var poolModule = on_chain.InitializeModulePool()
			var localPoolId string
			if detail.Target != os.Getenv(env.POOL_ID) {
				localPoolId = detail.Target
			} else {
				localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
			}

			module = poolModule.GetModule()
			function = poolModule.GetFunctionWithdrawFromPool()
			args = poolModule.ToWithdrawFromPoolArguments(on_chain.WithdrawFromPoolArguments{
				LocalPoolId:        localPoolId,
				WithdrawProposalId: detail.ProposalID,
			})
		} else {
			proposal, err := on_chain.GetOnChainObject[entities.WithdrawProposal](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  detail.ProposalID,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			pool, err := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  os.Getenv(env.POOL_ID),
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			localPools, err := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
				Client:    client,
				ObjectIds: pool.LocalPools,
				ErrLogger: p.errLogger,
			}, ctx)
			if err != nil {
				return util.GeneratePaymentRedirectUrl(util.GenerateRedirectParamRequest{
					OrderCode: payment.TransactionId,
					Status:    data.Status,
					Message:   err.Error(),
				}), nil
			}

			var localPoolId string
			for _, localPool := range localPools {
				if localPool.Region == proposal.PoolName {
					localPoolId = localPool.ID.ID
					break
				}
			}

			var childModule = on_chain.InitializeModuleChild()
			switch detail.Purpose {
			case string(entities.BOOKS_NEED_PURPOSE):
				function = childModule.GetFunctionWithdrawFromBooksNeedProposal()
			case string(entities.MEAL_NEED_PURPOSE):
				function = childModule.GetFunctionWithdrawFromMealNeedProposal()
			case string(entities.SPECIAL_NEED_PURPOSE):
				function = childModule.GetFunctionWithdrawFromSpecialNeedCampaign()
			}

			args = childModule.ToWithdrawFromNeedArguments(on_chain.WithdrawFromNeedArguments{
				LocalPool:  localPoolId,
				TargetID:   detail.Target,
				ProposalID: detail.ProposalID,
			})
		}
	}

	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    client,
		Sender:    payment.Actor,
		Module:    module,
		Function:  function,
		ErrLogger: p.errLogger,
		Arguments: args,
	}, ctx)

	var req = util.GenerateRedirectParamRequest{
		OrderCode: payment.TransactionId,
		Status:    data.Status,
		TxBytes:   txBytes,
	}

	if err != nil {
		req.Message = err.Error()
	} else {
		req.Message = "Success"
	}

	return util.GeneratePaymentRedirectUrl(req), nil
}

// func HandlePaymentCallback(w http.ResponseWriter, r *http.Request) {
//     // 1. Lấy thông tin từ PayOS (Query params)
//     params := r.URL.Query()
//     status := params.Get("status")
//     orderCode := params.Get("orderCode")

//     // 2. Kiểm tra thanh toán thành công
//     if status == "PAID" {
//         // Cập nhật DB (Off-chain)
//         updateOrderPaid(orderCode)

//         // 3. Logic BUILD TRANSACTION hiện tại của bạn
//         // Giả sử txBytes của bạn là chuỗi Hex hoặc Base64 dưới 100 ký tự
//         txBytes := blockchainService.BuildUnsignedTx(orderCode)

//         // 4. Redirect về Frontend kèm theo txBytes trên URL
//         // Redirect về trang chuyên biệt để xử lý ký ví
//         targetFrontend := fmt.Sprintf("https://frontend.com",
//                           txBytes, orderCode)

//         http.Redirect(w, r, targetFrontend, http.StatusSeeOther)
//         return
//     }

//     // Xử lý khi thanh toán thất bại
//     http.Redirect(w, r, "https://frontend.com", http.StatusSeeOther)
// }

// import (
// 	"fmt"
// 	"net/url"
// )

// func main() {
// 	orderCode := "123"
// 	txBytes := "Sui/Base64+Data==" // Giả sử txBytes chứa ký tự đặc biệt

// 	// Cách 1: Encode từng phần (Dùng Path Escape nếu bỏ vào giữa URL /123/txBytes)
// 	safeTx := url.PathEscape(txBytes)
// 	fmt.Println("https://frontend.com" + orderCode + "/" + safeTx)
// 	// Kết quả: https://frontend.com123/Sui%2FBase64+Data%3D%3D

// 	// Cách 2: Dùng Query Params (Khuyên dùng vì chuẩn hóa hơn)
// 	params := url.Values{}
// 	params.Add("tx", txBytes)
// 	finalURL := fmt.Sprintf("https://frontend.comconfirm?order=%s&%s", orderCode, params.Encode())
// 	fmt.Println(finalURL)
// 	// Kết quả: https://frontend.comconfirm?order=123&tx=Sui%2FBase64%2BData%3D%3D
// }
