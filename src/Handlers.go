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
  Name      string  // name of pod base object, empty for random
	Kind      string  // Kubernetes object to look for, if Name specified
  NameSpace string  // Namespace to use, otherwise will consider all but kube-system
  Target    int     // number of pods to kill at a time, defaults to 1
  Interval  int     // time between kills, unspecified for single kill
  Duration  int     // length of run, unspecified for single kill
}

// parse JSON from request body and return data struct
func readJSONData(r *http.Request) KillPodData {
  var data KillPodData
  body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
  if err != nil {
    panic(err)
  }
  if err := r.Body.Close(); err != nil {
    panic(err)
  }

  if err := json.Unmarshal(body, &data); err != nil {
    panic(err)
  }

  fmt.Println("request data: ", data)
  return data
}

func KillPodStart(w http.ResponseWriter, r *http.Request) {
  data := readJSONData(r)
  startKillPod(data)
}

func KillPodStop(w http.ResponseWriter, r *http.Request) {
  data := readJSONData(r)
  stopKillPod(data)
}

func KillPodStatus(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  fmt.Fprintln(w, "pod name: ", podName)
}

func KillPodStatuses(w http.ResponseWriter, r *http.Request) {
  statusesKillPod()
}
