package clickhouse

import (
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func (cl Clickhouse) GetChartsData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {

	q := sq.Select("start_of_period, toString(sum(tx_amount)) transaction_volume").
		From(dmodels.TransactionsTable).
		Where(sq.GtOrEq{"tx_time": params.From}).
		Where(sq.LtOrEq{"tx_time": params.To}).
		GroupBy("toStartOfDay(tx_time) as start_of_period").
		OrderBy("start_of_period asc")

	rawSql, args, err := q.ToSql()
	if err != nil {
		return resp, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		row := dmodels.ChartData{}

		err := rows.Scan(&row.BeginOfPeriod, &row.TransactionVolume)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetEscrowRatioChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {

	q := sq.Select("start_of_period, escrow_balance_active / (general_balance + escrow_balance_active + escrow_debonding_active) * 100 escrow_ratio").
		From(dmodels.DayTotalBalanceView).
		Where(sq.GtOrEq{"start_of_period": params.From}).
		Where(sq.LtOrEq{"start_of_period": params.To})

	rawSql, args, err := q.ToSql()
	if err != nil {
		return resp, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		row := dmodels.ChartData{}

		err := rows.Scan(&row.BeginOfPeriod, &row.EscrowRatio)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
