package main

import (
  "net/http"
  "github.com/gorilla/mux"
)

type Route struct {
  Name string
  Method string
  Pattern string
  HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

  router := mux.NewRouter().StrictSlash(true)
  for _, route := range routes {
    router.
      Methods(route.Method).
      Path(route.Pattern).
      Name(route.Name).
      Handler(route.HandlerFunc)
  }

  return router
}

var routes = Routes {
  Route{
    "KillPodStart",
    "POST",
    "/killpod/start/{podName}",
    KillPodStart,
  },
  Route{
    "KillPodStop",
    "POST",
    "/killpod/stop/{podName}",
    KillPodStop,
  },
  Route{
    "KillPodStatus",
    "GET",
    "/killpod/status/{podName}",
    KillPodStatus,
  },
  Route{
    "KillPodStatuses",
    "GET",
    "/killpod/status",
    KillPodStatuses,
  },
}
