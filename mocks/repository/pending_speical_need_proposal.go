package repository

import (
	"context"
	"raise-child/model/dtos/request"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type pendingChildSpecialNeedProposalMockRepo struct {
	mock.Mock
}

func InitializeChildPendingSpecialProposalMockRepo() *pendingChildSpecialNeedProposalMockRepo {
	return &pendingChildSpecialNeedProposalMockRepo{}
}

// CreatePendingChildSpecialNeedProposal implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalMockRepo) CreatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error {
	var mockData = p.Called(proposal, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.PendingChildSpecialNeedProposal, context.Context) error); ok {
		return mockFunc(proposal, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetPendingChildSpecialNeedProposal implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalMockRepo) GetPendingChildSpecialNeedProposal(id string, ctx context.Context) (*entities.PendingChildSpecialNeedProposal, error) {
	var mockData = p.Called(id, ctx)

	var res1 *entities.PendingChildSpecialNeedProposal
	if val := mockData.Get(0); val != nil {
		// 1. Check if it's a dynamic function
		if mockFunc, ok := val.(func(string, context.Context) *entities.PendingChildSpecialNeedProposal); ok {
			res1 = mockFunc(id, ctx)
		} else {
			// 2. Otherwise, assert it to the slice type
			res1 = val.(*entities.PendingChildSpecialNeedProposal)
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

// GetPendingChildSpecialNeedProposals implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalMockRepo) GetPendingChildSpecialNeedProposals(req request.GetPendingChildSpecialNeedProposalsRequest, ctx context.Context) ([]entities.PendingChildSpecialNeedProposal, int, error) {
	var mockData = p.Called(req, ctx)

	var res1 []entities.PendingChildSpecialNeedProposal
	if val := mockData.Get(0); val != nil {
		// 1. Check if it's a dynamic function
		if mockFunc, ok := val.(func(request.GetPendingChildSpecialNeedProposalsRequest, context.Context) []entities.PendingChildSpecialNeedProposal); ok {
			res1 = mockFunc(req, ctx)
		} else {
			// 2. Otherwise, assert it to the slice type
			res1 = val.([]entities.PendingChildSpecialNeedProposal)
		}
	} else {
		// If it is nil, res1 remains the zero-value for a slice (which is nil)
		res1 = nil
	}

	var res2 int
	if mockFunc, ok := mockData.Get(1).(func(request.GetPendingChildSpecialNeedProposalsRequest, context.Context) int); ok {
		res2 = mockFunc(req, ctx)
	} else {
		res2 = mockData.Get(1).(int)
	}

	var res3 error
	if mockFunc, ok := mockData.Get(2).(func(request.GetPendingChildSpecialNeedProposalsRequest, context.Context) error); ok {
		res3 = mockFunc(req, ctx)
	} else {
		res3 = mockData.Error(2)
	}

	return res1, res2, res3
}

// UpdatePendingChildSpecialNeedProposal implements repository.IPendingChildSpecialNeedProposalRepository.
func (p *pendingChildSpecialNeedProposalMockRepo) UpdatePendingChildSpecialNeedProposal(proposal entities.PendingChildSpecialNeedProposal, ctx context.Context) error {
	var mockData = p.Called(proposal, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.PendingChildSpecialNeedProposal, context.Context) error); ok {
		return mockFunc(proposal, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
