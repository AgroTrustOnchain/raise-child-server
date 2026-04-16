package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

type aiClient struct {
	geminiProvider *geminiProvider
	errLogger      *log.Logger
}

var _aiClient *aiClient

type ImageValidateCase string

const ()

type IAiClientProvider interface {
	ValidateUploadChildRequest(req ValidateUploadChildRequest, ctx context.Context) string
	ValidateCreateCenterRequest(req ValidateCreateCenterRequest, ctx context.Context) string
	ValidateRegistrationRequest(req ValidateRegistrationRequest, ctx context.Context) string
	ValidateTaskProof(req ValidateTaskProof, ctx context.Context) string
	ValidateProvideMealForChildTaskProof(req ValidateProvideMealForChildTaskProof, ctx context.Context) string
	ValidateWithdrawProposal(req ValidateWithdrawProposal, ctx context.Context) string
	ValidateChildSpecialNeedProposal(req ValidateChildSpecialNeedProposal, ctx context.Context) string
	ValidatePoolCampaign(req ValidatePoolCampaign, ctx context.Context) string
}

func InitializeAiProvider(ctx context.Context, errLogger *log.Logger) IAiClientProvider {
	if _aiClient == nil {
		_aiClient = &aiClient{
			geminiProvider: initializeGeminiClient(ctx, errLogger),
			errLogger:      errLogger,
		}
	}

	return _aiClient
}

// ValidateUploadChildRequest implements IAiClientProvider.
func (a *aiClient) ValidateUploadChildRequest(req ValidateUploadChildRequest, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", upload_child_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Child Avatar"), genai.ImageData("jpeg", req.AvatarBytesImage))
	prompt = append(prompt, genai.Text("Label: Child Birth Certificate"), genai.ImageData("jpeg", req.ChildBirthCertificateBytesImage))
	prompt = append(prompt, genai.Text("Label: Child First Guardian Identity Card"), genai.ImageData("jpeg", []byte(req.FirstGuardian.IdentityCardBytesImage)))
	if req.SecondGuardian != nil {
		prompt = append(prompt, genai.Text("Label: Child Second Guardian Identity Card"), genai.ImageData("jpeg", []byte(req.SecondGuardian.IdentityCardBytesImage)))
	}

	return a.processPrompt(prompt, ctx)
}

// ValidateCreateCenterRequest implements IAiClientProvider.
func (a *aiClient) ValidateCreateCenterRequest(req ValidateCreateCenterRequest, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", create_center_request_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Center Image"), genai.ImageData("jpeg", req.CenterBytesImage))

	return a.processPrompt(prompt, ctx)
}

// ValidateProvideMealForChildTaskProof implements IAiClientProvider.
func (a *aiClient) ValidateProvideMealForChildTaskProof(req ValidateProvideMealForChildTaskProof, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", registration_request_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Proof of Task Image"), genai.ImageData("jpeg", req.ProofBytesImage))
	prompt = append(prompt, genai.Text("Label: Child Avatar Who Provided Lunch"), genai.ImageData("jpeg", req.ChildAvatarBytesImage))

	return a.processPrompt(prompt, ctx)
}

// ValidateRegistrationRequest implements IAiClientProvider.
func (a *aiClient) ValidateRegistrationRequest(req ValidateRegistrationRequest, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", registration_request_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Identity Card Image"), genai.ImageData("jpeg", req.IdentityCardBytesImage))
	prompt = append(prompt, genai.Text("Label: Avatar Image"), genai.ImageData("jpeg", req.AvatarBytesImage))

	return a.processPrompt(prompt, ctx)
}

// ValidateTaskProof implements IAiClientProvider.
func (a *aiClient) ValidateTaskProof(req ValidateTaskProof, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", registration_request_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Proof of Task Image"), genai.ImageData("jpeg", req.ProofBytesImage))

	return a.processPrompt(prompt, ctx)
}

// ValidateWithdrawProposal implements IAiClientProvider.
func (a *aiClient) ValidateWithdrawProposal(req ValidateWithdrawProposal, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s - %s", withdraw_proposal_validate_case, req.Purpose)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Withdraw Proof Image"), genai.ImageData("jpeg", req.ProofBytesImage))

	return a.processPrompt(prompt, ctx)
}

// ValidateChildSpecialNeedProposal implements IAiClientProvider.
func (a *aiClient) ValidateChildSpecialNeedProposal(req ValidateChildSpecialNeedProposal, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", child_special_need_proposal_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Campaign Relevant Proof Image"), genai.ImageData("jpeg", req.ProofBytesImage))

	return a.processPrompt(prompt, ctx)
}

// ValidatePoolCampaign implements IAiClientProvider.
func (a *aiClient) ValidatePoolCampaign(req ValidatePoolCampaign, ctx context.Context) string {
	if !a.geminiProvider.limiter.Allow() {
		return ""
	}

	jsonData, _ := json.MarshalIndent(req, "", "  ")
	var prompt = []genai.Part{
		genai.Text(fmt.Sprintf("Case: %s", pool_campaign_validate_case)),
		genai.Text(string(jsonData)),
	}

	prompt = append(prompt, genai.Text("Label: Pool Campaign Relevant Proof Image"), genai.ImageData("jpeg", req.ProofBytesImage))

	return a.processPrompt(prompt, ctx)
}

func (a *aiClient) processPrompt(prompt []genai.Part, ctx context.Context) string {
	res, err := a.geminiProvider.model.GenerateContent(ctx, prompt...)
	if err != nil {
		a.errLogger.Println(err.Error())
		return ""
	}
	if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
		var raw string = fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])
		return strings.ToLower(strings.Trim(raw, " \n\r\t.\"'"))
	}

	return "uncertain"
}
