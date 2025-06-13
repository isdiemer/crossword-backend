package model

import "time"

type Session struct {
	Token     string    `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
