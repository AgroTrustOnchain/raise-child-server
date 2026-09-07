package repository

import (
	"context"
	"raise-child/interfaces/repository"
	"raise-child/model/entities"

	"github.com/stretchr/testify/mock"
)

type offChainWithdrawProposalMockRepo struct {
	mock.Mock
}

func InializeOffChainWithdrawProposalMockRepo() repository.IOffChainWithdrawProposalRepository {
	return &offChainWithdrawProposalMockRepo{}
}

// CreateOffChainWithdrawProposal implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalMockRepo) CreateOffChainWithdrawProposal(proposal entities.OffChainWithdrawProposal, ctx context.Context) error {
	var mockData = o.Called(proposal, ctx)

	if mockFunc, ok := mockData.Get(0).(func(entities.OffChainWithdrawProposal, context.Context) error); ok {
		return mockFunc(proposal, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}

// GetOffChainWithdrawProposal implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalMockRepo) GetOffChainWithdrawProposal(id string, ctx context.Context) (*entities.OffChainWithdrawProposal, error) {
	var mockData = o.Called(id, ctx)

	var res1 *entities.OffChainWithdrawProposal
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.OffChainWithdrawProposal); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.OffChainWithdrawProposal)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// GetOffChainWithdrawProposalByProposal implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalMockRepo) GetOffChainWithdrawProposalByProposal(id string, ctx context.Context) (*entities.OffChainWithdrawProposal, error) {
	var mockData = o.Called(id, ctx)

	var res1 *entities.OffChainWithdrawProposal
	if mockFunc, ok := mockData.Get(0).(func(string, context.Context) *entities.OffChainWithdrawProposal); ok {
		res1 = mockFunc(id, ctx)
	} else {
		res1 = mockData.Get(0).(*entities.OffChainWithdrawProposal)
	}

	var res2 error
	if mockFunc, ok := mockData.Get(1).(func(string, context.Context) error); ok {
		res2 = mockFunc(id, ctx)
	} else {
		res2 = mockData.Error(1)
	}

	return res1, res2
}

// SetOnChainProposalIdAfterExecuteTx implements repository.IOffChainWithdrawProposalRepository.
func (o *offChainWithdrawProposalMockRepo) SetOnChainProposalIdAfterExecuteTx(id string, proposalId string, ctx context.Context) error {
	var mockData = o.Called(id, proposalId, ctx)

	if mockFunc, ok := mockData.Get(0).(func(string, string, context.Context) error); ok {
		return mockFunc(id, proposalId, ctx)
	}

	if err, ok := mockData.Error(0).(error); ok {
		return err
	}

	return nil
}
