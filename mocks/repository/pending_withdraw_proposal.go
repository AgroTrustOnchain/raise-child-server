package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type pendingWithdrawProposalMockRepo struct {
	mock.Mock
}

func InitializePendingWithdrawProposalMockRepo() *pendingWithdrawProposalMockRepo {
	return &pendingWithdrawProposalMockRepo{}
}

// CreatePendingWithdrawProposal implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalMockRepo) CreatePendingWithdrawProposal(proposal entities.PendingWithdrawProposal, ctx context.Context) error {
	var mockData = p.Called(proposal, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.PendingWithdrawProposal, context.Context) error); ok {
		return mockFunc(proposal, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetPendingWithdrawProposal implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalMockRepo) GetPendingWithdrawProposal(id string, ctx context.Context) (*entities.PendingWithdrawProposal, error) {
	var mockData = p.Called(id, ctx)

	var res1 *entities.PendingWithdrawProposal
	if val := mockData.Get(0); val != nil {
		// 1. Check if it's a dynamic function
		if mockFunc, ok := val.(func(string, context.Context) *entities.PendingWithdrawProposal); ok {
			res1 = mockFunc(id, ctx)
		} else {
			// 2. Otherwise, assert it to the slice type
			res1 = val.(*entities.PendingWithdrawProposal)
		}
	} else {
		// If it is nil, res1 remains the zero-value for a slice (which is nil)
		res1 = nil
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetPendingWithdrawProposals implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalMockRepo) GetPendingWithdrawProposals(req request.GetPendingWithdrawProposalsRequest, ctx context.Context) ([]entities.PendingWithdrawProposal, int, error) {
	var mockData = p.Called(req, ctx)

	var res1 []entities.PendingWithdrawProposal
	if val := mockData.Get(0); val != nil {
		// 1. Check if it's a dynamic function
		if mockFunc, ok := val.(func(request.GetPendingWithdrawProposalsRequest, context.Context) []entities.PendingWithdrawProposal); ok {
			res1 = mockFunc(req, ctx)
		} else {
			// 2. Otherwise, assert it to the slice type
			res1 = val.([]entities.PendingWithdrawProposal)
		}
	} else {
		// If it is nil, res1 remains the zero-value for a slice (which is nil)
		res1 = nil
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(request.GetPendingWithdrawProposalsRequest, context.Context) int); ok {
		res2 = mockFunc(req, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetPendingWithdrawProposalsRequest, context.Context) error); ok {
		res3 = mockFunc(req, ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// IsPendingWithdrawProposalProposedWithSpecificInfo implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalMockRepo) IsPendingWithdrawProposalProposedWithSpecificInfo(purpose string, target string, description string, withdrawAmount int64, ctx context.Context) (bool, error) {
	var mockData = p.Called(purpose, target, description, withdrawAmount, ctx)

	var res1 bool
	if mockFunc, ok := mockData.Get(0).(func(string, string, string, int64, context.Context) bool); ok {
		res1 = mockFunc(purpose, target, description, withdrawAmount, ctx)
	} else {
		res1 = mockData.Get(0).(bool)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, string, string, int64, context.Context) error); ok {
		res2 = mockFunc(purpose, target, description, withdrawAmount, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// UpdatePendingWithdrawProposal implements repository.IPendingWithdrawProposalRepository.
func (p *pendingWithdrawProposalMockRepo) UpdatePendingWithdrawProposal(proposal entities.PendingWithdrawProposal, ctx context.Context) error {
	var mockData = p.Called(proposal, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.PendingWithdrawProposal, context.Context) error); ok {
		return mockFunc(proposal, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
