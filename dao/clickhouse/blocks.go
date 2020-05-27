package clickhouse

import (
	"fmt"
	sq "github.com/wedancedalot/squirrel"
	"log"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func (cl Clickhouse) CreateBlocks(blocks []dmodels.Block) error {
	log.Print("Len: ", len(blocks))
	var err error
	tx, err := cl.db.conn.Begin()
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

func (cl Clickhouse) CreateBlockSignatures(blocks []dmodels.BlockSignature) error {
	var err error

	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, sig_timestamp, sig_block_id_flag, sig_validator_address, sig_blk_signature)"+
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

func (cl Clickhouse) GetBlocksList(params smodels.BlockParams) ([]dmodels.RowBlock, error) {

	resp := make([]dmodels.RowBlock, 0, params.Limit)

	q := sq.Select("*").
		From(dmodels.BlocksRowView).
		JoinClause(fmt.Sprintf("ANY LEFT JOIN %s as sig USING blk_lvl", dmodels.BlocksSigCountView)).
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
		row := dmodels.RowBlock{}

		err := rows.Scan(&row.Height, &row.CreatedAt, &row.Hash, &row.ProposerAddress, &row.ValidatorHash, &row.Epoch, &row.GasUsed, &row.Fee, &row.TxsCount, &row.SigCount)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
