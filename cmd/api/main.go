package main

import (
	"flag"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/t-oki/pollen-api/internal/infra/hanako"
	"github.com/t-oki/pollen-api/internal/infra/handler"
)

var port = flag.String("port", "8080", "port to listen")

func main() {
	flag.Parse()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())

	//cors
	aos := []string{"*"}
	if os.Getenv("ALLOW_ORIGINS") != "" {
		aos = strings.Split(os.Getenv("ALLOW_ORIGINS"), ",")
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  aos,
		ExposeHeaders: []string{echo.HeaderXRequestID},
	}))

	pollenRepo := hanako.NewPollenRepositoryImpl()
	handler := handler.NewHandler(pollenRepo)
	e.GET("/areas", handler.ListAreas)
	e.GET("/areas/:id", handler.GetArea)
	e.GET("/observatories/:id", handler.GetObservatory)

	e.Start(":" + *port)
}
