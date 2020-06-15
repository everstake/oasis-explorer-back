package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func ValidatorsList(accs []dmodels.Validator) []smodels.Validator {
	accounts := make([]smodels.Validator, len(accs))
	for i := range accs {
		accounts[i] = Validator(accs[i])
	}
	return accounts
}

func Validator(a dmodels.Validator) smodels.Validator {

	return smodels.Validator{
		Account:        a.EntityID,
		AccountName:    a.ValidatorName,
		Fee:            a.ValidatorFee,
		EscrowBalance:  a.EscrowBalance,
		AvailableScore: a.AvailableScore,
		CreatedAt:      a.ValidateSince.Unix(),
		ValidatorInfo: smodels.ValidatorInfo{
			CommissionScheduleRules: smodels.TestNetGenesis,
			Status:                  a.Status,
			NodeAddress:             a.NodeAddress,
			DepositorsCount:         a.DepositorsNum,
			BlocksCount:             a.BlocksCount,
			SignaturesCount:         a.SignaturesCount,
		},
	}
}
