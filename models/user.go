package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"default:driver"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
