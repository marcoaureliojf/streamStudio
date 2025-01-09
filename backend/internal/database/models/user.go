package models

import (
    "time"
    "gorm.io/gorm"
)
type User struct {
    gorm.Model
    ID        uint      `gorm:"primarykey" json:"id"`
    Name      string    `json:"name"`
    Email     string    `gorm:"unique" json:"email"`
    Password  string    `json:"password"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    TeamID    uint      `json:"teamId"`
    Team      Team      `gorm:"foreignKey:TeamID" json:"team"`
}