package business

import (
	"context"
	"fmt"
	"log"
	"os"
	"raise-child/constants/env"
	"raise-child/constants/shared"
	"raise-child/interfaces/business"
	i_repository "raise-child/interfaces/repository"
	"raise-child/model/entities"
	"raise-child/repository"
	"raise-child/util"
	"raise-child/util/db"
	on_chain "raise-child/util/on_chain"
	"strconv"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
)

type backgroundService struct {
	backgroundChildrenWithdrawRepo i_repository.IBackgroundChildrenWithdrawProposalRequestRepository
	leaderNotiRepo                 i_repository.ILeaderNotiRepository
	registrationRequestRepo        i_repository.IRegistrationRequestRepository
	centerRequestRepo              i_repository.ICenterRequestRepository
	uploadChildRequestRepo         i_repository.IUploadChildRequestRepository
	clients                        map[string]sui.ISuiAPI
	errLogger                      *log.Logger
}

func logMissingEnvKeys(logger *log.Logger, scope string, keys ...string) {
	for _, key := range keys {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			logger.Println(fmt.Sprintf("%s missing env: %s", scope, key))
		}
	}
}

func initializeBackgroundService(
	backgroundChildrenWithdrawRepo i_repository.IBackgroundChildrenWithdrawProposalRequestRepository,
	leaderNotiRepo i_repository.ILeaderNotiRepository,
	registrationRequestRepo i_repository.IRegistrationRequestRepository,
	centerRequestRepo i_repository.ICenterRequestRepository,
	uploadChildRequestRepo i_repository.IUploadChildRequestRepository,
	clients map[string]sui.ISuiAPI,
	errLogger *log.Logger,
) business.IBackgroundService {
	return &backgroundService{
		backgroundChildrenWithdrawRepo: backgroundChildrenWithdrawRepo,
		registrationRequestRepo:        registrationRequestRepo,
		centerRequestRepo:              centerRequestRepo,
		uploadChildRequestRepo:         uploadChildRequestRepo,
		clients:                        clients,
		errLogger:                      errLogger,
	}
}

func GenerateBackgroundService() (business.IBackgroundService, error) {
	var errLogger = util.GetLogConfig(shared.ERROR_LEVEL)

	cnn, err := db.ConnectDB(errLogger, db.InitializePostgreSQL())
	if err != nil {
		return nil, err
	}

	return initializeBackgroundService(
		repository.InitializeBackgroundChildrenWithdrawRequestRepository(cnn, errLogger),
		repository.InitializeLeaderNotiRepository(cnn, errLogger),
		repository.InitializeRegistrationRequestRepo(cnn, errLogger),
		repository.InitializeCenterRequestRepository(cnn, errLogger),
		repository.InitializeUploadChildRequestRepo(cnn, errLogger),
		_networkAliases,
		errLogger,
	), nil
}

// ProcessCreateChildrenWithdrawProposals implements business.IBackgroundService.
func (b *backgroundService) ProcessCreateChildrenWithdrawProposals(ctx context.Context) {
	b.errLogger.Println("Create children withdraws call")
	var curTime time.Time = time.Now()
	var curTimeEndOfDate time.Time = util.ToEndOfDate(curTime)
	if !isChildrenWithdrawProposalDateValid(curTime) {
		return
	}

	reqs, _ := b.backgroundChildrenWithdrawRepo.GetCurrentPendingRequests(ctx)
	if reqs == nil {
		return
	}

	var client = b.clients[constant.SuiTestnet]
	manage, _ := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
		ErrLogger: b.errLogger,
	}, ctx)
	if manage == nil {
		b.errLogger.Println("Nil manage")
		return
	}

	pool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: b.errLogger,
	}, ctx)
	if pool == nil {
		b.errLogger.Println("Nil Pool")
		return
	}

	localPools, _ := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.LocalPools,
		ErrLogger: b.errLogger,
	}, ctx)
	if localPools == nil {
		b.errLogger.Println("Nil Local Pools")
		return
	}

	booksNeedWithdrawDates, _ := on_chain.GetOnChainObject[entities.BooksNeedWithdrawDates](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.BOOKS_NEED_WITHDRAW_DATES_ID),
		ErrLogger: b.errLogger,
	}, ctx)
	if booksNeedWithdrawDates == nil {
		b.errLogger.Println("Nil Books Need Withdraw Dates")
		return
	}

	healthNeedWithdrawDate, _ := on_chain.GetOnChainObject[entities.HealthInsuranceNeedWithdrawDate](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.HEALTH_INSURANCE_NEED_WITHDRAW_DATE_ID),
		ErrLogger: b.errLogger,
	}, ctx)
	if healthNeedWithdrawDate == nil {
		b.errLogger.Println("Nil Health Insurance Need Withdraw Date")
		return
	}

	var modules []string
	var functions []string
	var args [][]interface{}
	var module = on_chain.InitializeModuleChild()
	for _, req := range reqs {
		var centerId string
		for i, region := range manage.LocalRegions {
			if region == req.Region {
				centerId = manage.ChildrenCenters[i]
				break
			}
		}

		center, _ := on_chain.GetOnChainObject[entities.Center](on_chain.GetOnChainObjectRequest{
			Client:    client,
			ObjectId:  centerId,
			ErrLogger: b.errLogger,
		}, ctx)
		if center == nil {
			return
		}

		children, _ := on_chain.GetOnChainObjects[entities.Child](on_chain.GetOnChainObjectsRequest{
			Client:    client,
			ObjectIds: center.ChildIDs,
			ErrLogger: b.errLogger,
		}, ctx)
		if children == nil {
			return
		}

		var localPoolId string
		for _, localPool := range localPools {
			if localPool.Region == req.Region {
				localPoolId = localPool.ID.ID
				break
			}
		}

		for _, child := range children {
			// Books Needs Withdraw
			for i := 0; i <= 1; i++ {
				if data := b.prepareDataCreateBooksNeedWithdrawProposal(
					child.BooksNeeds[i],
					req.ActorAddress,
					localPoolId,
					module,
					curTimeEndOfDate,
					client,
					booksNeedWithdrawDates,
					ctx,
				); data != nil {
					modules = append(modules, module.GetModule())
					functions = append(functions, module.GetFunctionCreateChildBooksNeedWithdrawProposal())
					args = append(args, data)
				}
			}

			// Health Insurance Need Withdraw
			if data := b.prepareDataCreateHealthInsuranceNeedWithdrawProposal(
				child.HealthInsuranceNeed,
				req.ActorAddress,
				localPoolId,
				module,
				curTimeEndOfDate,
				client,
				healthNeedWithdrawDate,
				ctx,
			); data != nil {
				modules = append(modules, module.GetModule())
				functions = append(functions, module.GetFunctionCreateChildHealthInsuranceNeedWithdrawProposal())
				args = append(args, data)
			}

			// Meal Need Withdraw
			if data := b.prepareDataCreateMealNeedWithdrawProposal(
				child.MealNeed,
				req.ActorAddress,
				localPoolId,
				module,
				curTimeEndOfDate,
				client,
				ctx,
			); data != nil {
				modules = append(modules, module.GetModule())
				functions = append(functions, module.GetFunctionCreateChildMealNeedWithdrawProposal())
				args = append(args, data)
			}
		}
	}

	var errExecuteTxs error
	if args != nil && len(functions) > 0 && len(modules) > 0 {
		for i := 1; i <= 3; i++ {
			if err := on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
				Client:    client,
				Modules:   modules,
				Functions: functions,
				Arguments: args,
				ErrLogger: b.errLogger,
			}, ctx); err == nil {
				break
			} else {
				errExecuteTxs = err
			}
		}
	}

	if errExecuteTxs == nil {
		if curTime.Hour() == 23 && curTime.Minute() >= 55 {
			for i := 1; i <= 3; i++ {
				if b.backgroundChildrenWithdrawRepo.SetRequestsExecuted(reqs, ctx) == nil {
					return
				}
			}
		}
	}
}

// ProcessRefundVotePower implements business.IBackgroundService.
func (b *backgroundService) ProcessRefundVotePower(ctx context.Context) {
	var client = b.clients[constant.SuiTestnet]

	pool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.POOL_ID),
		ErrLogger: b.errLogger,
	}, ctx)
	if pool == nil {
		return
	}

	pendingProposals, _ := on_chain.GetOnChainObjects[entities.WithdrawProposal](on_chain.GetOnChainObjectsRequest{
		Client:    client,
		ObjectIds: pool.PendingWithdrawProposals,
		ErrLogger: b.errLogger,
	}, ctx)
	if len(pendingProposals) == 0 {
		return
	}

	rateObj, _ := on_chain.GetOnChainObject[entities.AllowedFundedWithdrawRateObject](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  os.Getenv(env.ALLOWED_FUNDED_WITHDRAW_RATE_OBJECT_ID),
		ErrLogger: b.errLogger,
	}, ctx)
	if rateObj == nil {
		return
	}

	allowedRate, _ := strconv.ParseInt(rateObj.Rate, 10, 64)

	var curTime time.Time = time.Now()
	var modules []string
	var functions []string
	var args [][]interface{}
	var module = on_chain.InitializeModuleRefund()
	for _, proposal := range pendingProposals {
		miliSecsClosedAt, _ := strconv.ParseInt(proposal.ClosedAt, 10, 64)
		if util.MilliSecToTime(miliSecsClosedAt).After(curTime) {
			continue
		}

		withdrawAmount, _ := strconv.ParseInt(proposal.WithdrawAmount, 10, 64)
		totalApproveWeight, _ := strconv.ParseInt(proposal.TotalApproveWeight, 10, 64)
		if totalApproveWeight >= withdrawAmount*allowedRate/100 {
			continue
		}

		var targetId string
		switch proposal.Purpose {
		case string(entities.SPECIAL_NEED_CAMPAIGN_WITHDRAW_PROPOSAL_PURPOSE):
			functions = append(functions, module.GetFunctionRefundChildSpecialNeedCampaignVotePower())
			targetId = proposal.TargetID
		case string(entities.POOL_CAMPAIGN_WITHDRAW_PROPOSAL_PURPOSE):
			functions = append(functions, module.GetFunctionRefundPoolCampaignVotePower())
			targetId = proposal.TargetID
		case string(entities.POOL_WITHDRAW_PROPOSAL_PURPOSE):
			functions = append(functions, module.GetFunctionRefundPoolVotePower())
			if proposal.IsFromLocalPool {
				targetId = proposal.TargetID
			} else {
				targetId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
			}
		}

		modules = append(modules, module.GetModule())
		args = append(args, module.ToRefundVotePowerArguments(on_chain.RefundVotePowerArguments{
			TargetID:   targetId,
			ProposalID: proposal.ID.ID,
		}))
	}

	b.errLogger.Println("Refund vote power call")
	if args != nil && len(functions) > 0 && len(modules) > 0 {
		for i := 1; i <= 3; i++ {
			if on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
				Client:    client,
				Modules:   modules,
				Functions: functions,
				Arguments: args,
				ErrLogger: b.errLogger,
			}, ctx) == nil {
				break
			}
		}
	}
}

// ProcessBackgroundCenterRequests implements business.IBackgroundService.
func (b *backgroundService) ProcessBackgroundCenterRequests(ctx context.Context) {
	pendingRes, err := b.centerRequestRepo.GetPendingRequests(ctx)
	if err != nil {
		return
	}

	b.errLogger.Println("Background center call")

	if len(pendingRes) > 0 {
		var refusedReqs, approvedRes []entities.CenterRequest
		for _, req := range pendingRes {
			var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
			if rate >= approve_rate_limit {
				approvedRes = append(approvedRes, req)
			} else {
				refusedReqs = append(refusedReqs, req)
			}
		}

		b.centerRequestRepo.SetRefusedStatuses(refusedReqs, ctx)

		if len(approvedRes) > 0 {
			var client = b.clients[constant.SuiTestnet]
			manage, _ := on_chain.GetOnChainObject[entities.Manage](on_chain.GetOnChainObjectRequest{
				Client:    client,
				ObjectId:  os.Getenv(env.MANAGE_OBJECT_ID),
				ErrLogger: b.errLogger,
			}, ctx)
			if manage == nil {
				return
			}

			var modules []string
			var functions []string
			var args [][]interface{}
			var module = on_chain.InitializeModuleChild()
			for _, req := range approvedRes {
				var startIdx int
				for i, region := range manage.LocalRegions {
					if region == req.Region {
						startIdx = i
						break
					}
				}

				leaders, _ := on_chain.GetOnChainObjects[entities.StaffNft](on_chain.GetOnChainObjectsRequest{
					Client:    client,
					ObjectIds: manage.LocalLeaderNfts[startIdx:],
					ErrLogger: b.errLogger,
				}, ctx)
				if leaders == nil {
					return
				}

				var regionLeaders []string
				for _, leader := range leaders {
					if leader.Region == req.Region {
						regionLeaders = append(regionLeaders, leader.Owner)
					}
				}

				modules = append(modules, module.GetModule())
				functions = append(functions, module.GetFunctionUploadCenter())
				args = append(args, module.ToCreateCenterArguments(on_chain.CreateCenterArguments{
					Region:      req.Region,
					Address:     req.Address,
					PhoneNumber: req.PhoneNumber,
					ImageBlobID: req.ImageBlobID,
					Leaders:     regionLeaders,
					Sender:      req.CreatedBy,
				}))
			}

			if err := on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
				Client:    client,
				Modules:   modules,
				Functions: functions,
				Arguments: args,
				ErrLogger: b.errLogger,
			}, ctx); err != nil {
				return
			}

			for i := 1; i <= 3; i++ {
				if b.centerRequestRepo.SetApprovedStatuses(approvedRes, ctx) == nil {
					return
				}
			}
		}
	}
}

// ProcessBackgroundRegistrationRequests implements business.IBackgroundService.
func (b *backgroundService) ProcessBackgroundRegistrationRequests(ctx context.Context) {
	pendingRes, err := b.registrationRequestRepo.GetPendingRequestsV2(ctx)
	if err != nil {
		return
	}

	b.errLogger.Println("Registration background call")

	if len(pendingRes) > 0 {
		var refusedReqs, approvedReqs []entities.RegistrationRequest
		for _, req := range pendingRes {
			var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
			if rate >= approve_rate_limit {
				approvedReqs = append(approvedReqs, req)
			} else {
				refusedReqs = append(refusedReqs, req)
			}
		}

		b.registrationRequestRepo.SetRefusedStatusesV2(refusedReqs, ctx)

		var modules []string
		var functions []string
		var args [][]interface{}
		var staffModule = on_chain.InitializeModuleStaff()
		var client = b.clients[constant.SuiTestnet]

		for _, req := range approvedReqs {
			modules = append(modules, staffModule.GetModule())

			switch req.RegisterRole {
			case admin_role:
				functions = append(functions, staffModule.GetFunctionRegisterAdmin())
				args = append(args, staffModule.ToRegisterAdminArguments(on_chain.RegisterAdminArguments{
					IdentityCode: req.IdentityCode,
					AvatarBlobID: req.AvatarBlobID,
					FirstName:    req.FirstName,
					LastName:     req.LastName,
					Gender:       req.Gender,
					DateOfBirth:  req.DateOfBirth,
					PhoneNumber:  req.PhoneNumber,
					Email:        req.Email,
					Owner:        req.CreatedBy,
				}))
			case local_leader_role:
				pool, _ := on_chain.GetOnChainObject[entities.MainPool](on_chain.GetOnChainObjectRequest{
					Client:    client,
					ObjectId:  os.Getenv(env.POOL_ID),
					ErrLogger: b.errLogger,
				}, ctx)

				if pool == nil {
					return
				}

				localPools, _ := on_chain.GetOnChainObjects[entities.LocalPool](on_chain.GetOnChainObjectsRequest{
					Client:    client,
					ObjectIds: pool.LocalPools,
					ErrLogger: b.errLogger,
				}, ctx)

				if len(localPools) == 0 {
					return
				}

				var localPoolId string
				for _, localPool := range localPools {
					if localPool.Region == req.Region {
						localPoolId = localPool.ID.ID
						break
					}
				}

				if localPoolId == "" {
					localPoolId = os.Getenv(env.SHARED_LOCAL_POOL_ID)
				}

				functions = append(functions, staffModule.GetFunctionRegisterLeader())
				args = append(args, staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
					LocalPoolID: localPoolId,
					Region:      req.Region,
					RegisterAdminArguments: on_chain.RegisterAdminArguments{
						IdentityCode: req.IdentityCode,
						AvatarBlobID: req.AvatarBlobID,
						FirstName:    req.FirstName,
						LastName:     req.LastName,
						Gender:       req.Gender,
						DateOfBirth:  req.DateOfBirth,
						PhoneNumber:  req.PhoneNumber,
						Email:        req.Email,
						Owner:        req.CreatedBy,
					},
				}))
			case volunteer_role:
				functions = append(functions, staffModule.GetFunctionRegisterVolunteer())
				args = append(args, staffModule.ToRegisterNormalStaffArguments(on_chain.RegisterNormalStaffArguments{
					Region: req.Region,
					RegisterAdminArguments: on_chain.RegisterAdminArguments{
						IdentityCode: req.IdentityCode,
						AvatarBlobID: req.AvatarBlobID,
						FirstName:    req.FirstName,
						LastName:     req.LastName,
						Gender:       req.Gender,
						DateOfBirth:  req.DateOfBirth,
						PhoneNumber:  req.PhoneNumber,
						Email:        req.Email,
						Owner:        req.CreatedBy,
					},
				}))
			}
		}

		b.errLogger.Println("Registration Background call")
		if args != nil && len(functions) > 0 && len(modules) > 0 {
			for i := 1; i <= 3; i++ {
				if on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
					Client:    client,
					Modules:   modules,
					Functions: functions,
					Arguments: args,
					ErrLogger: b.errLogger,
				}, ctx) == nil {
					break
				}
			}
		}

		b.registrationRequestRepo.SetApprovedStatusesV2(approvedReqs, ctx)
	}
}

// ProcessBackgroundUploadChildRequests implements business.IBackgroundService.
func (b *backgroundService) ProcessBackgroundUploadChildRequests(ctx context.Context) {
	// pendingRes, approvedRes, err := b.uploadChildRequestRepo.GetPendingRequests(ctx)
	// if err != nil {
	// 	return
	// }

	// var refusedReqs []entities.BackgroundRecord
	// for _, req := range pendingRes {
	// 	var rate float32 = float32(len(req.Approvers)) / float32(len(req.Approvers)+len(req.Refusers))
	// 	if rate >= approve_rate_limit {
	// 		approvedRes = append(approvedRes, req)
	// 	} else {
	// 		refusedReqs = append(refusedReqs, req)
	// 	}
	// }

	// b.uploadChildRequestRepo.SetRefusedStatuses(refusedReqs, ctx)

	// var modules []string
	// var functions []string
	// var args [][]interface{}
	// var module = on_chain.InitializeModuleManage()
	// for i, req := range approvedRes {
	// 	modules[i] = module.GetModule()
	// 	functions[i] = module.GetFunctionMintUploadCenterCap()
	// 	args = append(args, module.ToMintCapArguments(on_chain.MintCapArguments{
	// 		Recipient: req.Sender,
	// 	}))
	// }

	// if err := on_chain.BuildMultiBackgroundTransactions(on_chain.BuildMultiBackgroundTransactionsRequest{
	// 	Client:    b.clients[constant.SuiTestnet],
	// 	Modules:   modules,
	// 	Functions: functions,
	// 	Arguments: args,
	// 	ErrLogger: b.errLogger,
	// }, ctx); err != nil {
	// 	return
	// }

	// b.uploadChildRequestRepo.SetApprovedStatuses(approvedRes, ctx)
}

func (b *backgroundService) prepareDataCreateBooksNeedWithdrawProposal(needId, sender, poolId string, module on_chain.IModuleChild, curTime time.Time, client sui.ISuiAPI, withdrawDates *entities.BooksNeedWithdrawDates, ctx context.Context) []interface{} {
	need, _ := on_chain.GetOnChainObject[entities.BooksNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  needId,
		ErrLogger: b.errLogger,
	}, ctx)
	if need == nil {
		return nil
	}

	if len(need.WithdrawProposals) > len(need.WithdrawsForNeed) {
		log.Println("Check withdraw ok")
		return nil
	}

	if len(need.Donations) == len(need.WithdrawsForNeed) {
		return nil
	}

	var rawExpectedWithdrawDate string
	if need.Semster == "1" {
		rawExpectedWithdrawDate = fmt.Sprintf("%s/%d", withdrawDates.FirstSemesterDate, curTime.Year())
	} else {
		rawExpectedWithdrawDate = fmt.Sprintf("%s/%d", withdrawDates.SecondSemesterDate, curTime.Year())
	}

	var expectedWithdrawDate time.Time = util.ToStartOfDate(util.RawDateToTime(rawExpectedWithdrawDate))
	if !expectedWithdrawDate.After(curTime) && int(expectedWithdrawDate.AddDate(0, 0, 7).Month()-curTime.Month()) <= 1 {
		var description string = fmt.Sprintf("Withdraw Books Need Semester %s - %d for child %s", need.Semster, curTime.Year(), util.FormatAddress(need.ChildID))
		return module.ToCreateChildNormalNeedWithdrawProposalArguments(on_chain.CreateChildNormalNeedWithdrawProposalArguments{
			NeedID:      needId,
			ChildID:     need.ChildID,
			LocalPool:   poolId,
			Description: description,
			Sender:      sender,
		})
	}

	return nil
}

func (b *backgroundService) prepareDataCreateHealthInsuranceNeedWithdrawProposal(needId, sender, poolId string, module on_chain.IModuleChild, curTime time.Time, client sui.ISuiAPI, withdrawDate *entities.HealthInsuranceNeedWithdrawDate, ctx context.Context) []interface{} {
	need, _ := on_chain.GetOnChainObject[entities.HealthInsuranceNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  needId,
		ErrLogger: b.errLogger,
	}, ctx)
	if need == nil {
		return nil
	}

	if len(need.WithdrawProposals) > len(need.WithdrawsForNeed) {
		log.Println("Check withdraw ok")
		return nil
	}

	if len(need.Donations) == len(need.WithdrawsForNeed) {
		return nil
	}

	var rawExpectedWithdrawDate string = fmt.Sprintf("%s/%d", withdrawDate.ExpectedDate, curTime.Year())
	var expectedWithdrawDate time.Time = util.ToStartOfDate(util.RawDateToTime(rawExpectedWithdrawDate))
	if !expectedWithdrawDate.After(curTime) && int(expectedWithdrawDate.AddDate(0, 0, 7).Month()-curTime.Month()) <= 1 {
		var description string = fmt.Sprintf("Withdraw Health Insurance Need %d for child %s", curTime.Year(), util.FormatAddress(need.ChildID))
		return module.ToCreateChildNormalNeedWithdrawProposalArguments(on_chain.CreateChildNormalNeedWithdrawProposalArguments{
			NeedID:      needId,
			ChildID:     need.ChildID,
			LocalPool:   poolId,
			Description: description,
			Sender:      sender,
		})
	}

	return nil
}

func (b *backgroundService) prepareDataCreateMealNeedWithdrawProposal(needId, sender, poolId string, module on_chain.IModuleChild, curTime time.Time, client sui.ISuiAPI, ctx context.Context) []interface{} {
	need, _ := on_chain.GetOnChainObject[entities.MealNeed](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  needId,
		ErrLogger: b.errLogger,
	}, ctx)
	if need == nil {
		b.errLogger.Println("Nil Meal Need")
		return nil
	}

	if len(need.WithdrawProposals) > len(need.WithdrawsForNeed) {
		log.Println("Check withdraw ok")
		return nil
	}

	totalSupportedMonths, _ := strconv.Atoi(need.TotalSupportedMonths)
	var expectedDuration int = totalSupportedMonths - len(need.WithdrawsForNeed)
	if expectedDuration == 0 {
		return nil
	}

	var previousDuration int = 0
	var expectedDate time.Time
	var foundIdx int = -1
	for i := len(need.Durations) - 1; i >= 0; i-- {
		var duration = need.Durations[0]
		var startPeriod time.Time = util.ToStartOfDate(util.RawDateToTime(duration.Fields.StartPeriod))
		var endPeriod time.Time = util.ToEndOfDate(util.RawDateToTime(duration.Fields.EndPeriod))
		var startMonth int = int(startPeriod.Month())
		var endMonth int = int(endPeriod.Month())
		if endMonth == 1 { // To next year
			endMonth = 13
		}

		var currentDuration int = endMonth - startMonth
		var totalDuration int = currentDuration + previousDuration
		var months int = totalDuration - expectedDuration
		if months >= 0 {
			var startDate = startPeriod.AddDate(0, months, 0)
			expectedDate = startDate.AddDate(0, 0, -3)
			foundIdx = i
			break
		}

		previousDuration = totalDuration
	}

	if !expectedDate.After(curTime) && int(expectedDate.AddDate(0, 0, 7).Month()-curTime.Month()) <= 1 {
		var description string = fmt.Sprintf("Withdraw Meal Need for child %s from %s to %s", util.FormatAddress(need.ChildID), need.Durations[foundIdx].Fields.StartPeriod, need.Durations[foundIdx].Fields.EndPeriod)
		return module.ToCreateChildNormalNeedWithdrawProposalArguments(on_chain.CreateChildNormalNeedWithdrawProposalArguments{
			NeedID:      needId,
			ChildID:     need.ChildID,
			LocalPool:   poolId,
			Description: description,
			Sender:      sender,
		})
	}

	return nil
}
