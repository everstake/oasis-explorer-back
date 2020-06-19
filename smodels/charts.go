package smodels

type ChartParams struct {
	From  uint64     `schema:"from"`
	To    uint64     `schema:"to"`
	Frame ChartFrame `schema:"frame"`
}

type ChartFrame string

const (
	FrameDay ChartFrame = "D"
)

func (p ChartParams) Validate() error {
	if p.Frame == "" {
		p.Frame = FrameDay
	}

	return nil
}

type ChartData struct {
	Timestamp         int64   `json:"timestamp"`
	TransactionVolume string  `json:"transaction_volume,omitempty"`
	EscrowRatio       float64 `json:"escrow_ratio,omitempty"`
}

type BalanceChartData struct {
	Timestamp        int64  `json:"timestamp"`
	TotalBalance     uint64 `json:"total_balance"`
	EscrowBalance    uint64 `json:"escrow_balance"`
	DebondingBalance uint64 `json:"debonding_balance"`
	SelfStakeBalance uint64 `json:"self_stake_balance"`
}
