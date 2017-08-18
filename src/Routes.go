package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	name        string
	method      string
	pattern     string
	handlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/list", listGroups)
	router.HandleFunc("/list/{group}", listJobs)
	router.HandleFunc("/show/{group}/{name}", showJob)
	router.HandleFunc("/start/{group}/{name}", startJob)
	router.HandleFunc("/stop/{group}/{name}", stopJob)
	router.HandleFunc("/status/{group}", statusGroup)
	router.HandleFunc("/status/{group}/{name}", statusJob)
	return router
}

/*
	Route{
		"killPodStatus",
		"GET",
		"/killpod/status/{podName}",
		killPodStatus,
	},
	Route{
		"killPodStatuses",
		"GET",
		"/killpod/status/",
		killPodStatuses,
	},
	Route{
		"nodeExecStatus",
		"POST",
		"/nodeexec/status/",
		nodeExecStatus,
	},
	Route{
		"nodeExecStatusDetails",
		"POST",
		"/nodeexec/status/{id}",
		nodeExecStatusDetails,
	},
*/
