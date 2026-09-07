package response

import "time"

type TransactionResponse struct {
	ID           string    `json:"id"`
	ActorAddress string    `json:"actor_address"`
	ActionType   string    `json:"action_type"` // e.g., "Deposit", "Withdraw"
	PoolName     string    `json:"pool_name"`
	Amount       int64     `json:"amount"`
	Message      string    `json:"message"`
	CoinType     string    `json:"coin_type"` // e.g., "VND"
	CreatedAt    time.Time `json:"created_at"`
}
