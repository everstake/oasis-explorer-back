package smodels

type Reward struct {
	BlockLevel int64  `json:"block_level"`
	Epoch      uint64 `json:"epoch"`
	Amount     uint64 `json:"amount"`
	CreatedAt  int64  `json:"created_at"`
}

type RewardStat struct {
	EntityAddress string `json:"entity_address"`
	TotalAmount   uint64 `json:"total_reward"`
	DayAmount     uint64 `json:"day_reward"`
	WeekAmount    uint64 `json:"week_reward"`
	MonthAmount   uint64 `json:"month_reward"`
}
