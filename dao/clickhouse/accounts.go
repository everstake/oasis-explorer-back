package clickhouse

import (
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
)

func (cl Clickhouse) GetAccountTiming(accountID string) (resp dmodels.AccountTime, err error) {

	q := sq.Select("min(tx_time) created_at, max(tx_time) last_active").
		From(dmodels.TransactionsTable)

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
		err := rows.Scan(&resp.CreatedAt, &resp.LastActive)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}
