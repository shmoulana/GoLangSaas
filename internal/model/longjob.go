package model

import "time"

type LongJob struct {
	ID         int        `gorm:"column:id;type:int(11);primaryKey;autoIncrement" json:"id"`
	EntityId   int        `gorm:"column:entity_id;size:11"`
	Entity     string     `gorm:"column:entity;size:255"`
	TenantId   *string    `gorm:"column:tenant_id;size:255"`
	Meta       *string    `gorm:"column:meta;size:255"`
	Status     string     `gorm:"column:status;size:255"`
	RetryCount int        `gorm:"column:retry_count;size:255"`
	FailedAt   *time.Time `gorm:"column:failed_at;size:255"`
	FinishedAt *time.Time `gorm:"column:finished_at;size:255"`
	CreatedAt  time.Time  `gorm:"column:created_at;"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;"`
}
