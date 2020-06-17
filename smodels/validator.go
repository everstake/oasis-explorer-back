package smodels

func NewValidatorListParams() ValidatorParams {
	return ValidatorParams{
		CommonParams: CommonParams{
			Limit:  50,
			Offset: 0,
		},
	}
}

const (
	StatusActive   = "active"
	StatusInActive = "inactive"
)

type ValidatorParams struct {
	CommonParams
	ValidatorID string
}

type Validator struct {
	Account        string `json:"account_id"`
	AccountName    string `json:"account_name"`
	Fee            uint64 `json:"fee"`
	EscrowBalance  uint64 `json:"escrow_balance"`
	AvailableScore uint64 `json:"available_score"`
	CreatedAt      int64  `json:"validate_since"`
	ValidatorInfo
}

type ValidatorStats struct {
	Timestamp         int64  `json:"timestamp"`
	AvailabilityScore uint64 `json:"availability_score,omitempty"`
	BlocksCount       uint64 `json:"blocks_count,omitempty"`
	SignaturesCount   uint64 `json:"signatures_count,omitempty"`
}
