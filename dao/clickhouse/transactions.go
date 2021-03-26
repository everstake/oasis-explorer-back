package clickhouse

import (
	"fmt"
	"log"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"

	"github.com/ClickHouse/clickhouse-go"
	sq "github.com/wedancedalot/squirrel"
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
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_hash, tx_time, tx_hash, tx_amount, tx_escrow_amount,  tx_escrow_reclaim_amount, tx_type, tx_status, tx_error, tx_sender, tx_receiver, tx_nonce, tx_fee, tx_gas_limit, tx_gas_price)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", dmodels.TransactionsTable))
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
			transfers[i].Type,
			transfers[i].Status,
			transfers[i].Error,
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

func (cl Clickhouse) CreateRegisterNodeTransactions(txs []dmodels.NodeRegistryTransaction) error {
	if len(txs) == 0 {
		return nil
	}

	var err error
	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, tx_time, tx_hash, reg_id, reg_address, reg_entity_id, reg_entity_address, reg_expiration, reg_p2p_id, reg_consensus_id, reg_consensus_address, reg_physical_address, reg_roles)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", dmodels.RegisterNodeTable))
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
			txs[i].Address,
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

func (cl Clickhouse) CreateRegisterEntityTransactions(txs []dmodels.EntityRegistryTransaction) error {
	if len(txs) == 0 {
		return nil
	}

	var err error
	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, tx_time, tx_hash, reg_entity_id, reg_entity_address, reg_nodes)"+
			"VALUES (?, ?, ?, ?, ?, ?)", dmodels.RegisterEntityTable))
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
			txs[i].Address,
			clickhouse.Array(txs[i].Nodes),
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
		OrderBy("blk_lvl desc").
		Limit(params.Limit).
		Offset(params.Offset)

	if len(params.OperationID) > 0 {
		q = q.Where(sq.Eq{"tx_hash": params.OperationID})
	}

	if len(params.OperationKind) > 0 {
		q = q.Where(sq.Eq{"tx_type": params.OperationKind})
	}

	if len(params.BlockLevel) > 0 {
		q = q.Where(sq.Eq{"blk_lvl": params.BlockLevel})
	}

	if len(params.BlockID) > 0 {
		q = q.Where(sq.Eq{"blk_hash": params.BlockID})
	}

	if params.From > 0 {
		q = q.Where(sq.GtOrEq{"tx_time": params.From})
	}

	if params.To > 0 {
		q = q.Where(sq.Lt{"tx_time": params.To})
	}

	if params.AccountID != "" {
		q = q.Where(sq.Or{sq.Eq{"tx_sender": params.AccountID}, sq.Eq{"tx_receiver": params.AccountID}})
	}

	if params.Sender != "" {
		q = q.Where(sq.Eq{"tx_sender": params.Sender})
	}

	if params.Receiver != "" {
		q = q.Where(sq.Eq{"tx_receiver": params.Receiver})
	}

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
		row := dmodels.Transaction{}

		err = rows.Scan(&row.BlockLevel, &row.BlockHash, &row.Time, &row.Hash, &row.Amount, &row.EscrowAmount, &row.EscrowReclaimAmount, &row.Type, &row.Status, &row.Error, &row.Sender, &row.Receiver, &row.Nonce, &row.Fee, &row.GasLimit, &row.GasPrice)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
