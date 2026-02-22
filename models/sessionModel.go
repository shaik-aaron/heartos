package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	UserID          uint      `gorm:"not null;index" json:"userId"`
	User            User      `gorm:"foreignKey:UserID" json:"-"`
	CompletedAt     time.Time `json:"completedAt"`
	DurationMinutes int       `json:"durationMinutes"`
	DurationSeconds int       `json:"durationSeconds"`
	SessionType     string    `json:"sessionType"`
	Status          string    `json:"status"`
}
