package main

import (
	"flag"
	"log"
	"net/http"
)

// running inside or outside of kubernetes
var inKubeCluster bool

func main() {
	flag.BoolVar(&inKubeCluster, "inKubeCluster", false, "indicate if we are running inside k8s cluster or not")
	flag.Parse()

	log.Printf("Running in Kubernetes cluster: %t", inKubeCluster)

	http.HandleFunc("/cmd", runner)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
