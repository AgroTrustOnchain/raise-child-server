package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
)

// type Transaction struct {
// 	ID           string    `json:"id"`
// 	ActorAddress string    `json:"actor_address"`
// 	ActionType   string    `json:"action_type"` // e.g., "Deposit", "Withdraw"
// 	Amount       int64     `json:"amount"`
// 	Message      string    `json:"message"`
// 	CoinType     string    `json:"coin_type"` // e.g., "Sui", "USDC"
// 	CreatedAt    time.Time `json:"created_at"`
// }

type Transaction struct {
	ID           ID     `json:"id"`
	ActorAddress string `json:"actor_address"`
	ActionType   string `json:"action_type"` // e.g., "Deposit", "Withdraw"
	PoolName     string `json:"pool_name"`
	Amount       int64  `json:"amount"`
	Message      string `json:"message"`
	CoinType     string `json:"coin_type"` // e.g., "Sui", "USDC"
	CreatedAt    int64  `json:"created_at"`
}

func (t Transaction) ToTransactionResponse() response.TransactionResponse {
	if t.ID.ID == "" {
		return response.TransactionResponse{}
	}

	return response.TransactionResponse{
		ID:           t.ID.ID,
		ActorAddress: t.ActorAddress,
		ActionType:   t.ActionType,
		Amount:       t.Amount,
		Message:      t.Message,
		CoinType:     t.CoinType,
		CreatedAt:    util.MilliSecToTime(t.CreatedAt),
	}
}
