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
	x := -1
	for i, podState := range activePods {
		if podState.name == data.Name {
			x = i
		}
	}

	if x == -1 {
		log.Printf("Starting to snack on %s\n", data.Name)
		podState := PodState{}
		podState.name = data.Name
		podState.startTime = time.Now()
		podState.ticker = time.NewTicker(time.Second * time.Duration(data.Interval))
		log.Println("COOKIES!!!!")
		go func() {
			for range podState.ticker.C {
				killPod(data.Name, data.Kind, data.Namespace, data.Target)
			}
		}()

		activePods = append(activePods, podState)
	} else {
		log.Printf("pod %s is already being munched, ignoring request\n", data.Name)
	}

}

func stopKillPod(data KillPodData) {
	x := -1
	for i, podState := range activePods {
		if podState.name == data.Name {
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
