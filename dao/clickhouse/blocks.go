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
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_created_at, blk_hash, blk_proposer_address, blk_validator_hash, blk_epoch) VALUES (?, ?, ?, ?, ?, ?)", dmodels.BlocksTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i := range blocks {

		if blocks[i].CreatedAt.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}
		log.Printf("Inserting block %v from epoch %v", blocks[i].Height, blocks[i].Epoch)
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

func (cl Clickhouse) GetBlocksList(params smodels.BlockParams) ([]dmodels.RowBlock, error) {

	resp := make([]dmodels.RowBlock, 0, params.Limit)

	s := (params.Limit * 12) + 900
	if params.Offset != 0 {
		s += params.Offset * 12
	}

	q := sq.Select("*").
		From(dmodels.BlocksRowView).OrderBy("blk_lvl DESC").
		JoinClause(fmt.Sprintf("ANY LEFT JOIN %s as sig USING blk_lvl", dmodels.BlocksSigCountView)).
		Where(fmt.Sprintf("blk_created_at >= now() - INTERVAL %d SECOND", s)).
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

	row := dmodels.RowBlock{}

	for rows.Next() {
		err := rows.Scan(&row.Height, &row.CreatedAt, &row.Hash, &row.ProposerAddress, &row.ValidatorHash, &row.Epoch, &row.GasUsed, &row.Fee, &row.TxsCount, &row.SigCount)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetLastBlock() (block dmodels.Block, err error) {
	q := sq.Select("*").
		From(dmodels.BlocksTable).
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
		err = rows.Scan(&block.Height, &block.CreatedAt, &block.Hash, &block.ProposerAddress, &block.ValidatorHash, &block.Epoch)
		if err != nil {
			return block, err
		}
	}

	return block, nil
}

func (cl Clickhouse) BlocksCount(params smodels.BlockParams) (count uint64, err error) {
	q := sq.Select("count()").
		From(dmodels.BlocksTable)

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
