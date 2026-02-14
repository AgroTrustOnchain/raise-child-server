package business

import (
	"context"
	"errors"
	"log"
	"raise-child/constants/noti"
	"raise-child/model/entities"
	on_chain "raise-child/util/on_chain"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
)

const (
	child_max_age_accepted int = 17
	child_min_age_accepted int = 2
)

func validateGetOnChainObject[T any](client sui.ISuiAPI, id string, errLogger *log.Logger, ctx context.Context) error {
	obj, err := on_chain.GetOnChainObject[T](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: errLogger,
	}, ctx)

	// Object not found
	if obj == nil {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return err
}

func getOnChainObject[T any](client sui.ISuiAPI, id string, errLogger *log.Logger, ctx context.Context) (*T, error) {
	obj, err := on_chain.GetOnChainObject[T](on_chain.GetOnChainObjectRequest{
		Client:    client,
		ObjectId:  id,
		ErrLogger: errLogger,
	}, ctx)

	// Object not found
	if obj == nil {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return obj, err
}

func isChildAgeInSupport(yearOfBirth int) bool {
	var curTime time.Time = time.Now()
	return curTime.Year()-yearOfBirth <= child_max_age_accepted && curTime.Year()-yearOfBirth >= child_min_age_accepted
}

func isProposalRateAvailableToConfirm(dao entities.DaoStruct, approverNum, refuserNum int, approveWeight, refuseWeight string) bool {
	minVoters, _ := strconv.Atoi(dao.MinVoters)
	numberRate, _ := strconv.ParseInt(dao.MinApprovedRate, 10, 64)
	var minRate float32 = float32(numberRate / 10000)

	approveWeightNumber, _ := strconv.ParseInt(approveWeight, 10, 64)
	refuseWeightNumber, _ := strconv.ParseInt(refuseWeight, 10, 64)
	var totalWeight int64 = approveWeightNumber + refuseWeightNumber

	return approverNum+refuserNum >= minVoters && float32(approveWeightNumber/totalWeight) >= minRate
}
