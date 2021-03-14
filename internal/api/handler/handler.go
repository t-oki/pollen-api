package handler

import (
	"errors"
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
	Details []ListAreasOutputDetail
}
type ListAreasOutputDetail struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToListAreasOutput(areas []entity.Area) ListAreasOutput {
	res := make([]ListAreasOutputDetail, len(areas))
	for i, v := range areas {
		res[i] = ListAreasOutputDetail{
			ID:   v.ID,
			Name: v.Name,
		}
	}
	return ListAreasOutput{Details: res}
}

func (h *Handler) ListAreas(c echo.Context) error {
	return APIResponseOK(c, ToListAreasOutput(entity.ListAreas()))
}

type GetAreaOutput struct {
	ID            int64                      `json:"id"`
	Name          string                     `json:"name"`
	Observatories []GetAreaOutputObservatory `json:"observatories"`
}

type GetAreaOutputObservatory struct {
	ID         int64  `json:"id"`
	Prefecture string `json:"prefecture"`
	Name       string `json:"name"`
}

func ToGetAreaOutput(area entity.Area, observatories []entity.Observatory) GetAreaOutput {
	observatoriesRes := make([]GetAreaOutputObservatory, len(observatories))
	for i, v := range observatories {
		observatoriesRes[i] = GetAreaOutputObservatory{
			ID:         v.ID,
			Prefecture: v.Prefecture,
			Name:       v.Name,
		}
	}
	return GetAreaOutput{
		ID:            area.ID,
		Name:          area.Name,
		Observatories: observatoriesRes,
	}
}

func (h *Handler) GetArea(c echo.Context) error {
	areaID, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", err)
	}
	log.Info(fmt.Sprintf("ID: %v", areaID))
	area, err := entity.GetArea(areaID)
	if err != nil {
		var e *entity.NotFoundError
		if errors.As(err, &e) {
			return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", e)
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
	From string `json:"from" form:"from" query:"from"`
	To   string `json:"to" form:"to" query:"to"`
}

func (h *Handler) GetObservatory(c echo.Context) error {
	areaID, err := strconv.ParseInt(c.Param("area_id"), 0, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	observatoryID, err := strconv.ParseInt(c.Param("observatory_id"), 0, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	param := GetObservatoryInput{}
	if err := c.Bind(&param); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	log.Info(fmt.Sprintf("areaID: %d, observatoryID: %d, param: %v", areaID, observatoryID, param))

	from, err := time.Parse(datetimeFormat, param.From)
	if err != nil {
		return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", err)
	}
	to, err := time.Parse(datetimeFormat, param.To)
	if err != nil {
		return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", err)
	}
	if from.After(time.Now()) || from.After(to.Add(-1*time.Hour)) {
		return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", fmt.Errorf("from or to is wrong, from: %v, to: %v", from, to))
	}

	area, err := entity.GetArea(areaID)
	if err != nil {
		var e *entity.NotFoundError
		if errors.As(err, &e) {
			return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", e)
		}
		return APIResponseError(c, http.StatusInternalServerError, "サーバエラーが起きました", err)
	}
	observatory, err := entity.GetObservatory(area, observatoryID)
	if err != nil {
		var e *entity.NotFoundError
		if errors.As(err, &e) {
			return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", e)
		}
		return APIResponseError(c, http.StatusInternalServerError, "サーバエラーが起きました", err)
	}
	res, err := h.pollenRepo.FetchPollen(area, observatory, from, to)
	if err != nil {
		var e *entity.NotFoundError
		if errors.As(err, &e) {
			return APIResponseError(c, http.StatusBadRequest, "リクエストが不正です", e)
		}
		return APIResponseError(c, http.StatusInternalServerError, "サーバエラーが起きました", err)
	}
	return APIResponseOK(c, res)
}
