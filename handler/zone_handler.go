package handler

import (
	"net/http"
	"spotsync-api/dto"
	"spotsync-api/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
	validate    *validator.Validate
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{
		zoneService: zoneService,
		validate:    validator.New(),
	}
}

func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest

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

	zone, err := h.zoneService.Create(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    zone,
	})
}

func (h *ZoneHandler) GetAll(c echo.Context) error {
	zones, err := h.zoneService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    zones,
	})
}

func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid zone ID",
		})
	}

	zone, err := h.zoneService.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    zone,
	})
}
