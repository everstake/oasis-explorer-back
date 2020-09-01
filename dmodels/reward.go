package dmodels

import "time"

const (
	RewardsTable    = "rewards"
	RewardsStatView = "validator_rewards_stat_view"
)

type Reward struct {
	EntityAddress string    `db="reg_entity_address"`
	BlockLevel    int64     `db="blk_lvl"`
	Epoch         uint64    `db="blk_epoch"`
	Amount        uint64    `db="rwd_amount"`
	CreatedAt     time.Time `db:"created_at"`
}

type RewardsStat struct {
	EntityAddress string `db="reg_entity_address"`
	TotalAmount   uint64 `db="total_amount"`
	DayAmount     uint64 `db="day_amount"`
	WeekAmount    uint64 `db="week_amount"`
	MonthAmount   uint64 `db="month_amount"`
}
