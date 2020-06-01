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

	data := smodels.ChartData{
		Timestamp:         chd.BeginOfPeriod.Unix(),
		TransactionVolume: chd.TransactionVolume,
		EscrowRatio:       chd.EscrowRatio,
	}

	return data
}
