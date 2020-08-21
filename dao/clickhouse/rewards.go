package clickhouse

import (
	"fmt"
	"oasisTracker/dmodels"
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
