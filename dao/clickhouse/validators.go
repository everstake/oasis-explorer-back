package clickhouse

import (
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func (cl Clickhouse) GetValidatorsList(params smodels.ValidatorParams) (resp []dmodels.Validator, err error) {

	q := sq.Select("reg_entity_address,reg_consensus_address,node_address,created_time,start_blk_lvl,blocks,signatures, acb_escrow_balance_active, acb_general_balance,acb_escrow_balance_share,acb_escrow_debonding_active,depositors_num,is_active,pvl_name,pvl_fee,pvl_info").
		From(dmodels.ValidatorsTable).
		OrderBy("acb_escrow_balance_active desc").
		Limit(params.Limit).
		Offset(params.Offset)

	if params.ValidatorID != "" {
		q = q.Where(sq.Eq{"reg_entity_address": params.ValidatorID})
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
		row := dmodels.Validator{}

		err = rows.Scan(&row.EntityID, &row.ConsensusAddress, &row.NodeAddress, &row.ValidateSince, &row.StartBlockLevel, &row.BlocksCount, &row.SignaturesCount, &row.EscrowBalance, &row.GeneralBalance, &row.EscrowBalanceShare, &row.DebondingBalance, &row.DepositorsNum, &row.IsActive, &row.ValidatorName, &row.ValidatorFee, &row.ValidatorMediaInfo)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetValidatorDayStats(consensusAddress string, params smodels.ChartParams) (resp []dmodels.ValidatorStats, err error) {

	q := sq.Select("day,signatures,blocks,blk_lvl").
		From(dmodels.ValidatorStatsView).
		OrderBy("day asc").
		Where(sq.Eq{"reg_consensus_address": consensusAddress}).
		Where(sq.GtOrEq{"day": params.From}).
		Where(sq.LtOrEq{"day": params.To})

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
		row := dmodels.ValidatorStats{}

		err = rows.Scan(&row.BeginOfPeriod, &row.SignaturesCount, &row.BlocksCount, &row.LastBlock)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetValidatorDelegators(validatorID string, params smodels.CommonParams) (resp []dmodels.Delegator, err error) {

	q := sq.Select("tx_sender,escrow_since,balance").
		From(dmodels.DepositorsView).
		OrderBy("balance desc").
		Where(sq.Gt{"balance": 0}).
		Where(sq.Eq{"tx_receiver": validatorID})

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
		row := dmodels.Delegator{}

		err = rows.Scan(&row.Address, &row.DelegateSince, &row.EscrowAmount)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) PublicValidatorsSearchList() (resp []dmodels.Validator, err error) {
	q := sq.Select("reg_entity_address,pvl_name").
		From(dmodels.PublicValidatorsTable)

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
		row := dmodels.Validator{}

		err = rows.Scan(&row.EntityID, &row.ValidatorName)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
