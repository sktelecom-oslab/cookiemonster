package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/seungkyua/cookiemonster2/src/handler"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler.PodHandler{}.SetHandler(e.Group("/api/v1/pod"))
	handler.ConfigHandler{}.SetHandler(e.Group("/api/v1/config"))

	e.Logger.Debug(e.Start(":8080"))
}
