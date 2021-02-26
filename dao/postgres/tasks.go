package postgres

import (
	"github.com/jinzhu/gorm"
	"oasisTracker/dmodels"
)

func (d DAO) GetTasks(OnlyActive bool) (tasks []dmodels.Task, err error) {

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

func (d DAO) GetLastTask() (task dmodels.Task, found bool, err error) {

	err = d.db.Select("*").
		Model(dmodels.Task{}).
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

func (d DAO) UpdateTask(task dmodels.Task) (err error) {
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

func (d DAO) CreateTask(task dmodels.Task) (err error) {
	return d.db.Model(dmodels.Task{}).Create(&task).Error
}
