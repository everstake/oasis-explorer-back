package postgres

import (
	"github.com/jinzhu/gorm"
	"oasisTracker/common/helpers"
	"oasisTracker/dmodels"
	"time"
)

func (d *Postgres) SaveSignatures(signatures []dmodels.BlockSignature) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		vs := new(dmodels.ValidatorInfo)
		vds := new(dmodels.ValidatorDayInfo)
		for i := range signatures {
			if err := tx.Select("*").
				Table(dmodels.ValidatorsPostgresTable).
				Where("address = ?", signatures[i].ValidatorAddress).
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
					vs.Address = signatures[i].ValidatorAddress
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
					"total_sig_count": gorm.Expr("total_sig_count + 1"),
					"last_sig_time":   time.Now(),
				}).
				Error; err != nil {
				return err
			}

			if err := tx.Select("*").
				Table(dmodels.ValidatorsDayStatsPostgresTable).
				Where("val_id = ? and day = ?", vs.ID, helpers.TruncateToDay(time.Now())).
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
					vds.Day = helpers.TruncateToDay(time.Now())
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
				Update("day_sig_count", gorm.Expr("day_sig_count + 1")).
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
