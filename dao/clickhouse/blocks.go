package clickhouse

import (
	"fmt"
	"log"
	"oasisTracker/dmodels"
)

func (db DB) CreateBlocks(blocks []dmodels.Block) error {
	log.Print("Len: ", len(blocks))
	var err error
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_created_at, blk_hash, blk_proposer_address, blk_validator_hash, blk_epoch)"+
			"VALUES (?, ?, ?, ?, ?, ?)", dmodels.BlocksTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i := range blocks {

		if blocks[i].CreatedAt.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}

		_, err = stmt.Exec(
			blocks[i].Height,
			blocks[i].CreatedAt,
			blocks[i].Hash,
			blocks[i].ProposerAddress,
			blocks[i].ValidatorHash,
			blocks[i].Epoch,
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

func (db DB) CreateBlockSignatures(blocks []dmodels.BlockSignature) error {
	var err error

	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (sig_blk_lvl, sig_timestamp, sig_block_id_flag, sig_validator_address, sig_blk_signature)"+
			"VALUES (?, ?, ?, ?, ?)", dmodels.BlockSignaturesTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i := range blocks {

		if blocks[i].Timestamp.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}

		_, err = stmt.Exec(
			blocks[i].BlockHeight,
			blocks[i].Timestamp,
			blocks[i].BlockIDFlag,
			blocks[i].ValidatorAddress,
			blocks[i].Signature,
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
