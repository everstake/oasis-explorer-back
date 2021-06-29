package mysql

import (
	"oasisTracker/common/dao"
	"oasisTracker/dmodels"

	"github.com/wedancedalot/squirrel"
)

func (md MysqlDAO) GetTasks(OnlyActive bool) (tasks []dmodels.Task, err error) {
	q := squirrel.Select("*").From(dmodels.TasksTable)
	if OnlyActive {
		q = q.Where(squirrel.Eq{"tsk_active": OnlyActive})
	}
	err = md.mysql.Find(&tasks, q)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (md MysqlDAO) GetLastTask(title string) (task dmodels.Task, found bool, err error) {
	q := squirrel.Select("*").From(dmodels.TasksTable).Where("tsk_title = ?", title).OrderBy("tsk_id desc")

	err = md.mysql.FindFirst(&task, q)
	if err != nil {
		if err == dao.ErrNoRows {
			return task, false, nil
		}
		return task, false, err
	}

	return task, true, nil
}

func (md MysqlDAO) UpdateTask(task dmodels.Task) (err error) {
	q := squirrel.Update(dmodels.TasksTable).
		Where(squirrel.Eq{"tsk_id": task.ID}).
		SetMap(map[string]interface{}{
			"tsk_current_height": task.CurrentHeight,
			"tsk_active":         task.IsActive,
		})

	_, err = md.mysql.Exec(q.ToSql())
	if err != nil {
		return err
	}
	return nil
}

func (md MysqlDAO) CreateTask(task dmodels.Task) (err error) {
	q := squirrel.Insert(dmodels.TasksTable).SetMap(squirrel.Eq{
		"tsk_active":         task.IsActive,
		"tsk_title":          task.Title,
		"tsk_start_height":   task.StartHeight,
		"tsk_current_height": task.CurrentHeight,
		"tsk_end_height":     task.EndHeight,
		"tsk_batch":          task.Batch,
	})

	_, err = md.mysql.Exec(q.ToSql())
	if err != nil {
		return err
	}
	return nil
}
