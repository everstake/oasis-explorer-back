package dmodels

import "time"

const RewardsTable = "rewards"

type Reward struct {
	EntityAddress string    `db="reg_entity_address"`
	BlockLevel    int64     `db="blk_lvl"`
	Epoch         uint64    `db="blk_epoch"`
	Amount        uint64    `db="rwd_amount"`
	CreatedAt     time.Time `db:"created_at"`
}
