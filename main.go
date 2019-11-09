package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cretz/bine/tor"
	"github.com/wille/badonions/internal/exitnodes"
)

type Result struct {
	job *Job
}

func worker(id int, jobs <-chan *Job, results chan<- Result) {
	datadir, _ := ioutil.TempDir("", fmt.Sprintf("badonions-%d", id))

	for job := range jobs {
		t, err := tor.Start(nil, &tor.StartConf{
			DataDir: datadir,
			ExtraArgs: []string{
				"ExitNodes",
				job.ExitNode.Fingerprint,
			},
		})

		if err != nil {
			log.Fatalf("failed to start %s: %s", job.ExitNode, err.Error())
		}

		log.Printf("Worker %d using %s %s \n", id, job.ExitNode.Fingerprint, job.ExitNode.ExitAddress)
		job.Run(t)
		t.Close()
		results <- Result{job}
	}
}

const workers = 1

func main() {
	log.Println("Fetching active exit nodes")
	exits, err := exitnodes.Get()
	if err != nil {
		log.Fatalf("Failed to list active exit nodes: %s", err.Error())
	}
	jobcount := len(exits)

	log.Printf("Iterating %d exits with %d workers\n", jobcount, workers)

	jobs := make(chan *Job, jobcount)
	results := make(chan Result, jobcount)

	for i := 0; i < workers; i++ {
		go worker(i, jobs, results)
	}

	for _, exit := range exits {
		jobs <- &Job{
			ExitNode: exit,
		}
	}

	close(jobs)

	for i := 0; i < jobcount; i++ {
		<-results
	}
}
