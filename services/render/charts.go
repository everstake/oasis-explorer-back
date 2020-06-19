package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func ChartData(chd []dmodels.ChartData) []smodels.ChartData {
	chds := make([]smodels.ChartData, len(chd))
	for i := range chd {
		chds[i] = ChartElement(chd[i])
	}
	return chds
}

func ChartElement(chd dmodels.ChartData) smodels.ChartData {

	return smodels.ChartData{
		Timestamp:         chd.BeginOfPeriod.Unix(),
		TransactionVolume: chd.TransactionVolume,
		EscrowRatio:       chd.EscrowRatio,
	}
}

func BalanceChartData(bcd []dmodels.BalanceChartData) []smodels.BalanceChartData {
	bcds := make([]smodels.BalanceChartData, len(bcd))
	for i := range bcd {
		bcds[i] = BalanceChartElement(bcd[i])
	}
	return bcds
}

func BalanceChartElement(bcd dmodels.BalanceChartData) smodels.BalanceChartData {

	return smodels.BalanceChartData{
		Timestamp:        bcd.BeginOfPeriod.Unix(),
		TotalBalance:     bcd.GeneralBalance,
		EscrowBalance:    bcd.EscrowBalance,
		DebondingBalance: bcd.DebondingBalance,
		SelfStakeBalance: 0,
	}
}
