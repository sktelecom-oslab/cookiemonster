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
	for _, route := range routes {
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(route.handlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"killPodStart",
		"POST",
		"/killpod/start/",
		killPodStart,
	},
	Route{
		"killPodStop",
		"POST",
		"/killpod/stop/",
		killPodStop,
	},
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
		"nodeExecStart",
		"POST",
		"/nodeexec/start/",
		nodeExecStart,
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
}
