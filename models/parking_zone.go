package models

import "time"

type ParkingZone struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"not null"`
	Type          string  `gorm:"not null"`
	TotalCapacity int     `gorm:"not null"`
	PricePerHour  float64 `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
