package clickhouse

import (
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
)

func (cl Clickhouse) GetEntityActiveDepositorsCount(accountID string) (count uint64, err error) {

	q := sq.Select("depositors_num").
		From(dmodels.EntityActiveDepositorsCounterView).
		Where(sq.Eq{"reg_entity_id": accountID})

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

func (cl Clickhouse) GetAccountValidatorInfo(accountID string) (resp dmodels.EntityNodesContainer, err error) {

	q := sq.Select("*").
		From(dmodels.EntityNodesView).
		Where(sq.Or{sq.Eq{"reg_entity_id": accountID}, sq.Eq{"reg_id": accountID}}).
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

		err := rows.Scan(&row.EntityID, &row.NodeID, &row.ConsensusAddress, &row.CreatedTime, &row.LastRegBlock, &row.Expiration, &row.LastBlockTime, &row.BlocksCount, &row.LastSignatureTime, &row.BlockSignaturesCount)
		if err != nil {
			return resp, err
		}

		resp.Nodes = append(resp.Nodes, row)
	}

	return resp, nil
}
