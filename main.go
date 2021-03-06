package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	flags "github.com/jessevdk/go-flags"

	"github.com/cretz/bine/tor"
	"github.com/wille/badonions/internal/check"
	"github.com/wille/badonions/internal/exitnodes"
	"github.com/wille/badonions/internal/nodetest"
)

var checks = make(map[string]nodetest.Test)

func init() {
	checks["example"] = &check.ExampleTest{}
	checks["http-file"] = &check.HTTPFileCheck{
		// Good URL because some exit nodes are getting rejected by it
		URL: "https://jigsaw.w3.org/icons/jigsaw",
	}
	checks["http-basic-auth"] = &check.HTTPBasicAuthCheck{
		URL: "http://httpbin.org/basic-auth",
	}
	checks["ssh"] = &check.SSHFingerprintCheck{
		Host: "github.com:22",
	}
}

// Job is sent to a worker for processing
type Job struct {
	// ExitNode is the target exit node to check for malicious behavior
	ExitNode exitnodes.ExitNode

	// Test is the test suite this job will run
	Test nodetest.Test
}

// Result is either OK, errored or failed
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

		timeout, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		dialer, err := t.Dialer(timeout, nil)
		if err != nil {
			t.Close()
			if err == context.DeadlineExceeded {
				log.Printf("Worker %d: timeout while etablishing circuit to %s\n", id, job.ExitNode.Fingerprint)
			} else {
				log.Printf("Worker %d: init: %s\n", id, err.Error())
			}
			continue
		}

		log.Printf("Worker %d using %s %s\n", id, job.ExitNode.Fingerprint, job.ExitNode.ExitAddress)
		err = job.Test.Run(&nodetest.T{
			DialContext: dialer.DialContext,
			ExitNode:    job.ExitNode,
		})
		if err != nil {
			if err == context.DeadlineExceeded {
				log.Printf("Worker %d: timeout while connecting\n", id)
			} else {
				log.Printf("Worker %d: run: %s\n", id, err.Error())
			}
		}

		t.Close()

		results <- Result{job}
	}
}

var opts struct {
	Workers   int      `short:"w" long:"workers" description:"Concurrent workers"`
	TestNames []string `short:"t" long:"test" choice:"example" choice:"http-basic-auth" choice:"http-file" choice:"ssh" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(-1)
	}
	if opts.Workers == 0 {
		opts.Workers = runtime.NumCPU()
	}

	log.Println("Fetching active exit nodes")
	exits, err := exitnodes.Get()
	if err != nil {
		log.Fatalf("Failed to list active exit nodes: %s", err.Error())
	}
	jobcount := len(exits) * len(opts.TestNames)

	log.Printf("Iterating %d exits with %d workers\n", jobcount, opts.Workers)

	jobs := make(chan *Job, jobcount)
	results := make(chan Result, jobcount)

	for i := 0; i < opts.Workers; i++ {
		go worker(i, jobs, results)
	}

	for _, name := range opts.TestNames {
		check := checks[name]
		err := check.Init()

		if err != nil {
			log.Fatalf("Failed to initialize %s: %s\n", name, err.Error())
		}
	}

	for _, name := range opts.TestNames {
		for _, exit := range exits {
			jobs <- &Job{
				ExitNode: exit,
				Test:     checks[name],
			}
		}
	}

	close(jobs)

	for i := 0; i < jobcount; i++ {
		<-results
	}
}
