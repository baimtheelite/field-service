package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Field struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	UUID           string         `gorm:"type:uuid;not null"`
	Code           string         `gorm:"type:varchar(15);not null"`
	Name           string         `gorm:"type:varchar(100);not null"`
	PricePerHour   int            `gorm:"type:int;not null"`
	Images         pq.StringArray `gorm:"type:text[];not null"`
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	DeletedAt      *gorm.DeletedAt
	FieldSchedules []FieldSchedule `gorm:"foreignKey:field_id;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
