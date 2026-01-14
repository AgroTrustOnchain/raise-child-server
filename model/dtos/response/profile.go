package response

type PersonalProfileResponse struct {
	ID           string `json:"id"`
	IdentityCode string `json:"identity_code"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
}
