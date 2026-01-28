package business

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
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
	paymentRepo i_repository.IPaymentRepository
	profileRepo i_repository.IProfileRepository
	clients     map[string]sui.ISuiAPI
	errLogger   *log.Logger
}

func InitializePaymentService(db *sql.DB, errLogger *log.Logger) business.IPaymentService {
	return &paymentService{
		paymentRepo: repository.InitializePaymentRepository(db, errLogger),
		profileRepo: repository.InitializeProfileRepository(db, errLogger),
		clients:     _networkAliases,
		errLogger:   errLogger,
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

// CallbackTx implements business.IPaymentService.
func (p *paymentService) CallbackTx(id string, ctx context.Context) (response.BuildTransactionResponse, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)

	payment, err := p.paymentRepo.GetPaymentById(id, ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if payment == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	// Expired
	if payment.ExpiredAt.Before(time.Now()) {
		return response.BuildTransactionResponse{}, errors.New(noti.PAYMENT_EXPIRED_MESSAGE)
	}

	data, err := payos.GetPaymentLinkInformation(payment.TransactionId)
	if err != nil {
		p.errLogger.Println("Error while get payos payment link information: " + err.Error())
		return response.BuildTransactionResponse{}, errors.New(noti.INTERNALL_ERR_MSG)
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
		return response.BuildTransactionResponse{}, genericErr
	}

	profile, err := p.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if profile == nil {
		return response.BuildTransactionResponse{}, genericErr
	}

	payment.UpdatedAt = time.Now()
	if err := p.paymentRepo.UpdatePayment(*payment, ctx); err != nil {
		return response.BuildTransactionResponse{}, err
	}

	if data.Status == shared.PAYOS_CANCELLED_STATUS {
		return response.BuildTransactionResponse{}, errors.New(noti.PAYMENT_CANCEL_MESSAGE)
	}

	var module = on_chain.InitializeModulePool()
	txBytes, err := on_chain.BuildTransaction(on_chain.BuildTransactionRequest{
		Client:    p.clients[constant.SuiTestnet],
		Sender:    ctx.Value("address").(string),
		Module:    module.GetModule(),
		Function:  module.GetFunctionDonateToPool(),
		ErrLogger: p.errLogger,
		Arguments: module.ToDonateToPoolArguments(on_chain.DonateToPoolArguments{
			Amount:      payment.Amount,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Gender:      profile.Gender,
			PhoneNumber: profile.PhoneNumber,
			Email:       profile.Email,
			Message:     payment.Message,
		}),
	}, ctx)

	return response.BuildTransactionResponse{
		TxBytes: txBytes,
	}, err
}

// Donate implements business.IPaymentService.
func (p *paymentService) Donate(req request.DonateRequest, ctx context.Context) (string, error) {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	if !utils.IsValidSuiAddress(models.SuiAddress(req.PoolId)) {
		return "", genericErr
	}
	profile, err := p.profileRepo.GetProfile(ctx.Value("sub").(string), ctx)
	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", genericErr
	}

	if profile.IdentityCode == "" {
		return "", errors.New(noti.NOT_UPLOADED_PROFILE_MESSAGE)
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
		return "", errors.New(noti.INTERNALL_ERR_MSG)
	}

	var curTime time.Time = time.Now()
	if err := p.paymentRepo.CreatePayment(entities.Payment{
		ID:            paymentId,
		Actor:         ctx.Value("address").(string),
		Target:        req.PoolId,
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
	}, ctx); err != nil {
		return "", err
	}

	return data.CheckoutUrl, nil
}
