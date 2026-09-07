package business

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/dtos/request"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"time"

	"raise-child/constants/env"
	"raise-child/constants/noti"
	"raise-child/constants/shared"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type onChainService struct {
	withdrawRepo     i_repository.IOffChainWithdrawProposalRepository
	centerRepo       i_repository.ICenterRequestRepository
	uploadChildRepo  i_repository.IUploadChildRequestRepository
	registrationRepo i_repository.IRegistrationRequestRepository
	leaderNotiRepo   i_repository.ILeaderNotiRepository
	clients          map[string]sui.ISuiAPI
	errLogger        *log.Logger
}

func initializeOnChainService(
	withdrawRepo i_repository.IOffChainWithdrawProposalRepository,
	centerRepo i_repository.ICenterRequestRepository,
	uploadChildRepo i_repository.IUploadChildRequestRepository,
	registrationRepo i_repository.IRegistrationRequestRepository,
	leaderNotiRepo i_repository.ILeaderNotiRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IOnChainService {
	return &onChainService{
		withdrawRepo:     withdrawRepo,
		centerRepo:       centerRepo,
		uploadChildRepo:  uploadChildRepo,
		registrationRepo: registrationRepo,
		leaderNotiRepo:   leaderNotiRepo,
		clients:          clients,
		errLogger:        errLogger,
	}
}

func GenerateOnChainService() (business.IOnChainService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeOnChainService(
		repository.InitializeOffChainWithdrawProposalRepository(cnn, errLogger),
		repository.InitializeCenterRequestRepository(cnn, errLogger),
		repository.InitializeUploadChildRequestRepo(cnn, errLogger),
		repository.InitializeRegistrationRequestRepo(cnn, errLogger),
		repository.InitializeLeaderNotiRepository(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ExecuteTransaction implements business.IOnChainService.
func (o *onChainService) ExecuteTransaction(req request.ExecuteTransactionRequest, ctx context.Context) error {
	var genericErr error = errors.New(noti.GENERIC_ERROR_WARN_MSG)
	var curTime time.Time = time.Now()
	if req.ProposalID != "" {
		proposal, err := o.withdrawRepo.GetOffChainWithdrawProposal(req.ProposalID, ctx)
		if err != nil {
			return err
		}

		if proposal.ProposalID != nil {
			return genericErr
		}
	}

	res, err := on_chain.ExecuteTransaction(on_chain.ExecuteTransactionRequest{
		Client:    o.clients[constant.SuiTestnet],
		TxBytes:   req.TxBytes,
		Signature: []string{req.Signature},
		ErrLogger: o.errLogger,
	}, ctx)
	if err != nil {
		return err
	}

	if req.ProposalID != "" {
		var events = res.Events
		if len(events) == 0 {
			return genericErr
		}

		var module = on_chain.InitializeModulePool()
		var eventType string = fmt.Sprintf("%s::%s::%s", os.Getenv(env.PACKAGE_ID), module.GetModule(), module.GetWithdrawProposalEventEmittedStruct())
		for _, event := range events {
			if event.Type == eventType {
				if onChainProposal, ok := event.ParsedJson["id"].(string); ok {
					o.withdrawRepo.SetOnChainProposalIdAfterExecuteTx(req.ProposalID, onChainProposal, ctx)
					break
				}
			}
		}
	} else if req.CenterReq != "" {
		req, err := o.centerRepo.GetRequest(req.CenterReq, ctx)
		if err != nil {
			return err
		}

		req.UpdatedAt = curTime

		o.centerRepo.UpdateRegistrationRequest(*req, ctx)
	} else if req.UploadChildReq != "" {
		req, err := o.uploadChildRepo.GetUploadChildRequest(req.UploadChildReq, ctx)
		if err != nil {
			return err
		}

		if req.IsConfirmUpload {
			return genericErr
		}

		req.IsConfirmUpload = true
		req.UpdatedAt = curTime

		o.uploadChildRepo.UpdateUploadChildRequest(*req, ctx)
	} else if req.RegistrationReq != "" {
		req, err := o.registrationRepo.GetRegistrationRequest(req.RegistrationReq, ctx)
		if err != nil {
			return err
		}

		if req.IsConfirmRegister {
			return genericErr
		}

		req.IsConfirmRegister = true
		req.UpdatedAt = curTime

		o.registrationRepo.UpdateRegistrationRequest(*req, ctx)

		if req.RegisterRole == local_leader_role {
			o.leaderNotiRepo.AssignLeader(req.CreatedBy, req.Region, ctx)
		}
	}

	return nil
}
