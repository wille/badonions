package exitnodes

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// ExitNode parsed from check.torproject.org
type ExitNode struct {
	Fingerprint string
	Published   string
	LastStatus  string
	ExitAddress string
}

// Get returns the exit nodes published at check.torproject.org
func Get() ([]ExitNode, error) {
	resp, err := http.Get("https://check.torproject.org/exit-addresses")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(body), "\n")
	return parse(lines)
}

func parse(lines []string) ([]ExitNode, error) {
	count := len(lines) / 4

	nodes := make([]ExitNode, count)

	for i := 0; i < count; i++ {
		m := make(map[string]string, 1)

		for l := i * 4; l < i*4+4; l++ {
			kv := strings.SplitN(lines[l], " ", 2)
			m[kv[0]] = kv[1]
		}

		nodes[i] = ExitNode{
			Fingerprint: m["ExitNode"],
			Published:   m["Published"],
			LastStatus:  m["LastStatus"],
			ExitAddress: m["ExitAddress"],
		}
	}

	return nodes, nil
}
