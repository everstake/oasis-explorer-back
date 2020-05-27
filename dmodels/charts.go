package dmodels

import "time"

type ChartData struct {
	BeginOfPeriod     time.Time `db:"start_of_period"`
	TransactionVolume string    `db:"transaction_volume"`
}
