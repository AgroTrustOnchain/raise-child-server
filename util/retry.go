package util

import "time"

type RetryFunc[T any] func() (T, error)

type RetryFuncNoResult func() error

type RetryOption struct {
	MaxAttempts int
	Delay       time.Duration
}

var DefaultRetryOption = RetryOption{
	MaxAttempts: 3,
	Delay:       1 * time.Second,
}

func Retry[T any](opt RetryOption, fn RetryFunc[T]) (T, error) {
	if opt.MaxAttempts <= 0 {
		opt.MaxAttempts = DefaultRetryOption.MaxAttempts
	}

	if opt.Delay <= 0 {
		opt.Delay = DefaultRetryOption.Delay
	}

	var res T
	var err error
	for i := 1; i <= opt.MaxAttempts; i++ {
		res, err = fn()
		if err == nil {
			break
		}

		if i < opt.MaxAttempts {
			time.Sleep(opt.Delay)
		}
	}

	return res, err
}

func RetryWithNoResult(opt RetryOption, fn RetryFuncNoResult) error {
	if opt.MaxAttempts <= 0 {
		opt.MaxAttempts = DefaultRetryOption.MaxAttempts
	}

	if opt.Delay <= 0 {
		opt.Delay = DefaultRetryOption.Delay
	}

	var err error
	for i := 1; i <= opt.MaxAttempts; i++ {
		if err = fn(); err == nil {
			break
		}

		if i < opt.MaxAttempts {
			time.Sleep(opt.Delay)
		}
	}

	return err
}
