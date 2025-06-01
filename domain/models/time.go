package models

import "time"

type Time struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UUID      string `gorm:"type:uuid;not null"`
	StartTime string `gorm:"type:time without timezone;not null"`
	EndTime   string `gorm:"type:time without timezone;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
