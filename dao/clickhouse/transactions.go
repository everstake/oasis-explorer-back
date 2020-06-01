package clickhouse

import (
	"fmt"
	sq "github.com/wedancedalot/squirrel"
	"log"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func (cl Clickhouse) CreateTransfers(transfers []dmodels.Transaction) error {
	if len(transfers) == 0 {
		return nil
	}

	log.Print("Len Transfers: ", len(transfers))
	var err error
	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_hash, tx_time, tx_hash, tx_amount, tx_escrow_amount,  tx_escrow_reclaim_amount, tx_escrow_account, tx_type, tx_sender, tx_receiver, tx_nonce, tx_fee, tx_gas_limit, tx_gas_price)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", dmodels.TransactionsTable))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range transfers {

		if transfers[i].Time.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}

		_, err = stmt.Exec(
			transfers[i].BlockLevel,
			transfers[i].BlockHash,
			transfers[i].Time,
			transfers[i].Hash,
			transfers[i].Amount,
			transfers[i].EscrowAmount,
			transfers[i].EscrowReclaimAmount,
			transfers[i].EscrowAccount,
			transfers[i].Type,
			transfers[i].Sender,
			transfers[i].Receiver,
			transfers[i].Nonce,
			transfers[i].Fee,
			transfers[i].GasLimit,
			transfers[i].GasPrice,
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

func (cl Clickhouse) CreateRegisterTransactions(txs []dmodels.RegistryTransaction) error {
	if len(txs) == 0 {
		return nil
	}

	var err error
	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, tx_time, tx_hash, reg_id, reg_entity_id, reg_entity_address, reg_expiration, reg_p2p_id, reg_consensus_id, reg_consensus_address, reg_physical_address, reg_roles)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", dmodels.RegisterTransactionsTable))
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range txs {

		if txs[i].Time.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}

		_, err = stmt.Exec(
			txs[i].BlockLevel,
			txs[i].Time,
			txs[i].Hash,
			txs[i].ID,
			txs[i].EntityID,
			txs[i].EntityAddress,
			txs[i].Expiration,
			txs[i].P2PID,
			txs[i].ConsensusID,
			txs[i].ConsensusAddress,
			txs[i].PhysicalAddress,
			txs[i].Roles,
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

func (cl Clickhouse) GetTransactionsList(params smodels.TransactionsParams) ([]dmodels.Transaction, error) {

	resp := make([]dmodels.Transaction, 0, params.Limit)

	q := sq.Select("*").
		From(dmodels.TransactionsTable).
		OrderBy("blk_lvl").
		Limit(params.Limit).
		Offset(params.Offset)

	if len(params.BlockLevel) > 0 {
		q = q.Where(sq.Eq{"blk_lvl": params.BlockLevel})
	}

	if len(params.BlockID) > 0 {
		q = q.Where(sq.Eq{"blk_hash": params.BlockID})
	}

	rawSql, args, err := q.ToSql()
	if err != nil {
		return resp, err
	}

	log.Print(rawSql)

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		row := dmodels.Transaction{}

		err := rows.Scan(&row.BlockLevel, &row.BlockHash, &row.Time, &row.Hash, &row.Amount, &row.EscrowAmount, &row.EscrowReclaimAmount, &row.EscrowAccount, &row.Type, &row.Sender, &row.Receiver, &row.Nonce, &row.Fee, &row.GasLimit, &row.GasPrice)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
