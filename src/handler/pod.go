package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/seungkyua/cookiemonster2/src/domain"
	"golang.org/x/net/context"
)

type PodHandler struct{}

var m = &domain.PodManage{
	Started: false,
}

/*
 * group.GET("", PodHandler.List)
 * group.GET("/:name", PodHandler.Get)
 * group.POST("/start", PodHandler.Start)
 * group.POST("/stop", PodHandler.Stop)
 */
func (h PodHandler) SetHandler(group *echo.Group) {
	group.GET("", h.List)
	group.POST("/start", h.Start)
	group.POST("/stop", h.Stop)
}

// List running pods
func (h PodHandler) List(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World\n")
}

// Start a job to delete random pod
func (h PodHandler) Start(c echo.Context) error {
	if m.Started {
		log.Println("Pod is already being munched, ignoring request\n")
		return c.String(http.StatusOK, "Pod is already being munched, ignoring request\n")
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		m.Ctx = ctx
		m.Cancel = cancel
		m.Started = true
	}

	err := m.Start(domain.GetConfig())
	if err != nil {
		log.Println(err)
		m.Started = false
		return c.String(http.StatusOK, "Cookie Monster Start Error !!! \n")
	}

	return c.String(http.StatusOK, "Cookie Monster Start!!! \n")
}

func (h PodHandler) Stop(c echo.Context) error {
	if !m.Started {
		log.Println("Pod is not currently getting munched, ignoring request\n")
		return c.String(http.StatusOK, "Pod is not currently getting munched, ignoring request\n")
	}

	m.Stop(domain.GetConfig())

	return c.String(http.StatusOK, "Cookie Monster Stop!!! \n")
}
