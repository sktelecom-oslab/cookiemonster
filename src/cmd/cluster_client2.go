package main

import (
	"fmt"
	"time"

	"github.com/seungkyua/cookiemonster2/src/domain"
	"log"
)

func main() {
	config := domain.GetConfig()
	err := config.ReadConfig("../../config/")
	if err != nil {
		fmt.Println(err)
	}

	pm := &domain.PodManage{
		Started: false,
	}

	startTime := time.Now().Add(time.Duration(config.Duration) * time.Second)
	ticker := time.NewTicker(time.Second * time.Duration(config.Interval))

	err = pm.MainLoop(config)
	if err != nil {
		panic(err.Error())
	}

	for range ticker.C {
		if startTime.Before(time.Now()) {
			log.Println("Duration time out !!!")
			return
		}
		err = pm.MainLoop(config)
		if err != nil {
			panic(err.Error())
		}
	}
}
