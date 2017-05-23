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
	Kind string
  Target int
  Interval int
  Duration int
}

func KillPodStart(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  var data KillPodData
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

/*
  t := RepoCreateTodo(todo)
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  if err := json.NewEncoder(w).Encode(t); err != nil {
      panic(err)
  }
*/

  fmt.Fprintln(w, "pod name: ", podName)
  fmt.Fprintln(w, "data: ", data)
}

func KillPodStop(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  fmt.Fprintln(w, "pod name:", podName)
}

func KillPodStatus(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  podName := vars["podName"]

  fmt.Fprintln(w, "pod name:", podName)
}
