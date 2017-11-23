/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
