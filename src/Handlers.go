package main

import (
  "encoding/json"
  "fmt"
  "io"
  "io/ioutil"
  "net/http"

  "github.com/gorilla/mux"
)

type KillPodData struct {
  Name      string  // name of pod base object
	Kind      string  // Kubernetes object to look for
  Target    int     // number of pods to kill at a time
  Interval  int     // time between kills
  Duration  int     // length of run
}

func KillPodStart(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  var data KillPodData
  data.Name = podName
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err)
  }
  if err := r.Body.Close(); err != nil {
    panic(err)
  }
  if err := json.Unmarshal(body, &data); err != nil {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(422) // unprocessable entity
    if err := json.NewEncoder(w).Encode(err); err != nil {
      panic(err)
    }
  }

  fmt.Fprintln(w, "pod name: ", podName)
  fmt.Fprintln(w, "data: ", data)

  startKillPod(data)
}

func KillPodStop(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  stopKillPod(podName)
}

func KillPodStatus(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  fmt.Fprintln(w, "pod name:", podName)
}

func KillPodStatuses(w http.ResponseWriter, r *http.Request) {
  statusesKillPod()
}
