package smodels

func NewValidatorListParams() ValidatorParams {
	return ValidatorParams{
		CommonParams: CommonParams{
			Limit:  50,
			Offset: 0,
		},
	}
}

type ValidatorParams struct {
	CommonParams
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