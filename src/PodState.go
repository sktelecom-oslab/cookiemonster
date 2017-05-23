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
  podData KillPodData
  StartTime time.Time
  StopTime time.Time
  Kills int
  Victims []Victim
}

var activePods []PodState

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
    podState := PodState{data, time.Now(), time.Now(), 0, nil}
    activePods = append(activePods, podState)
  } else {
    fmt.Printf("pod %s is already being munched, ignoring request\n", data.Name)
  }
}

func stopKillPod(podName string) {
  // find pod data
  x := -1
  for i, podState := range activePods {
    if podState.podData.Name == podName {
      x = i
    }
  }
  if x > 0 {
    fmt.Printf("Done snacking on %s, removing from position %d\n", podName, x)
    activePods = activePods[:x+copy(activePods[x:], activePods[x+1:])]
  } else {
    fmt.Printf("%s is not currently getting munched\n", podName)
  }
}

func statusesKillPod() {
  fmt.Print("Currently Running Pods: ")
  for _, podState := range activePods {
    fmt.Print(podState.podData.Name + " ")
  }
  fmt.Println("")
}
