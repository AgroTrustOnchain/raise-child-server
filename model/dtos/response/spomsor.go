package response

type DonorResponse struct {
	ID              string                `json:"id"`
	Owner           string                `json:"owner"`
	FirstName       string                `json:"first_name"`
	LastName        string                `json:"last_name"`
	Gender          string                `json:"gender"`
	PhoneNumber     string                `json:"phone_number"`
	Email           string                `json:"email"`
	Name            string                `json:"name"`
	TotalDonation   int64                 `json:"total_donation"`
	SupportedChilds []string              `json:"supported_childs"`
	Url             string                `json:"url"`
	Contributions   []TransactionResponse `json:"contributions"`
}
