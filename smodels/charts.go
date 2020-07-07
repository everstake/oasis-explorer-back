package smodels

type ChartParams struct {
	From  uint64     `schema:"from"`
	To    uint64     `schema:"to"`
	Frame ChartFrame `schema:"frame"`
}

type ChartFrame string

const (
	FrameHour ChartFrame = "H"
	FrameDay  ChartFrame = "D"
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
	AccountsCount     uint64  `json:"accounts_count,omitempty"`
	AvgBlockTime      float64 `json:"avg_block_time,omitempty"`
	Fees              uint64  `json:"fees,omitempty"`
	OperationsCount   uint64  `json:"operations_count,omitempty"`
	ReclaimAmount     uint64  `json:"reclaim_amount,omitempty"`
}

type TopEscrowRatioChart struct {
	AccountID   string  `json:"account_id"`
	AccountName string  `json:"account_name,omitempty"`
	EsrowRatio  float64 `json:"esrow_ratio"`
}

type BalanceChartData struct {
	Timestamp        int64  `json:"timestamp"`
	GeneralBalance   uint64 `json:"general_balance"`
	EscrowBalance    uint64 `json:"escrow_balance"`
	DebondingBalance uint64 `json:"debonding_balance"`
	SelfStakeBalance uint64 `json:"self_stake_balance"`
}
