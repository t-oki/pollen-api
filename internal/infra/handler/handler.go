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

type ListAreasOutput struct {
	details []ListAreasOutputDetail
}
type ListAreasOutputDetail struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToListAreasOutput(areas []entity.Area) ListAreasOutput {
	res := make([]ListAreasOutputDetail, len(areas))
	for _, v := range areas {
		res = append(res, ListAreasOutputDetail{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return ListAreasOutput{details: res}
}

func (h *Handler) ListAreas(c echo.Context) error {
	return APIResponseOK(c, ToListAreasOutput(entity.ListAreas()))
}

type GetAreaInput struct {
	ID int64 `path:"id"`
}

type GetAreaOutput struct {
	ID            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Observatories []GetAreaOutputObservatory `json:"observatories"`
}

type GetAreaOutputObservatory struct {
	ID         int64  `json:"id"`
	Prefecture string `json:"prefecture"`
	Name       string `json:"name`
}

func ToGetAreaOutput(area entity.Area, observatories []entity.Observatory) GetAreaOutput {
	observatoriesRes := make([]GetAreaOutputObservatory, len(observatories))
	for _, v := range observatories {
		observatoriesRes = append(observatoriesRes, GetAreaOutputObservatory{
			ID:         v.ID,
			Prefecture: v.Prefecture,
			Name:       v.Name,
		})
	}
	return GetAreaOutput{
		ID:            area.ID,
		Name:          area.Name,
		Observatories: observatoriesRes,
	}
}

func (h *Handler) GetArea(c echo.Context) error {
	param := GetAreaInput{}
	c.Bind(&param)
	area, err := entity.GetArea(param.ID)
	if err != nil {
		if err == entity.ErrAreaNotExist {
			return APIResponse(c, http.StatusBadRequest, "リクエストが不正です")
		}
		return APIResponseError(c, http.StatusInternalServerError, "サーバエラーが起きました", err)
	}
	observatories, err := area.ListObservatories()
	if err != nil {
		return APIResponseError(c, http.StatusInternalServerError, "サーバエラーが起きました", err)
	}
	return APIResponseOK(c, ToGetAreaOutput(area, observatories))
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
