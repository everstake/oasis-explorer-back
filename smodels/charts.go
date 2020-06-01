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
