package main

import (
	"fmt"
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
		fmt.Printf("Starting to snack on %s\n", data.Name)
		podState := PodState{}
		podState.name = data.Name
		podState.startTime = time.Now()
		podState.ticker = time.NewTicker(time.Second * time.Duration(data.Interval))
		//podState.ticker = time.NewTicker(time.Second * 30)
		go func() {
			for t := range podState.ticker.C {
				fmt.Println("COOKIES!!!!", t)
				killPod(data.Kind, data.Namespace, data.Target)
			}
		}()

		activePods = append(activePods, podState)
	} else {
		fmt.Printf("pod %s is already being munched, ignoring request\n", data.Name)
	}

}

func stopKillPod(data KillPodData) {
	// find pod data
	x := -1
	for i, podState := range activePods {
		if podState.name == data.Name {
			podState.ticker.Stop()
			x = i
		}
	}
	if x != -1 {
		fmt.Printf("Done snacking on %s, removing from position %d\n", data.Name, x)
		activePods = activePods[:x+copy(activePods[x:], activePods[x+1:])]
	} else {
		fmt.Printf("%s is not currently getting munched\n", data.Name)
	}
}

func statusesKillPod() {
	fmt.Print("Currently Running Pods: ")
	for _, podState := range activePods {
		fmt.Print(podState.name + " ")
	}
	fmt.Println()
}
