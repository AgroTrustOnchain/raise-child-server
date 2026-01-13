package response

type StaffNftResponse struct {
	ID           string `json:"id"`
	IdentityCode string `json:"identity_code"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Name         string `json:"name"`
	Url          string `json:"url"`
}
