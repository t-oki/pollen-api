package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/t-oki/pollen-api/internal/domain/entity"
)

type Handler struct {
	pollenRepo entity.PollenRepository
}

func NewHandler(pollenRepo entity.PollenRepository) *Handler {
	return &Handler{pollenRepo: pollenRepo}
}

func (h *Handler) ListAreas(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, nil)
}

func (h *Handler) GetArea(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, nil)
}

func (h *Handler) GetObservatory(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, nil)
}
