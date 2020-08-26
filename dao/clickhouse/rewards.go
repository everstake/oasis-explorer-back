package clickhouse

import (
	"fmt"
	sq "github.com/wedancedalot/squirrel"
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func (cl Clickhouse) CreateRewards(rewards []dmodels.Reward) error {
	tx, err := cl.db.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		fmt.Sprintf("INSERT INTO %s (blk_lvl, blk_epoch, reg_entity_address, rwd_amount, created_at)"+
			"VALUES (?, ?, ?, ?, ?, ?)", dmodels.RewardsTable))
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
			rewards[i].Amount,
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

func (cl Clickhouse) GetValidatorRewards(validatorID string, params smodels.CommonParams) (resp []dmodels.Reward, err error) {
	q := sq.Select("*").
		From(dmodels.RewardsTable).
		Where(sq.Eq{"reg_entity_address": validatorID}).
		OrderBy("blk_lvl desc")

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

		err = rows.Scan(&row.BlockLevel, &row.Epoch, &row.CreatedAt, &row.Amount, &row.EntityAddress)
		if err != nil {
			return resp, err
		}

		resp = append(resp, row)
	}

	return resp, nil
}

func (cl Clickhouse) GetValidatorRewardsStat(validatorID string) (resp dmodels.RewardsStat, err error) {
	q := sq.Select("*").
		From(dmodels.RewardsStatView).
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
