package model

import (
	"time"
)

type Template struct {
	ID        int       `gorm:"column:id;type:int(11);primaryKey;autoIncrement" json:"id"`
	Type      string    `gorm:"column:type;size:255;"`
	Body      string    `gorm:"column:body;type:text;"`
	CreatedAt time.Time `gorm:"column:created_at;index;"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;"`
}
