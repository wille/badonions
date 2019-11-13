package check

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wille/badonions/internal/nodetest"
)

// ExampleTest
type ExampleTest struct{}

// Init is where the test suite is setup
// Do some initial hashing of resources
func (*ExampleTest) Init() error {
	return nil
}

// Run runs the test through each exit node
// Compare the results here with what was calculated in Init()
func (e *ExampleTest) Run(t *nodetest.T) error {
	transport := &http.Transport{
		DialContext: t.DialContext,
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	ip := string(data)

	if err == nil && ip != t.ExitNode.ExitAddress {
		t.Fail(fmt.Errorf("ExitAddress mismatch: expected %s, got %s", t.ExitNode.ExitAddress, ip))
		return nil
	}

	return err
}
