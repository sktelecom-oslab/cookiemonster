package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type RunnerExec struct {
	Command    string   // Command name
	Parameters []string // Command parameters
}

// parse JSON from request body and return data struct
func readJSONData(r *http.Request) ([]RunnerExec, error) {
	data := []RunnerExec{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, err
	}
	if err := r.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func runner(w http.ResponseWriter, r *http.Request) {
	data, err := readJSONData(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range data {
		s := d.Command + " " + strings.Join(d.Parameters, " ")
		log.Println("Command: " + s)
		cmd := exec.Command(d.Command, d.Parameters...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			fmt.Fprintf(w, "Command failed: %s", err)
		}

		o, _ := ioutil.ReadAll(stdout)
		fmt.Fprintf(w, "%s", o)
		log.Printf("%s", o)
	}
}
