package models

import "time"

type Reservation struct {
	ID           uint        `gorm:"primaryKey;autoIncrement"`
	UserID       uint        `gorm:"not null"`
	ZoneID       uint        `gorm:"not null"`
	LicensePlate string      `gorm:"not null;size:15"`
	Status       string      `gorm:"default:active"`
	User         User        `gorm:"foreignKey:UserID"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
