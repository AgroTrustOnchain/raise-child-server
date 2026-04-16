package entities

type MainPool struct {
	ID                                     ID       `json:"id"`
	LocalPools                             []string `json:"local_pools"`
	WithdrawProposals                      []string `json:"withdraw_proposals"`
	TotalAmount                            string   `json:"total_amount"`
	TotalMealDonationAmount                string   `json:"total_meal_donation_amount"`
	TotalBooksDonationAmount               string   `json:"total_books_donation_amount"`
	TotalHealthInsuranceNeedDonationAmount string   `json:"total_health_insurance_donation_amount"`
	Campaigns                              []string `json:"campaigns"`
	AllCampaigns                           []string `json:"all_campaigns"`
}

type LocalPool struct {
	ID                                     ID       `json:"id"`
	Region                                 string   `json:"region"`
	Mods                                   []string `json:"mods"`
	TotalAmount                            string   `json:"total_amount"`
	TotalMealDonationAmount                string   `json:"total_meal_donation_amount"`
	TotalBooksDonationAmount               string   `json:"total_books_donation_amount"`
	TotalHealthInsuranceNeedDonationAmount string   `json:"total_health_insurance_donation_amount"`
	Campaigns                              []string `json:"campaigns"`
}
