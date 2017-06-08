package main

import (
	"log"
	"time"
)

type Result struct {
	nodeName  string
	startTime time.Time
	output    string
}

type ExecState struct {
	id        int
	startTime time.Time
	stopTime  time.Time
	results   []Result
}

var state []ExecState

func startNodeExec(data ExecData) {
	execState := ExecState{
		id:        123,
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
}

func nodeStatus() {
	log.Println("Node Exec Status")
	for _, s := range state {
		log.Printf("Id: %d", s.id)
	}
}

func nodeDetails(id int) {
	for _, s := range state {
		log.Printf("Id: %d", s.id)
		if id == s.id {
			log.Println("Match!")
		}
	}
}
