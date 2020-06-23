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
	Account            string              `json:"account_id"`
	AccountName        string              `json:"account_name,omitempty"`
	NodeID             string              `json:"node_id"`
	Fee                uint64              `json:"fee"`
	EscrowBalance      uint64              `json:"escrow_balance"`
	EscrowBalanceShare uint64              `json:"escrow_shares"`
	GeneralBalance     uint64              `json:"general_balance"`
	DebondingBalance   uint64              `json:"debonding_balance"`
	AvailableScore     uint64              `json:"available_score"`
	CreatedAt          int64               `json:"validate_since"`
	MediaInfo          *ValidatorMediaInfo `json:"media_info"`
	ValidatorInfo
}

type ValidatorMediaInfo struct {
	WebsiteLink  string `json:"website_link,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	TwitterAcc   string `json:"twitter_acc,omitempty"`
	TGChat       string `json:"tg_chat,omitempty"`
	MediumLink   string `json:"medium_link,omitempty"`
}

type ValidatorStats struct {
	Timestamp         int64  `json:"timestamp"`
	AvailabilityScore uint64 `json:"availability_score,omitempty"`
	BlocksCount       uint64 `json:"blocks_count,omitempty"`
	SignaturesCount   uint64 `json:"signatures_count,omitempty"`
}

type Delegator struct {
	Account       string `json:"account_id"`
	EscrowAmount  uint64 `json:"escrow_amount"`
	DelegateSince int64  `json:"delegate_since"`
}
