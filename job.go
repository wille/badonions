package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/cretz/bine/tor"
	"github.com/wille/badonions/internal/exitnodes"
)

type Job struct {
	ExitNode exitnodes.ExitNode
}

func (job *Job) Run(tor *tor.Tor) {
	dialer, err := tor.Dialer(nil, nil)
	if err != nil {
		log.Fatalf("failed to start %s: %s", job.ExitNode, err.Error())
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://api.ipify.org")
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(data))
}
