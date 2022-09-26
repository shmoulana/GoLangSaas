package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int            `gorm:"column:id;type:int(11);primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"column:name;size:255"`
	Email     string         `gorm:"column:email;size:255"`
	Password  string         `gorm:"column:password;type:text"`
	CreatedAt time.Time      `gorm:"column:created_at;index;"`
	UpdatedAt time.Time      `gorm:"column:updated_at;index;"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;"`
}
