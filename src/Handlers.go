package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func sendCommands(ip string, data []RunnerExecData) ([]byte, error) {
	url := "http://" + ip + ":8081/cmd"

	json_data, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// list groups
func listGroups(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", formatArray(jobGroups))
}

// list jobs in a group
func listJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tests := listAvailableTests(vars["group"])
	fmt.Fprintf(w, "%v\n", formatArray(tests))
}

//show job contents
func showJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contents := loadFile(vars["group"], vars["name"])
	fmt.Fprintf(w, "%s\n", contents)
}

// start a job
func startJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["group"] == "killpod" {
		job := loadKillPod(vars["name"])
		startKillPod(vars["name"], job)
	} else if vars["group"] == "nodeexec" {
		job := loadNodeExec(vars["name"])
		startNodeExec(vars["name"], job)
	}
	fmt.Fprintln(w, "start initiated")
}

// stop a job
func stopJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["group"] == "killpod" {
		job := loadKillPod(vars["name"])
		stopKillPod(vars["name"], job)
		fmt.Fprintln(w, "stop initiated")
	} else if vars["group"] == "nodeexec" {
		fmt.Fprintln(w, "can not stop a nodeexec job")
	} else {
		fmt.Fprintln(w, "invalid group specified")
	}
}

// show the current status for a group
func statusGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["group"] == "killpod" {
		status := statusKillPod()
		fmt.Fprintln(w, status)
	} else if vars["group"] == "nodeexec" {
		status := statusNodeExec()
		fmt.Fprintln(w, status)
	} else {
		fmt.Fprintln(w, "invalid group specified")
	}
}

// show the current status for a job
func statusJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["group"] == "killpod" {
		status := statusKillPodJob(vars["name"])
		fmt.Fprintln(w, status)
	} else if vars["group"] == "nodeexec" {
		status := statusNodeExecJob(vars["name"])
		fmt.Fprintln(w, status)
	} else {
		fmt.Fprintln(w, "invalid group specified")
	}
}
