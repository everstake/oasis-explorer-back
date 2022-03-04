package clickhouse

import (
	"fmt"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"

	sq "github.com/wedancedalot/squirrel"
)

func (cl *Clickhouse) CreateRewards(rewards []dmodels.Reward) error {
	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_epoch, reg_entity_address, acb_account, rwd_amount, rwd_type, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", dmodels.RewardsTable))
	if err != nil {
		return err
	}

	defer stmt.Close()

	for i := range rewards {

		if rewards[i].CreatedAt.IsZero() {
			return fmt.Errorf("timestamp can not be 0")
		}

		_, err = stmt.Exec(
			rewards[i].BlockLevel,
			rewards[i].Epoch,
			rewards[i].EntityAddress,
			rewards[i].AccountAddress,
			rewards[i].Amount,
			rewards[i].Type,
			rewards[i].CreatedAt,
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

func (cl *Clickhouse) GetAccountRewards(accountID string, params smodels.CommonParams) (resp []dmodels.Reward, err error) {
	q := sq.Select("*").
		From(dmodels.RewardsTable).
		Where(sq.Eq{"acb_account": accountID}).
		OrderBy("blk_lvl desc").
		Limit(params.Limit).
		Offset(params.Offset)

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
		row := dmodels.Reward{}

		err = rows.Scan(&row.BlockLevel, &row.Epoch, &row.CreatedAt, &row.Amount, &row.Type, &row.EntityAddress, &row.AccountAddress)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl *Clickhouse) GetAccountRewardsStat(accountID string) (resp dmodels.RewardsStat, err error) {
	q := sq.Select("*").
		From(dmodels.AccountRewardsStatView).
		Where(sq.Eq{"acb_account": accountID})

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
		err = rows.Scan(&resp.AccountAddress, &resp.TotalAmount, &resp.DayAmount, &resp.WeekAmount, &resp.MonthAmount)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func (cl *Clickhouse) GetValidatorRewards(accountID string, params smodels.CommonParams) (resp []dmodels.Reward, err error) {
	q := sq.Select("reg_entity_address, blk_epoch, anyLast(blk_lvl), anyLast(created_at), sum(rwd_amount)").
		From(dmodels.RewardsTable).
		Where(sq.Eq{"reg_entity_address": accountID}).
		GroupBy("reg_entity_address, blk_epoch").
		OrderBy("blk_epoch desc").
		Limit(params.Limit).
		Offset(params.Offset)

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
		row := dmodels.Reward{}

		err = rows.Scan(&row.EntityAddress, &row.Epoch, &row.BlockLevel, &row.CreatedAt, &row.Amount)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl *Clickhouse) GetValidatorRewardsStat(validatorID string) (resp dmodels.RewardsStat, err error) {
	q := sq.Select("*").
		From(dmodels.ValidatorRewardsStatView).
		Where(sq.Eq{"reg_entity_address": validatorID})

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
		err = rows.Scan(&resp.EntityAddress, &resp.TotalAmount, &resp.DayAmount, &resp.WeekAmount, &resp.MonthAmount)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}
