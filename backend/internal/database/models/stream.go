package models

import (
	"time"

	"gorm.io/gorm"
)

type Stream struct {
	gorm.Model
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	TeamID      uint      `json:"teamId"`
	Team        Team      `gorm:"foreignKey:TeamID" json:"team"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
