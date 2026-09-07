package response

type GetTransactionRecordsResponse struct {
	PoolName          string                `json:"pool_name"`
	PoolID            string                `json:"pool_id"`
	PoolTotalDonation int64                 `json:"pool_total_donation"`
	Data              []TransactionResponse `json:"data"`
	Amount            int                   `json:"amount"`
	Page              int                   `json:"page"`
	TotalPages        int                   `json:"total_pages"`
}
