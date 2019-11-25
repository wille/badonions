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

// parse does not handle cases with multiple ExitAddress
func parse(lines []string) ([]ExitNode, error) {
	nodes := make([]ExitNode, 0)

	for i := 0; i+1 < len(lines); {
		j := i + 1

		for ; j < len(lines) && strings.Split(lines[j], " ")[0] != "ExitNode"; j++ {
		}

		nodes = append(nodes, ExitNode{
			Fingerprint: strings.SplitN(lines[i], " ", 2)[1],
			Published:   strings.SplitN(lines[i+1], " ", 2)[1],
			LastStatus:  strings.SplitN(lines[i+2], " ", 2)[1],
			ExitAddress: strings.Split(lines[i+3], " ")[1],
		})

		i = j
	}

	return nodes, nil
}
