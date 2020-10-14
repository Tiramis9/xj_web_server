package model

type SignReward struct {
	DayID       int     `json:"day_id"`
	RewardScore float64 `json:"reward_score"`
	Append      float64 `json:"append"`
}

type SignRewardContinuous struct {
	SeriesDay   int     `json:"series_day"`
	RewardScore float64 `json:"reward_score"`
}
