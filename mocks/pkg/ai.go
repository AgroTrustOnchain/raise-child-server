package pkg

import (
	"context"
	"raise-child/util/ai"

	"github.com/stretchr/testify/mock"
)

type aiMockClient struct {
	mock.Mock
}

func InitializeAiMockClient() *aiMockClient {
	return &aiMockClient{}
}

// ValidateChildSpecialNeedProposal implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateChildSpecialNeedProposal(req ai.ValidateChildSpecialNeedProposal, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateChildSpecialNeedProposal, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidatePoolCampaign implements ai.IAiClientProvider.
func (a *aiMockClient) ValidatePoolCampaign(req ai.ValidatePoolCampaign, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidatePoolCampaign, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidateCreateCenterRequest implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateCreateCenterRequest(req ai.ValidateCreateCenterRequest, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateCreateCenterRequest, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidateProvideMealForChildTaskProof implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateProvideMealForChildTaskProof(req ai.ValidateProvideMealForChildTaskProof, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateProvideMealForChildTaskProof, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidateRegistrationRequest implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateRegistrationRequest(req ai.ValidateRegistrationRequest, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateRegistrationRequest, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidateTaskProof implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateTaskProof(req ai.ValidateTaskProof, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateTaskProof, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidateUploadChildRequest implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateUploadChildRequest(req ai.ValidateUploadChildRequest, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateUploadChildRequest, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}

// ValidateWithdrawProposal implements ai.IAiClientProvider.
func (a *aiMockClient) ValidateWithdrawProposal(req ai.ValidateWithdrawProposal, ctx context.Context) string {
	var mockData = a.Called(req, ctx)
	if mockFunc, ok := mockData.Get(0).(func(ai.ValidateWithdrawProposal, context.Context) string); ok {
		return mockFunc(req, ctx)
	}

	return ""
}
