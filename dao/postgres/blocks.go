package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"oasisTracker/common/helpers"
	"oasisTracker/dmodels"
	"time"
)

func (d *Postgres) SaveBlocks(blocks []dmodels.Block) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		b := new(dmodels.BlockInfo)
		bd := new(dmodels.BlockDayInfo)

		vs := new(dmodels.ValidatorInfo)
		vds := new(dmodels.ValidatorDayInfo)

		if len(blocks) > 0 {
			if err := tx.Select("*").
				Table(dmodels.BlocksPostgresTable).
				Order("id desc").
				Limit(1).
				Scan(&b).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {
					b.ID = 1
					b.TotalBlocks = 0
					b.LastLvl = 0
					b.Epoch = 0
					if err = tx.Table(dmodels.BlocksPostgresTable).
						Create(b).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			bInfo := map[string]interface{}{
				"total_count": gorm.Expr(fmt.Sprintf("total_count + %d", len(blocks))),
			}

			if b.LastLvl < blocks[len(blocks)-1].Height {
				bInfo["last_lvl"] = blocks[len(blocks)-1].Height
				bInfo["epoch"] = blocks[len(blocks)-1].Epoch
			}

			if err := tx.Table(dmodels.BlocksPostgresTable).
				Where("id = ?", b.ID).
				Updates(bInfo).
				Error; err != nil {
				return err
			}
		}

		for i := range blocks {
			if err := tx.Select("*").
				Table(dmodels.BlocksDayPostgresTable).
				Where("day = ?", helpers.TruncateToDay(blocks[i].CreatedAt)).
				Scan(&bd).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {
					lastId := new(dmodels.BlockDayInfo)
					if err = tx.Table(dmodels.BlocksDayPostgresTable).
						Select("id").
						Order("id desc").
						First(&lastId).Error; err != nil {
						if !gorm.IsRecordNotFoundError(err) {
							return err
						}
					}
					bd.ID = lastId.ID + 1
					bd.TotalDayBlocks = 0
					bd.Day = helpers.TruncateToDay(blocks[i].CreatedAt)
					if err = tx.Table(dmodels.BlocksDayPostgresTable).
						Create(bd).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			if err := tx.Table(dmodels.BlocksDayPostgresTable).
				Where("id = ?", bd.ID).
				Update("day_total_count", gorm.Expr(fmt.Sprintf("day_total_count + %d", 1))).
				Error; err != nil {
				return err
			}

			if err := tx.Select("*").
				Table(dmodels.ValidatorsPostgresTable).
				Where("address = ?", blocks[i].ProposerAddress).
				Scan(&vs).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {
					lastId := new(dmodels.ValidatorInfo)
					if err = tx.Table(dmodels.ValidatorsPostgresTable).
						Select("id").
						Order("id desc").
						First(&lastId).Error; err != nil {
						if !gorm.IsRecordNotFoundError(err) {
							return err
						}
					}
					vs.ID = lastId.ID + 1
					vs.Address = blocks[i].ProposerAddress
					vs.TotalBlocks = 0
					vs.TotalSigs = 0
					vs.LastBlkTime = time.Time{}
					vs.LastSigTime = time.Time{}
					if err = tx.Table(dmodels.ValidatorsPostgresTable).
						Create(vs).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			if err := tx.Table(dmodels.ValidatorsPostgresTable).
				Where("id = ?", vs.ID).
				Updates(map[string]interface{}{
					"total_blk_count": gorm.Expr("total_blk_count + 1"),
					"last_blk_time":   blocks[i].CreatedAt,
				}).
				Error; err != nil {
				return err
			}

			if err := tx.Select("*").
				Table(dmodels.ValidatorsDayStatsPostgresTable).
				Where("val_id = ? and day = ?", vs.ID, helpers.TruncateToDay(blocks[i].CreatedAt)).
				Scan(&vds).Error; err != nil {
				if gorm.IsRecordNotFoundError(err) {
					lastId := new(dmodels.ValidatorDayInfo)
					if err = tx.Table(dmodels.ValidatorsDayStatsPostgresTable).
						Select("id").
						Order("id desc").
						First(&lastId).Error; err != nil {
						if !gorm.IsRecordNotFoundError(err) {
							return err
						}
					}
					vds.ID = lastId.ID + 1
					vds.ValidatorID = vs.ID
					vds.DayBlocks = 0
					vds.DaySigs = 0
					vds.Day = helpers.TruncateToDay(blocks[i].CreatedAt)
					if err = tx.Table(dmodels.ValidatorsDayStatsPostgresTable).
						Create(vds).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			if err := tx.Table(dmodels.ValidatorsDayStatsPostgresTable).
				Where("id = ?", vds.ID).
				Update("day_blk_count", gorm.Expr("day_blk_count + 1")).
				Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Postgres) GetBlocksInfo() (*dmodels.BlockInfo, error) {
	bi := new(dmodels.BlockInfo)

	if err := d.db.Table(dmodels.BlocksPostgresTable).
		Select("*").
		First(&bi).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return bi, nil
}

func (d *Postgres) GetBlocksDayInfo() (*dmodels.BlockDayInfo, error) {
	bdi := new(dmodels.BlockDayInfo)

	if err := d.db.Table(dmodels.BlocksDayPostgresTable).
		Select("*").
		Where("day = ?", helpers.TruncateToDay(time.Now())).
		First(&bdi).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return bdi, nil
}

func (d *Postgres) SaveTotalBlocksCount(count uint64) error {
	if err := d.db.Table(dmodels.BlocksPostgresTable).
		Where("id <> 0").
		Update("total_count", count).
		Error; err != nil {
		return err
	}

	return nil
}
