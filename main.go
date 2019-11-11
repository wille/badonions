package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	flags "github.com/jessevdk/go-flags"

	"github.com/cretz/bine/tor"
	"github.com/wille/badonions/internal/checkbin"
	example "github.com/wille/badonions/internal/example_test"
	"github.com/wille/badonions/internal/exitnodes"
	"github.com/wille/badonions/internal/nodetest"
)

var checks = make(map[string]nodetest.Test)

func init() {
	checks["example"] = example.ExampleTest{}
	checks["http"] = checkbin.HTTPExecutableCheck{}
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

		dialer, err := t.Dialer(nil, nil)
		if err != nil {
			panic(err)
		}

		log.Printf("Worker %d using %s %s\n", id, job.ExitNode.Fingerprint, job.ExitNode.ExitAddress)
		err = job.Test.Run(&nodetest.T{
			DialContext: dialer.DialContext,
			ExitNode:    job.ExitNode,
		})
		log.Printf("%d: %s\n", id, err.Error())

		t.Close()

		results <- Result{job}
	}
}

var opts struct {
	Workers   int      `short:"w" long:"workers" description:"Concurrent workers"`
	TestNames []string `short:"t" long:"test" choice:"example" required:"true"`
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
