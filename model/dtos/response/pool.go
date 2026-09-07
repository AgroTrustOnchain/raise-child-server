package response

type PoolResponse struct {
	ID            string `json:"id"`
	PoolName      string `json:"pool_name"`
	TotalDonation int64  `json:"total_donation"`
}
