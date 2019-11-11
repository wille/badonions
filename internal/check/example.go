package check

import (
	"io/ioutil"
	"net/http"

	"github.com/wille/badonions/internal/nodetest"
)

type ExampleTest struct{}

func (e ExampleTest) Run(t *nodetest.T) error {
	transport := &http.Transport{
		DialContext: t.DialContext,
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://api.ipify.org")
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	ip := string(data)

	if ip != t.ExitNode.ExitAddress {
		t.Fail()
		return nil
	}

	return err
}
