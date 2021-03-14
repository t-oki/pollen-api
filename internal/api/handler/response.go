package handler

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/labstack/echo/v4"
)

type APIErrorResponse struct {
	Code   string   `json:"code"`
	Errors []string `json:"errors"`
}

func APIResponseOK(c echo.Context, data interface{}) error {
	return c.JSONPretty(http.StatusOK, data, "  ")
}

func APIResponseNoContent(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

func APIResponseError(c echo.Context, httpStatus int, message string, err error) error {
	if jsonErr := c.JSON(httpStatus, APIErrorResponse{Code: fmt.Sprintf("%d-000", httpStatus), Errors: []string{message}}); jsonErr != nil {
		log.Error(jsonErr.Error())
	}
	log.Error(err.Error())
	return err
}
