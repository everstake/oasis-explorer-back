package clickhouse

import (
	"oasisTracker/dmodels"

	sq "github.com/wedancedalot/squirrel"
)

func (cl Clickhouse) GetEntityActiveDepositorsCount(accountAddress string) (count uint64, err error) {

	q := sq.Select("depositors_num").
		From(dmodels.EntityActiveDepositorsCounterView).
		Where(sq.Eq{"reg_entity_address": accountAddress})

	rawSql, args, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	rows, err := cl.db.conn.Query(rawSql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}

	}

	return count, nil
}

func (cl Clickhouse) GetEntity(address string) (resp dmodels.EntityRegistryTransaction, err error) {
	q := sq.Select("*").
		From(dmodels.RegisterEntityTable).
		Where(sq.Eq{"reg_entity_address": address}).
		OrderBy("blk_lvl desc").
		Limit(1)

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
		err := rows.Scan(&resp.BlockLevel, &resp.Time, &resp.Hash, &resp.ID, &resp.Address, &resp.Nodes)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (cl Clickhouse) GetAccountValidatorInfo(accountAddress string) (resp dmodels.EntityNodesContainer, err error) {

	q := sq.Select("*").
		From(dmodels.EntityNodesView).
		Where(sq.Or{sq.Eq{"reg_entity_address": accountAddress}, sq.Eq{"reg_address": accountAddress}}).
		OrderBy("blk_lvl asc")

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
		row := dmodels.EntityNode{}

		err = rows.Scan(&row.EntityID, &row.EntityAddress, &row.NodeID, &row.Address, &row.ConsensusAddress, &row.CreatedTime, &row.LastRegBlock, &row.Expiration, &row.LastBlockTime, &row.BlocksCount, &row.LastSignatureTime, &row.BlocksSigned, &row.BlockSignaturesCount)
		if err != nil {
			return resp, err
		}

		resp.Nodes = append(resp.Nodes, row)
	}

	return resp, nil
}
