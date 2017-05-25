package main

import (
	"flag"
	"log"
	"net/http"
)

// running inside or outside of kubernetes
var inCluster bool

func main() {
	flag.BoolVar(&inCluster, "inCluster", false, "indicate if we are running inside k8s cluster or not")
	log.Printf("Running in cluster: %t", inCluster)

	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
