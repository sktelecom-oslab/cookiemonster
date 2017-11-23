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
	"time"
)

type Victim struct {
	podName     string
	timeOfDeath time.Time
}

type PodState struct {
	name      string
	ticker    *time.Ticker
	startTime time.Time
	stopTime  time.Time
	kills     int
	victims   []Victim
}

var activePods []PodState

func startKillPod(jobName string, data KillPodData) {
	x := false
	for _, podState := range activePods {
		if podState.name == jobName {
			x = true
		}
	}

	if x == true {
		log.Printf("pod %s is already being munched, ignoring request\n", jobName)
		return
	}

	podState := PodState{
		name:      jobName,
		startTime: time.Now(),
		ticker:    time.NewTicker(time.Second * time.Duration(data.Interval)),
	}
	if data.Name == "" {
		log.Printf("Cookie Time!!! Random feast starting on %s in namespace %s", data.Kind, data.Namespace)
	} else {
		log.Printf("Cookie Time!!! Feast starting on %s %s in namespace %s", data.Name, data.Kind, data.Namespace)
	}
	go func() {
		for range podState.ticker.C {
			if victimName := killPod(data.Name, data.Kind, data.Namespace, data.Slack); victimName != "" {
				v := Victim{podName: victimName, timeOfDeath: time.Now()}
				podState.victims = append(podState.victims, v)
			}
		}
	}()

	activePods = append(activePods, podState)
}

func stopKillPod(name string, data KillPodData) {
	x := -1
	for i, podState := range activePods {
		if podState.name == name {
			podState.ticker.Stop()
			x = i
		}
	}
	if x != -1 {
		log.Printf("Done snacking on %s\n", name)
		activePods = activePods[:x+copy(activePods[x:], activePods[x+1:])]
	} else {
		log.Printf("%s is not currently getting munched\n", name)
	}
}

func statusKillPod() string {
	var status string
	status = "Running Jobs:\n"
	for _, p := range activePods {
		status += p.name + "\n"
	}
	return status
}

func statusKillPodJob(name string) string {
	var status string
	status = "Job Status:\n"
	for _, p := range activePods {
		if p.name == name {
			status += fmt.Sprintf("%v", p)
		}
	}
	return status
}
