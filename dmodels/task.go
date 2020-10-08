package dmodels

const TasksTable = "tasks"

type Task struct {
	ID            uint64 `gorm:"column:tsk_id;PRIMARY_KEY;DEFAULT" `
	IsActive      bool   `gorm:"column:tsk_active;DEFAULT"`
	Title         string `gorm:"column:tsk_title"`
	StartHeight   uint64 `gorm:"column:tsk_start_height"`
	CurrentHeight uint64 `gorm:"column:tsk_current_height"`
	EndHeight     uint64 `gorm:"column:tsk_end_height"`
	Batch         uint64 `gorm:"column:tsk_batch"`
}
