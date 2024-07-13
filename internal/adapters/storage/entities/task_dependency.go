package entities

type TaskDependency struct {
	TaskID          uint `gorm:"primaryKey"`
	DependentTaskID uint `gorm:"primaryKey"`
}
