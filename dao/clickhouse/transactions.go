package clickhouse

import (
	"fmt"
	"log"
	"oasisTracker/dmodels"
)

func (db DB) CreateTransfers(transfers []dmodels.Transaction) error {
	log.Print("Len Transfers: ", len(transfers))
	var err error
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (tx_blk_lvl, tx_time, tx_hash, tx_amount, tx_escrow_amount,  tx_escrow_reclaim_amount, tx_escrow_account, tx_type, tx_sender, tx_receiver, tx_nonce, tx_fee, tx_gas_limit, tx_gas_price)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", dmodels.TransactionsTable))
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
