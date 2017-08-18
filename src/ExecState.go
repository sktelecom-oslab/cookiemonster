package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type Result struct {
	nodeName  string
	startTime time.Time
	output    string
}

type ExecState struct {
	name      string
	startTime time.Time
	stopTime  time.Time
	results   []Result
}

var state []ExecState
var id int = 0

func startNodeExec(name string, data ExecData) {
	n := name + "-" + strconv.Itoa(id)
	execState := ExecState{
		name:      n,
		startTime: time.Now(),
	}
	if data.Target == "" {
		log.Print("Cookie Time!!! Feast starting on random node")
	} else if data.Target == "all" {
		log.Print("Cookie Time!!! Feast starting on all nodes")
	} else {
		log.Printf("Cookie Time!!! Feast starting on node %s", data.Target)
	}

	go func() {
		if resp, err := doExec(data); err != nil {
			panic(err.Error())
		} else {
			log.Printf("reply from %s", data.Target)
			log.Println(string(resp[:]))
		}
	}()

	state = append(state, execState)
	id++
}

func statusNodeExec() string {
	var status string
	status = "Jobs:\n"
	for _, s := range state {
		status += s.name + "\n"
	}
	return status
}

func statusNodeExecJob(name string) string {
	var status string
	status = "Job Status:\n"
	for _, s := range state {
		if s.name == name {
			status += fmt.Sprintf("%v", s)
		}
	}
	return status
}
