package render

import (
	"encoding/json"
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
		Account:            a.EntityID,
		AccountName:        a.ValidatorName,
		NodeID:             a.NodeAddress,
		Fee:                a.ValidatorFee,
		GeneralBalance:     a.GeneralBalance,
		EscrowBalance:      a.EscrowBalance,
		EscrowBalanceShare: a.EscrowBalanceShare,
		DebondingBalance:   a.DebondingBalance,
		AvailableScore:     a.AvailabilityScore,
		CreatedAt:          a.ValidateSince.Unix(),
		MediaInfo:          ValidatorMediaInfo(a.ValidatorMediaInfo),
		ValidatorInfo: smodels.ValidatorInfo{
			Status:          a.Status,
			DepositorsCount: a.DepositorsNum,
			BlocksCount:     a.BlocksCount,
			SignaturesCount: a.SignaturesCount,
		},
	}
}

func ValidatorMediaInfo(validatorMediaInfoString string) *smodels.ValidatorMediaInfo {
	if validatorMediaInfoString == "" {
		return nil
	}

	var mediaInfo smodels.ValidatorMediaInfo
	if validatorMediaInfoString != "" {
		json.Unmarshal([]byte(validatorMediaInfoString), &mediaInfo)
	}

	return &mediaInfo
}

func ValidatorStatList(sts []dmodels.ValidatorStats) []smodels.ValidatorStats {
	stats := make([]smodels.ValidatorStats, len(sts))
	for i := range sts {
		stats[i] = ValidatorStat(sts[i])
	}
	return stats
}

func ValidatorStat(s dmodels.ValidatorStats) smodels.ValidatorStats {

	return smodels.ValidatorStats{
		Timestamp:         s.BeginOfPeriod.Unix(),
		AvailabilityScore: s.AvailabilityScore,
		BlocksCount:       s.BlocksCount,
		SignaturesCount:   s.SignaturesCount,
	}
}

func DelegatorList(sts []dmodels.Delegator) []smodels.Delegator {
	stats := make([]smodels.Delegator, len(sts))
	for i := range sts {
		stats[i] = Delegator(sts[i])
	}
	return stats
}

func Delegator(s dmodels.Delegator) smodels.Delegator {

	return smodels.Delegator{
		Account:       s.Address,
		EscrowAmount:  s.EscrowAmount,
		DelegateSince: s.DelegateSince.Unix(),
	}
}
