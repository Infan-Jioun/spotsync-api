package handler

import (
	"net/http"
	"spotsync-api/dto"
	"spotsync-api/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService service.ReservationService
	validate           *validator.Validate
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		validate:           validator.New(),
	}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	userID := c.Get("userID").(uint)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	reservation, err := h.reservationService.Create(userID, req)
	if err != nil {
		if err.Error() == "zone_full" {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false,
				Message: "Sorry, this zone is full",
			})
		}
		if err.Error() == "zone not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "Zone not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    reservation,
	})
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID := c.Get("userID").(uint)

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    reservations,
	})
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID := c.Get("userID").(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid reservation ID",
		})
	}

	if err := h.reservationService.Cancel(uint(id), userID); err != nil {
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Message: "You can only cancel your own reservations",
			})
		}
		if err.Error() == "reservation not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "Reservation not found",
			})
		}
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}

func (h *ReservationHandler) GetAll(c echo.Context) error {
	reservations, err := h.reservationService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "All reservations retrieved successfully",
		Data:    reservations,
	})
}
