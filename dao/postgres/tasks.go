package postgres

import (
	"oasisTracker/dmodels"

	"github.com/jinzhu/gorm"
)

func (d *Postgres) GetTasks(OnlyActive bool) (tasks []dmodels.Task, err error) {

	db := d.db.Select("*").
		Model(dmodels.Task{})
	if OnlyActive {
		db = db.Where("tsk_active = ?", OnlyActive)
	}

	err = db.Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (d *Postgres) GetLastTask(title string) (task dmodels.Task, found bool, err error) {

	err = d.db.Select("*").
		Model(dmodels.Task{}).
		Where("tsk_title = ?", title).
		Order("tsk_id desc").
		First(&task).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return task, false, nil
		}
		return task, false, err
	}

	return task, true, nil
}

func (d *Postgres) UpdateTask(task dmodels.Task) (err error) {
	err = d.db.Model(dmodels.Task{}).Where("tsk_id = ?", task.ID).
		Updates(map[string]interface{}{
			"tsk_current_height": task.CurrentHeight,
			"tsk_active":         task.IsActive,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Postgres) CreateTask(task dmodels.Task) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		oldTask := new(dmodels.Task)
		if err := tx.Select("*").
			Model(dmodels.Task{}).
			Where("tsk_title = ? and tsk_start_height = ? and tsk_end_height = ?", task.Title, task.StartHeight, task.EndHeight).
			Order("tsk_id desc").
			Last(&oldTask).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				return err
			}
		}

		if oldTask.ID != 0 {
			return nil
		}

		if err := tx.Model(dmodels.Task{}).
			Create(&task).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
