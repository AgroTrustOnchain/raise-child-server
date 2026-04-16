package ai

import "time"

// Validate cases
const (
	upload_child_validate_case                      string = "Upload Child Request Validation"
	create_center_request_validate_case             string = "Create Center Request Validation"
	registration_request_validate_case              string = "Registration Request Validation"
	task_proof_validate_case                        string = "Task Proof Validation"
	provide_meal_for_child_task_proof_validate_case string = "Provide Lunch Meal For Child Task Proof Validation"
	withdraw_proposal_validate_case                 string = "Withdraw Proposal Validation"
	child_special_need_proposal_validate_case       string = "Child Special Need Campaign Validation"
	pool_campaign_validate_case                     string = "Pool Campaign Validation"
)

type ValidateUploadChildRequest struct {
	HomeBytesImage                  []byte                `json:"home_bytes_image"`
	IdentityCode                    string                `json:"identity_code"`
	ChildBirthCertificateBytesImage []byte                `json:"child_birth_certificate_bytes_image"`
	Region                          string                `json:"region"`
	FirstName                       string                `json:"first_name"`
	LastName                        string                `json:"last_name"`
	Gender                          string                `json:"gender"`
	DateOfBirth                     string                `json:"date_of_birth"`
	HomeAddress                     string                `json:"home_address"`
	AvatarBytesImage                []byte                `json:"avatar_bytes_image"`
	FirstGuardian                   ChildGuardianProfile  `json:"first_guardian"`
	SecondGuardian                  *ChildGuardianProfile `json:"second_guardian"`
}

type ChildGuardianProfile struct {
	FullName               string `json:"guardian_full_name"`
	PhoneNumber            string `json:"guardian_phone_number"`
	Relation               string `json:"guardian_relation_with_child"`
	IdentityCardBytesImage []byte `json:"identity_card_bytes_image"`
}

type ValidateCreateCenterRequest struct {
	Region           string `json:"region"`
	Address          string `json:"address"`
	PhoneNumber      string `json:"phone_number"`
	CenterBytesImage []byte `json:"center_bytes_image"`
}

type ValidateRegistrationRequest struct {
	IdentityCode           string `json:"identity_code"`
	IdentityCardBytesImage []byte `json:"identity_card_bytes_image"`
	AvatarBytesImage       []byte `json:"avatar_bytes_image"`
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	Gender                 string `json:"gender"`
	DateOfBirth            string `json:"date_of_birth"`
	PhoneNumber            string `json:"phone_number"`
}

type ValidateTaskProof struct {
	TaskDescription string    `json:"task_description"`
	ProofBytesImage []byte    `json:"proof_bytes_image"`
	CreatedAt       time.Time `json:"created_at"`
}

type ValidateProvideMealForChildTaskProof struct {
	ChildAvatarBytesImage []byte `json:"child_avatar_bytes_image"`
	ValidateTaskProof
}

type ValidateWithdrawProposal struct {
	Purpose         string `json:"withdraw_purpose"`
	WithdrawAmount  int64  `json:"withdraw_amount"`
	Description     string `json:"description"`
	ProofBytesImage []byte `json:"proof_bytes_image"`
}

type ValidateChildSpecialNeedProposal struct {
	CamapaignTarget int64  `json:"campaign_target"`
	Description     string `json:"campaign_description"`
	ProofBytesImage []byte `json:"proof_bytes_image"`
}

type ValidatePoolCampaign struct {
	CamapaignTarget int64  `json:"campaign_target"`
	Description     string `json:"campaign_description"`
	ProofBytesImage []byte `json:"proof_bytes_image"`
}
