package main

import (
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

func startKillPod(data KillPodData) {
	// make sure we're not already munching on this pod
	psName := data.Name + "-" + data.Namespace
	x := -1
	for i, podState := range activePods {
		if podState.name == psName {
			x = i
		}
	}

	if x != -1 {
		log.Printf("pod %s is already being munched, ignoring request\n", data.Name)
		return
	}

	podState := PodState{
		name:      psName,
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
			if victimName := killPod(data.Name, data.Kind, data.Namespace, data.Target, data.Slack); victimName != "" {
				v := Victim{podName: victimName, timeOfDeath: time.Now()}
				podState.victims = append(podState.victims, v)
			}
		}
	}()

	activePods = append(activePods, podState)
}

func stopKillPod(data KillPodData) {
	psName := data.Name + "-" + data.Namespace

	x := -1
	for i, podState := range activePods {
		if podState.name == psName {
			podState.ticker.Stop()
			x = i
		}
	}
	if x != -1 {
		log.Printf("Done snacking on %s, removing from position %d\n", data.Name, x)
		activePods = activePods[:x+copy(activePods[x:], activePods[x+1:])]
	} else {
		log.Printf("%s is not currently getting munched\n", data.Name)
	}
}

func statusesKillPod() {
	log.Print("Currently Running Pods: ")
	for _, podState := range activePods {
		log.Print(podState.name + " ")
	}
	log.Println()
}
