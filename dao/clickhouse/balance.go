package clickhouse

import (
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
)

func (cl Clickhouse) GetLastDayTotalBalance() (bal dmodels.DayBalance, err error) {
	q := sq.Select("*").
		From(dmodels.DayTotalBalanceView).
		Limit(1)

	rawSql, args, err := q.ToSql()
	if err != nil {
		return bal, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return bal, err
	}

	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&bal.StartOfPeriod, &bal.GeneralBalance, &bal.EscrowBalanceActive, &bal.EscrowDebondingActive)
		if err != nil {
			return bal, err
		}

	}

	return bal, nil
}
