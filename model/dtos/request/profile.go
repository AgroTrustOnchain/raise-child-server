package request

type UploadProfileRequest struct {
	IdentityCode string `json:"identity_code" validate:"required"`
	FirstName    string `json:"first_name" validate:"required"`
	LastName     string `json:"last_name" validate:"required"`
	Gender       string `json:"gender" validate:"required"`
	DateOfBirth  string `json:"date_of_birth" validate:"required"`
	PhoneNumber  string `json:"phone_number" validate:"required"`
	Email        string `json:"email" validate:"required"`
}
