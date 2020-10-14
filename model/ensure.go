package model

type Ensure struct {
	ID int `json:"id"`
	LossScore float64 `json:"loss_score"`
	RewardScore float64 `json:"reward_score"`
	BalanceLimit float64 `json:"balance_limit"`
}
