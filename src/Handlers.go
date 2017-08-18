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

// list all available tests
func listTestAll(w http.ResponseWriter, r *http.Request) {
	tests := listAvailableTests("")
	fmt.Fprintf(w, "available tests:\n%v", formatArray(tests))
}

func listTestGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tests := listAvailableTests(vars["group"])
	fmt.Fprintf(w, "available tests:\n%v", formatArray(tests))
}

func showTest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contents := loadFile(vars["group"], vars["name"])
	fmt.Fprintf(w, "test contents:\n%s", contents)
}

func startAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["group"] == "killpod" {
		test := loadKillPod(vars["name"])
		startKillPod(test)
	} else if vars["group"] == "nodeexec" {
		test := loadNodeExec(vars["name"])
		startNodeExec(test)
	}
	fmt.Fprintln(w, "test start triggered")
}

func stopAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["group"] == "killpod" {
		test := loadKillPod(vars["name"])
		stopKillPod(test)
		fmt.Fprintln(w, "test stop triggered")
	} else if vars["group"] == "nodeexec" {
		fmt.Fprintln(w, "can not stop a nodeexec test")
	}
}

/*
func killPodStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podName := vars["podName"]

	data := &KillPodData{}
	readJSONData(r, data)

	fmt.Fprintln(w, "pod name: ", podName)
}

func killPodStatuses(w http.ResponseWriter, r *http.Request) {
	statusesKillPod()
}


// list of all previous and current jobs
func nodeExecStatus(w http.ResponseWriter, r *http.Request) {
	nodeStatus()
}

// show specific job details and output
func nodeExecStatusDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	execId, _ := strconv.Atoi(vars["id"])

	nodeDetails(execId)
	fmt.Fprintf(w, "id: %d\n", execId)
}
*/
