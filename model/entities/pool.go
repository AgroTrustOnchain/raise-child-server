package entities

type MainPool struct {
	ID                                     ID       `json:"id"`
	LocalPools                             []string `json:"local_pools"`
	ChildNeedWithdrawProposals             []string `json:"child_need_withdraw_proposals"`
	PendingWithdrawProposals               []string `json:"pending_withdraw_proposals"`
	ExecutedWithdrawProposals              []string `json:"executed_withdraw_proposals"`
	CancelledWithdrawProposals             []string `json:"cancelled_withdraw_proposals"`
	AllWithdrawProposals                   []string `json:"all_withdraw_proposals"`
	TotalAmount                            string   `json:"total_amount"`
	Donors                                 []string `json:"donors"`
	DonorTotalContributions                []string `json:"donor_total_contributions"`
	RemainVotePowers                       []string `json:"remain_vote_powers"`
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
	Donors                                 []string `json:"donors"`
	DonorTotalContributions                []string `json:"donor_total_contributions"`
	RemainVotePowers                       []string `json:"remain_vote_powers"`
	TotalMealDonationAmount                string   `json:"total_meal_donation_amount"`
	TotalBooksDonationAmount               string   `json:"total_books_donation_amount"`
	TotalHealthInsuranceNeedDonationAmount string   `json:"total_health_insurance_donation_amount"`
	Campaigns                              []string `json:"campaigns"`
}
