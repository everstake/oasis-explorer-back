package postgres

import (
	"github.com/jinzhu/gorm"
	"oasisTracker/common/helpers"
	"oasisTracker/dmodels"
	"time"
)

const BlocksOffsetPostgresTable = "blocks_progress"

type BlocksOffset struct {
	ID     uint64 `gorm:"column:id"`
	Offset uint64 `gorm:"column:current_offset"`
}

func (d *Postgres) MigrateValidatorsInfo(validators []dmodels.ValidatorView) error {
	vi := new(dmodels.ValidatorInfo)
	vdi := new(dmodels.ValidatorDayInfo)

	err := d.db.Transaction(func(tx *gorm.DB) error {
		for i := range validators {
			if err := tx.Select("*").
				Table(dmodels.ValidatorsPostgresTable).
				Where("address = ?", validators[i].ConsensusAddress).
				Scan(&vi).Error; err != nil {
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
					vi.ID = lastId.ID + 1
					vi.Address = validators[i].ConsensusAddress
					vi.TotalBlocks = validators[i].ProposedBlocksCount
					vi.TotalSigs = validators[i].SignedBlocksCount
					vi.LastBlkTime = validators[i].LastBlockTime
					vi.LastSigTime = validators[i].LastBlockTime
					if err = tx.Table(dmodels.ValidatorsPostgresTable).
						Create(vi).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			if err := tx.Table(dmodels.ValidatorsPostgresTable).
				Where("id = ?", vi.ID).
				Updates(map[string]interface{}{
					"total_blk_count": validators[i].ProposedBlocksCount,
					"total_sig_count": validators[i].SignedBlocksCount,
					"last_blk_time":   validators[i].LastBlockTime,
					"last_sig_time":   validators[i].LastBlockTime,
				}).
				Error; err != nil {
				return err
			}

			if err := tx.Select("*").
				Table(dmodels.ValidatorsDayStatsPostgresTable).
				Where("val_id = ? and day = ?", vi.ID, helpers.TruncateToDay(time.Now())).
				Scan(&vdi).Error; err != nil {
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
					vdi.ID = lastId.ID + 1
					vdi.ValidatorID = vi.ID
					vdi.DayBlocks = validators[i].DayBlocksCount
					vdi.DaySigs = validators[i].DaySignaturesCount
					vdi.Day = helpers.TruncateToDay(time.Now())
					if err = tx.Table(dmodels.ValidatorsDayStatsPostgresTable).
						Create(vdi).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			if err := tx.Table(dmodels.ValidatorsDayStatsPostgresTable).
				Where("id = ?", vdi.ID).
				Updates(map[string]interface{}{
					"day_blk_count": validators[i].DayBlocksCount,
					"day_sig_count": validators[i].DaySignaturesCount,
				}).
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

func (d *Postgres) UpdateBlocksMigrationOffset(offset uint64) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(BlocksOffsetPostgresTable).
			Where("id = ?", 1).
			Update("current_offset", offset).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Postgres) GetBlocksMigrationOffset() (uint64, error) {
	s := new(BlocksOffset)

	if err := d.db.Table(BlocksOffsetPostgresTable).
		Where("id = 1").
		Select("*").
		First(&s).
		Error; err != nil {
		return 0, err
	}

	return s.Offset, nil
}
