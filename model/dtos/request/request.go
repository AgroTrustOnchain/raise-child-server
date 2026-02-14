package request

type VoteRequest struct {
	IsVoteYes    bool   `json:"is_vote_yes" validate:"required"`
	RefuseReason string `json:"refuse_reason"`
}

// type CreateRegistrationRequest struct {
// 	RegisterRole       string `json:"register_role"`
// 	IdentityCode       string `json:"identity_code"`
// 	IdentityCardBlobID string `json:"identity_card_blob_id"`
// 	AvatarBlobID       string `json:"avatar_blob_id"`
// 	Region             string `json:"region"`
// 	FirstName          string `json:"first_name"`
// 	LastName           string `json:"last_name"`
// 	Gender             string `json:"gender"`
// 	DateOfBirth        string `json:"date_of_birth"`
// 	PhoneNumber        string `json:"phone_number"`
// 	Email              string `json:"email"`
// }

// Registraion requests

type AdminRegistrationRequest struct {
	IdentityCardBlobID string `json:"identity_card_blob_id" validate:"required"`
	AvatarBlobID       string `json:"avatar_blob_id" validate:"required"`
}

type VolunteerRegistrationRequest struct {
	Region string `json:"region" validate:"required"`
	AdminRegistrationRequest
}

type LocalLeaderRegistrationRequest struct {
	CenterAddress     string `json:"center_address" validate:"required"`
	CenterPhoneNumber string `json:"center_phone_number" validate:"required"`
	CenterImageBlobID string `json:"center_image_blob_id" validate:"required"`
	VolunteerRegistrationRequest
}

// For Admin role
type GetAdminRegistrationRequets struct {
	Keyword   string `json:"keyword"`
	Gender    string `json:"gender"`
	Status    string `json:"status"`
	IsClosed  *bool  `json:"is_closed"`
	IsConfirm *bool  `json:"is_confirm"`
	SortOrder string `json:"sort_order"`
	PageSize  int    `json:"page_size"`
	Page      int    `json:"page"`
}

// For Volunteer, Local Leader role
type GetNormalStaffRegistrationRequests struct {
	Region string `json:"region"`
	GetAdminRegistrationRequets
}

type CreateRegistrationRequest struct {
	RegisterRole       string `json:"register_role" validate:"required"`
	Region             string `json:"region" validate:"required"`
	IdentityCardBlobID string `json:"identity_card_blob_id" validate:"required"`
	AvatarBlobID       string `json:"avatar_blob_id" validate:"required"`
}

type GetRegistrationRequests struct {
	RegisterRole         string `json:"register_role"`
	IsAvailableToConfirm *bool  `json:"is_available_to_confirm"`
	GetUploadChildRequests
}

type RegistrationRoleRequest struct {
	RegisterRole       string `json:"register_role" validate:"required"`
	IdentityCode       string `json:"identity_code" validate:"required"`
	IdentityCardBlobID string `json:"identity_card_blob_id" validate:"required"`
	AvatarBlobID       string `json:"avatar_blob_id" validate:"required"`
	Region             string `json:"region" validate:"required"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	Gender             string `json:"gender" validate:"required"`
	DateOfBirth        string `json:"date_of_birth" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	Email              string `json:"email" validate:"required"`
}

// Upload-child requests
type GetUploadChildRequests struct {
	Keyword   string `json:"keyword"`
	Region    string `json:"region"`
	Gender    string `json:"gender"`
	Status    string `json:"status"`
	IsClosed  *bool  `json:"is_closed"`
	SortOrder string `json:"sort_order"`
	PageSize  int    `json:"page_size"`
	Page      int    `json:"page"`
}

// Publisher
type UpdatePublisherInfoRequest struct {
	IdentityCode       string `json:"identity_code" validate:"required"`
	IdentityCardBlobID string `json:"identity_card_blob_id" validate:"required"`
	AvatarBlobID       string `json:"avatar_blob_id" validate:"required"`
	FirstName          string `json:"first_name" validate:"required"`
	LastName           string `json:"last_name" validate:"required"`
	Gender             string `json:"gender" validate:"required"`
	DateOfBirth        string `json:"date_of_birth" validate:"required"`
	PhoneNumber        string `json:"phone_number" validate:"required"`
	Email              string `json:"email" validate:"required"`
}

type WithdrawRequest struct {
	Description string `json:"description" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,min=10000"`
}

type AdoptRequest struct {
	ChildID     string `json:"child_id" validate:"required"`
	Description string `json:"description" validate:"required"`
	Duration    int    `json:"duration" validate:"required,min=1"`
	Unit        string `json:"unit" validate:"required"` // e.g. "month", "year"
}
