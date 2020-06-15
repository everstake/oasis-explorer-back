package clickhouse

import (
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
)

func (cl Clickhouse) GetValidatorsList(params interface{}) (resp []dmodels.Validator, err error) {

	q := sq.Select("reg_entity_id,reg_id,created_time,blocks,signatures,acb_escrow_balance_active,depositors_num,is_active,pvl_name,pvl_fee,pvl_address").
		From(dmodels.ValidatorsTable).
		OrderBy("acb_escrow_balance_active desc")

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

		err := rows.Scan(&row.EntityID, &row.NodeAddress, &row.ValidateSince, &row.BlocksCount, &row.SignaturesCount, &row.EscrowBalance, &row.DepositorsNum, &row.IsActive, &row.ValidatorName, &row.ValidatorFee, &row.WebAddress)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}
