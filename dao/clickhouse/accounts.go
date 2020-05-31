package clickhouse

import (
	"fmt"
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

func (cl Clickhouse) CreateAccountBalances(balances []dmodels.AccountBalance) (err error) {
	if len(balances) == 0 {
		return nil
	}

	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_time, acb_account, acb_nonce, acb_general_balance, acb_escrow_balance_active, acb_escrow_balance_share, acb_escrow_debonding_active, acb_escrow_debonding_share)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", dmodels.AccountBalanceTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i := range balances {

		if balances[i].Time.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}

		_, err = stmt.Exec(
			balances[i].Height,
			balances[i].Time,
			balances[i].Account,
			balances[i].Nonce,
			balances[i].GeneralBalance,
			balances[i].EscrowBalanceActive,
			balances[i].EscrowBalanceShare,
			balances[i].EscrowDebondingActive,
			balances[i].EscrowDebondingShare,
		)

		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
