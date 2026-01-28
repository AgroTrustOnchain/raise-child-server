package business

import (
	"context"
	"errors"
	"log"
	"raise-child/constants/noti"
	on_chain "raise-child/util/on_chain"
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
