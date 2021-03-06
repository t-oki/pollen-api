package handler

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/labstack/echo/v4"
)

const (
	prettyIndent = "  "
)

type APIErrorResponse struct {
	Code   string   `json:"code"`
	Errors []string `json:"errors"`
}

func APIResponseOK(c echo.Context, data interface{}) error {
	return c.JSONPretty(http.StatusOK, data, prettyIndent)
}

func APIResponseNoContent(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

//APIResponse APIResponseErrorとは異なり、errorは返さない. 4xxとかを想定
func APIResponse(c echo.Context, status int, message string) error {
	return APIResponseWithSubCode(c, status, 0, message)
}

func APIResponseWithSubCode(c echo.Context, status int, statusSubCode int, message string) error {
	return c.JSON(status, APIErrorResponse{Code: fmt.Sprintf("%d-%03d", status, statusSubCode), Errors: []string{message}})
}

func APIResponseCustomCode(c echo.Context, status int, customCode string, message string) error {
	return c.JSON(status, APIErrorResponse{Code: fmt.Sprintf("%d-%s", status, customCode), Errors: []string{message}})
}

func APIResponseError(c echo.Context, httpStatus int, message string, err error) error {
	if jsonErr := c.JSON(httpStatus, APIErrorResponse{Code: fmt.Sprintf("%d-000", httpStatus), Errors: []string{message}}); jsonErr != nil {
		log.Error(jsonErr.Error())
	}
	log.Error(err.Error())
	return err
}
