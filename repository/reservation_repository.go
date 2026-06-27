package repository

import (
	"errors"
	"spotsync-api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReservationRepository interface {
	CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error)
	FindByUserID(userID uint) ([]models.Reservation, error)
	FindByID(id uint) (*models.Reservation, error)
	Cancel(id uint) error
	FindAll() ([]models.Reservation, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}

// ⚠️ Critical — concurrency lock দিয়ে reservation বানাও
func (r *reservationRepository) CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Step 1 — Zone row lock করো (FOR UPDATE)
		var zone models.ParkingZone
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, zoneID).Error; err != nil {
			return errors.New("zone not found")
		}

		// Step 2 — Active reservation count করো
		var activeCount int64
		tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount)

		// Step 3 — Capacity check করো
		if activeCount >= int64(zone.TotalCapacity) {
			return errors.New("zone_full")
		}

		// Step 4 — Reservation বানাও
		reservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}

		return tx.Create(&reservation).Error
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) FindByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) Cancel(id uint) error {
	return r.db.Model(&models.Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled").Error
}

func (r *reservationRepository) FindAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Zone").Preload("User").
		Find(&reservations).Error
	return reservations, err
}
