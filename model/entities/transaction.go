package entities

import (
	"raise-child/model/dtos/response"
	"raise-child/util"
	"strconv"
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
	Amount       string `json:"amount"`
	Message      string `json:"message"`
	CoinType     string `json:"coin_type"` // e.g., "Sui", "USDC"
	CreatedAt    string `json:"created_at"`
}

func (t Transaction) ToTransactionResponse() response.TransactionResponse {
	if t.ID.ID == "" {
		return response.TransactionResponse{}
	}

	amount, _ := strconv.ParseInt(t.Amount, 10, 64)
	createdAt, _ := strconv.ParseInt(t.CreatedAt, 10, 64)

	return response.TransactionResponse{
		ID:           t.ID.ID,
		ActorAddress: t.ActorAddress,
		ActionType:   t.ActionType,
		PoolName:     t.PoolName,
		Amount:       amount,
		Message:      t.Message,
		CoinType:     t.CoinType,
		CreatedAt:    util.MilliSecToTime(createdAt),
	}
}
