package postgres

import (
	"github.com/jinzhu/gorm"
	"oasisTracker/dmodels"
)

func (d *Postgres) GetValidatorsInfo() ([]dmodels.ValidatorInfoWithDay, error) {
	vi := make([]dmodels.ValidatorInfoWithDay, 0)

	if err := d.db.Raw(`select v.id,
								   v.address,
								   v.total_blk_count,
								   v.total_sig_count,
								   v.last_blk_time,
								   v.last_sig_time,
								   vds.day_blk_count,
								   vds.day_sig_count,
								   vds.day
							from validators v
							inner join validator_day_stats vds on v.id = vds.val_id
							where vds.day >= date_trunc('day', now())`).
		Scan(&vi).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return vi, nil
}
