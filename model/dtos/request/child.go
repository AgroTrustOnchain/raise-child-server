package request

type GetChildrenRequest struct {
	Keyword     string `form:"keyword"`
	Region      string `form:"region"`
	YearOfBirth *int   `form:"year_of_birth"`
	SortOrder   string `form:"sort_order"`
	Gender      string `form:"gender"`
	PageSize    int    `form:"page_size"`
	Page        int    `form:"page"`
}

type UploadChildRequest struct {
	HomeBlobID     string                `json:"home_blob_id" validate:"required"`
	IdentityCode   string                `json:"identity_code" validate:"required"`
	Region         string                `json:"region" validate:"required"`
	FirstName      string                `json:"first_name" validate:"required"`
	LastName       string                `json:"last_name" validate:"required"`
	Gender         string                `json:"gender" validate:"required"`
	DateOfBirth    string                `json:"date_of_birth" validate:"required"`
	HomeAddress    string                `json:"home_address" validate:"required"`
	AvatarBlobId   string                `json:"avatar_blob_id" validate:"required"`
	FirstGuardian  ChildGuardianProfile  `json:"first_guardian" validate:"required"`
	SecondGuardian *ChildGuardianProfile `json:"second_guardian"`
}

type ChildGuardianProfile struct {
	FullName           string `json:"guardian_full_name" validate:"required"`
	PhoneNumber        string `json:"guardian_phone_number" validate:"required"`
	Relation           string `json:"guardian_relation" validate:"required"`
	IdentityCardBlobID string `json:"identity_card_blob_id" validate:"required"`
}

type AddChildStringMetadaRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type AddChildNumberMetadaRequest struct {
	Key   string `json:"key" validate:"required"`
	Value int    `json:"value" validate:"required"`
}

// childs/books-need/{id}/support

// childs/meal-need/{id}/support
type SupportMealNeadRequest struct {
	Months int `json:"months" validate:"required,min=1,max=12"`
}

// childs/special-need/{id}/support
type SupportSpecialNeedRequest struct {
	Amount      int64  `json:"amount" validate:"required,min=2000"`
	Description string `json:"description" validate:"required"`
}

type ConfirmProvideMealForChildRequest struct {
	ImageBlobID string `json:"image_blob_id" validate:"required"`
}

// childs/books-need/withdraw-proposal
// childs/meal-need/withdraw-proposal
type CreateNormalNeedWithdrawProposalRequest struct {
	NeedID      string  `json:"need_id" validate:"required"`
	ProofBlobID *string `json:"proof_blob_id"`
}

// childs/special-need/withdraw-proposal
type CreateSpecialNeedWithdrawProposalRequest struct {
	CampaignID  string  `json:"campaign_id" validate:"required"`
	Amount      int64   `json:"amount" validate:"required,min=2000"`
	Description string  `json:"description" validate:"required"`
	ProofBlobID *string `json:"proof_blob_id"`
}

// childs/special-need/proposal
type CreateSpecialNeedProposalRequest struct { // For create special need campaign
	ChildID     string  `json:"child_id" validate:"required"`
	Target      int64   `json:"target" validate:"required,min=100000"`
	Description string  `json:"description" validate:"required"`
	ProofBlobID *string `json:"proof_blob_id"`
}

type UpdateChildNeedRequest struct {
	ChildID string `json:"child_id" validate:"required"`
	NeedID  string `json:"need_id" validate:"required"`
	Value   *int64 `json:"value"`
}

type UpdateChildEditNeedDatesRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// childs/special-need/proposal/{id}/confirm
