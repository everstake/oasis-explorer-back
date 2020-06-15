package dmodels

import "time"

const ValidatorsTable = "validators_list_view"

type Validator struct {
	EntityID          string    `db:"reg_entity_id"`
	NodeAddress       string    `db:"reg_id"`
	ValidateSince     time.Time `db:"created_time"`
	LastBlockTime     time.Time `db:"last_block_time"`
	BlocksCount       uint64    `db:"blocks"`
	LastSignatureTime uint64    `db:"last_signature_time"`
	SignaturesCount   uint64    `db:"signatures"`
	EscrowBalance     uint64    `db:"acb_escrow_balance_active"`
	DepositorsNum     uint64    `db:"depositors_num"`
	IsActive          bool      `db:"is_active"`
	ValidatorName     string    `db:"pvl_name"`
	ValidatorFee      uint64    `db:"pvl_fee"`
	WebAddress        string    `db:"pvl_address"`
	AvailableScore    uint64    `db:"-"`
	Status            string    `db:"-"`
}
