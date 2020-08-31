package clickhouse

import (
	"fmt"
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

func (cl Clickhouse) GetTotalAccountsCountChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {

	q := sq.Select("start_of_period, accounts_number").
		From(dmodels.DayAccountsCountView).
		Where(sq.GtOrEq{"start_of_period": params.From}).
		Where(sq.LtOrEq{"start_of_period": params.To}).
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

		err := rows.Scan(&row.BeginOfPeriod, &row.AccountNumber)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetAvgBlockTimeChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {

	q := sq.Select("start_of_period, avg(next_blk_created_at - blk_created_at) avg_delay").
		From(dmodels.BlocksTable).
		JoinClause("ANY LEFT JOIN (select toUInt64(blk_lvl - 1) blk_lvl, blk_created_at next_blk_created_at from blocks) s using blk_lvl").
		Where(sq.Gt{"next_blk_created_at": 0}).
		Where(sq.GtOrEq{"start_of_period": params.From}).
		Where(sq.LtOrEq{"start_of_period": params.To}).
		GroupBy(fmt.Sprintf("%s(next_blk_created_at) as start_of_period", getFrameFunc(params.Frame))).
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

		err := rows.Scan(&row.BeginOfPeriod, &row.AvgBlockTime)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetFeeVolumeChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {

	q := sq.Select("start_of_period, sum(tx_fee) tx_fee").
		From(dmodels.TransactionsTable).
		Where(sq.GtOrEq{"tx_time": params.From}).
		Where(sq.LtOrEq{"tx_time": params.To}).
		GroupBy(fmt.Sprintf("%s(tx_time) as start_of_period", getFrameFunc(params.Frame))).
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

		err := rows.Scan(&row.BeginOfPeriod, &row.Fees)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetOperationsCountChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {
	q := sq.Select("start_of_period, count() operations").
		From(dmodels.TransactionsTable).
		Where(sq.GtOrEq{"start_of_period": params.From}).
		Where(sq.LtOrEq{"start_of_period": params.To}).
		GroupBy(fmt.Sprintf("%s(tx_time) as start_of_period", getFrameFunc(params.Frame))).
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

		err := rows.Scan(&row.BeginOfPeriod, &row.OperationsCount)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetReclaimAmountChartData(params smodels.ChartParams) (resp []dmodels.ChartData, err error) {
	q := sq.Select("start_of_period, sum(tx_escrow_reclaim_amount) reclaim_amount").
		From(dmodels.TransactionsTable).
		Where(sq.GtOrEq{"start_of_period": params.From}).
		Where(sq.LtOrEq{"start_of_period": params.To}).
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

		err := rows.Scan(&row.BeginOfPeriod, &row.ReclaimAmount)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
