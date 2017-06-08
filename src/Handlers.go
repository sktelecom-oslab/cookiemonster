package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type KillPodData struct {
	Name      string // name of pod base object, empty for random
	Kind      string // Kubernetes object to abuse
	Namespace string // Namespace to use, otherwise will consider all but kube-system
	Interval  int    // time between kills, unspecified for single kill
	Slack     bool   // output to Slack
	Fatal     bool   // never perform an operation which is unrecoverable
}

type RunnerExecData struct {
	Command    string   `json:"command,omitempty" protobuf:"bytes,1,opt,name=command"`
	Parameters []string `json:"parameters" protobuf:"bytes,1,rep,name=parameters"`
}

type ExecData struct {
	Target   string           // name of node, empty for random, 'all' for all
	Commands []RunnerExecData // commands to be executed
}

// parse JSON from request body
func readJSONData(r *http.Request, data interface{}) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, data); err != nil {
		panic(err)
	}

	log.Println("request data: ", data)
}

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

// Kubernetes API relates endpoints
func killPodStart(w http.ResponseWriter, r *http.Request) {
	data := &KillPodData{}
	readJSONData(r, data)
	startKillPod(*data)
	fmt.Fprintln(w, "Munching has begun")
}

func killPodStop(w http.ResponseWriter, r *http.Request) {
	data := &KillPodData{}
	readJSONData(r, data)
	stopKillPod(*data)
}

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

// Node Exec relates endpoints

// start a node job
func nodeExecStart(w http.ResponseWriter, r *http.Request) {
	data := &ExecData{}
	readJSONData(r, data)
	startNodeExec(*data)
	fmt.Fprintln(w, "Munching has begun")
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
