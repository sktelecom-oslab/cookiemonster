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
		"listTestAll",
		"GET",
		"/list",
		listTestAll,
	},
	Route{
		"listTestGroup",
		"GET",
		"/list/{group}",
		listTestGroup,
	},
	Route{
		"showTest",
		"GET",
		"/show/{group}/{name}",
		showTest,
	},
	Route{
		"startAction",
		"GET",
		"/start/{group}/{name}",
		startAction,
	},
	Route{
		"stopAction",
		"GET",
		"/stop/{group}/{name}",
		stopAction,
	},
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
