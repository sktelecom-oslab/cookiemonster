package main

import (
  "fmt"
  "time"
)

type Victim struct {
  PodName string
  TimeOfDeath time.Time
}

type PodState struct {
  PodData KillPodData
  StartTime time.Time
  StopTime time.Time
  Kills int
  Victims []Victim
}

var activePods []PodState
var random string = "RANDOM"

func startKillPod(data KillPodData)  {
  // make sure we're not already munching on this pod
  x := -1
  for i, podState := range activePods {
    if podState.podData.Name == data.Name {
      x = i
    }
  }

  if x == -1 {
    fmt.Printf("Starting to snack on %s\n", data.Name)
    podState := PodState{}
    podState.PodData = data
    podState.StartTime = time.Now()
    
    activePods = append(activePods, podState)
  } else {
    fmt.Printf("pod %s is already being munched, ignoring request\n", data.Name)
  }
}

func stopKillPod(data KillPodData) {
  // find pod data
  x := -1
  for i, podState := range activePods {
    if podState.podData.Name == data.Name {
      x = i
    }
  }
  if x > 0 {
    fmt.Printf("Done snacking on %s, removing from position %d\n", data.Name, x)
    activePods = activePods[:x+copy(activePods[x:], activePods[x+1:])]
  } else {
    fmt.Printf("%s is not currently getting munched\n", data.Name)
  }
}

func statusesKillPod() {
  fmt.Print("Currently Running Pods: ")
  for _, podState := range activePods {
    fmt.Print(podState.podData.Name + " ")
  }
  fmt.Println("")
}
