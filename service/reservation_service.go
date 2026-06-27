package service

import (
	"errors"
	"spotsync-api/dto"
	"spotsync-api/repository"
	"time"
)

type ReservationService interface {
	Create(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	Cancel(reservationID, userID uint) error
	GetAll() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	zoneRepo        repository.ZoneRepository
}

func NewReservationService(
	reservationRepo repository.ReservationRepository,
	zoneRepo repository.ZoneRepository,
) ReservationService {
	return &reservationService{reservationRepo, zoneRepo}
}

func (s *reservationService) Create(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {

	_, err := s.zoneRepo.FindByID(req.ZoneID)
	if err != nil {
		return nil, errors.New("zone not found")
	}

	reservation, err := s.reservationRepo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		if err.Error() == "zone_full" {
			return nil, errors.New("zone_full")
		}
		return nil, errors.New("failed to create reservation")
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    reservation.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to fetch reservations")
	}

	var result []dto.MyReservationResponse
	for _, r := range reservations {
		result = append(result, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *reservationService) Cancel(reservationID, userID uint) error {
	// Reservation খোঁজো
	reservation, err := s.reservationRepo.FindByID(reservationID)
	if err != nil {
		return errors.New("reservation not found")
	}

	// নিজের reservation কিনা check করো
	if reservation.UserID != userID {
		return errors.New("forbidden")
	}

	// Already cancelled কিনা check করো
	if reservation.Status == "cancelled" {
		return errors.New("reservation already cancelled")
	}

	return s.reservationRepo.Cancel(reservationID)
}

func (s *reservationService) GetAll() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch reservations")
	}

	var result []dto.AdminReservationResponse
	for _, r := range reservations {
		result = append(result, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			User: dto.UserResponse{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			CreatedAt: r.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}
