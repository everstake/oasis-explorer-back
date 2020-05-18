package dmodels

const TasksTable = "tasks"

type Task struct {
	ID            uint64 `db:"tsk_id"`
	IsActive      bool   `db:"tsk_active"`
	Title         string `db:"tsk_title"`
	StartHeight   uint64 `db:"tsk_start_height"`
	CurrentHeight uint64 `db:"tsk_current_height"`
	EndHeight     uint64 `db:"tsk_end_height"`
	Batch         uint64 `db:"tsk_batch"`
}
