package handler

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/seungkyua/cookiemonster2/src/domain"
)

type ConfigHandler struct{}

/*
 * group.GET("", config.Get)
 */
func (h ConfigHandler) SetHandler(group *echo.Group) {
	group.GET("", h.Get)
}

// Get a config
func (h ConfigHandler) Get(context echo.Context) error {
	log.Println("###########", randomInt(1))
	return context.JSONPretty(http.StatusOK, domain.GetConfig(), "    ")
}

func randomInt(i int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(i)
}
