package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	ID            uint      `gorm:"primarykey" json:"id"`
	StreamID      uint      `json:"streamId"`
	Stream        Stream    `gorm:"foreignKey:StreamID" json:"stream"`
	ScheduledTime time.Time `json:"scheduledTime"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
