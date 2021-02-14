package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/labstack/echo/v4"
	"github.com/t-oki/pollen-api/internal/domain/entity"
)

var datetimeFormat = "2006010215"

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

type GetObservatoryInput struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (h *Handler) GetObservatory(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	param := GetObservatoryInput{}
	if err := c.Bind(&param); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	from, err := time.Parse(datetimeFormat, param.From)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	to, err := time.Parse(datetimeFormat, param.From)
	if err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}
	log.Info(fmt.Sprintf("id: %d, param: %v", id, param))
	if err := h.pollenRepo.FetchPollen("関東地域", id, from, to); err != nil {
		log.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusNotImplemented, nil)
}
