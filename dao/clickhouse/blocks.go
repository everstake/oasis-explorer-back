package clickhouse

import (
	"fmt"
	"log"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"

	sq "github.com/wedancedalot/squirrel"
)

func (cl Clickhouse) CreateBlocks(blocks []dmodels.Block) (err error) {
	log.Print("Len blocks: ", len(blocks))

	if len(blocks) == 0 {
		return nil
	}

	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf(`INSERT INTO %s 
							(blk_lvl, 
							 blk_created_at, 
							 blk_hash, 
							 blk_proposer_address, 
							 blk_validator_hash, 
							 blk_epoch, 
							 blk_number_of_txs, 
							 blk_number_of_signatures, 
							 blk_fees, 
							 blk_gas_used) 
							 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, dmodels.BlocksNewTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i := range blocks {

		if blocks[i].CreatedAt.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}
		//log.Printf("Inserting block %v from epoch %v", blocks[i].Height, blocks[i].Epoch)
		_, err = stmt.Exec(
			blocks[i].Height,
			blocks[i].CreatedAt,
			blocks[i].Hash,
			blocks[i].ProposerAddress,
			blocks[i].ValidatorHash,
			blocks[i].Epoch,
			blocks[i].NumberOfTxs,
			blocks[i].NumberOfSignatures,
			blocks[i].Fees,
			blocks[i].GasUsed,
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
		fmt.Sprintf("INSERT INTO %s (blk_lvl, sig_timestamp, sig_block_id_flag, sig_validator_address, sig_blk_signature) VALUES (?, ?, ?, ?, ?)", dmodels.BlockSignaturesTable))
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

func (cl Clickhouse) GetBlocksList(params smodels.BlockParams) ([]dmodels.Block, error) {
	resp := make([]dmodels.Block, 0, params.Limit)

	s := (params.Limit * 7) + 86400
	if params.Offset != 0 {
		s += params.Offset * 7
	}

	q := sq.Select("*").
		From(dmodels.BlocksRowView).OrderBy("blk_lvl asc").
		JoinClause(fmt.Sprintf("ANY LEFT JOIN %s as sig USING blk_lvl", dmodels.BlocksSigCountView)).
		Limit(params.Limit).
		Offset(params.Offset)

	q = getBlocksQueryParam(q, params)

	rawSql, args, err := q.ToSql()
	if err != nil {
		return resp, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	row := dmodels.Block{}

	for rows.Next() {
		err := rows.Scan(
			&row.Height,
			&row.CreatedAt,
			&row.Hash,
			&row.ProposerAddress,
			&row.ValidatorHash,
			&row.Epoch,
			&row.NumberOfTxs,
			&row.GasUsed,
			&row.Fees,
			&row.NumberOfSignatures)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetBlocksListNew(params smodels.BlockParams) ([]dmodels.Block, error) {
	resp := make([]dmodels.Block, 0, params.Limit)

	s := (params.Limit * 7) + 86400
	if params.Offset != 0 {
		s += params.Offset * 7
	}

	q := sq.Select("*").
		From(dmodels.BlocksNewTable).OrderBy("blk_lvl desc").
		Limit(params.Limit).
		Offset(params.Offset)

	q = getBlocksQueryParam(q, params)

	rawSql, args, err := q.ToSql()
	if err != nil {
		return resp, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return resp, err
	}
	defer rows.Close()

	row := dmodels.Block{}

	for rows.Next() {
		err := rows.Scan(
			&row.Height,
			&row.CreatedAt,
			&row.Hash,
			&row.ProposerAddress,
			&row.ValidatorHash,
			&row.Epoch,
			&row.NumberOfTxs,
			&row.NumberOfSignatures,
			&row.Fees,
			&row.GasUsed)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetLastBlock() (block dmodels.Block, err error) {
	q := sq.Select("*").
		From(dmodels.BlocksNewTable).
		Where("blk_created_at >= now() - INTERVAL 1 DAY").
		Limit(1).
		OrderBy("blk_lvl desc")

	rawSql, args, err := q.ToSql()
	if err != nil {
		return block, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return block, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&block.Height,
			&block.CreatedAt,
			&block.Hash,
			&block.ProposerAddress,
			&block.ValidatorHash,
			&block.Epoch,
			&block.NumberOfTxs,
			&block.NumberOfSignatures,
			&block.Fees,
			&block.GasUsed)
		if err != nil {
			return block, err
		}
	}

	return block, nil
}

func (cl Clickhouse) BlocksCount(params smodels.BlockParams) (count uint64, err error) {
	//todo switch table to new
	q := sq.Select("count()").
		From(dmodels.BlocksOldTable)

	q = getBlocksQueryParam(q, params)

	rawSql, args, err := q.ToSql()
	if err != nil {
		return count, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return count, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (cl Clickhouse) BlockSignatures(params smodels.BlockParams) (count uint64, err error) {
	q := sq.Select("count()").
		From(dmodels.BlocksSigCountView)

	q = getBlocksQueryParam(q, params)

	rawSql, args, err := q.ToSql()
	if err != nil {
		return count, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return count, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func getBlocksQueryParam(q sq.SelectBuilder, params smodels.BlockParams) sq.SelectBuilder {

	if len(params.BlockLevel) > 0 {
		q = q.Where(sq.Eq{"blk_lvl": params.BlockLevel})
	}

	if len(params.BlockID) > 0 {
		q = q.Where(sq.Eq{"blk_hash": params.BlockID})
	}

	if len(params.Proposer) > 0 {
		q = q.Where(sq.Eq{"blk_proposer_address": params.Proposer})
	}

	if params.From > 0 {
		q = q.Where(sq.GtOrEq{"blk_created_at": params.From})
	}

	if params.To > 0 {
		q = q.Where(sq.Lt{"blk_created_at": params.To})
	}

	return q
}
