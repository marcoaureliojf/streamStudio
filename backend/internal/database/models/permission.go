    package models

    import (
        "time"
        "gorm.io/gorm"
    )
    type Permission struct {
        gorm.Model
        ID        uint      `gorm:"primarykey" json:"id"`
        Name      string    `gorm:"unique" json:"name"`
        CreatedAt time.Time `json:"createdAt"`
        UpdatedAt time.Time `json:"updatedAt"`
    }