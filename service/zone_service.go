package service

import (
	"errors"
	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
	"time"

	"gorm.io/gorm"
)

type ZoneService interface {
	Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAll() ([]dto.ZoneResponse, error)
	GetByID(id uint) (*dto.ZoneResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo}
}

func (s *zoneService) Create(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(zone); err != nil {
		return nil, errors.New("failed to create zone")
	}

	return &dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *zoneService) GetAll() ([]dto.ZoneResponse, error) {
	zones, err := s.zoneRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch zones")
	}

	var result []dto.ZoneResponse
	for _, zone := range zones {
		// Available spots calculate করো
		activeCount, _ := s.zoneRepo.CountActiveReservations(zone.ID)
		available := zone.TotalCapacity - int(activeCount)

		result = append(result, dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: available,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *zoneService) GetByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.zoneRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("zone not found")
		}
		return nil, errors.New("failed to fetch zone")
	}

	activeCount, _ := s.zoneRepo.CountActiveReservations(zone.ID)
	available := zone.TotalCapacity - int(activeCount)

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt.Format(time.RFC3339),
	}, nil
}
