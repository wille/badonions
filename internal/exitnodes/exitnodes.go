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
		ln := i * 4

		nodes[i] = ExitNode{
			Fingerprint: strings.SplitN(lines[ln], " ", 2)[1],
			Published:   strings.SplitN(lines[ln+1], " ", 2)[1],
			LastStatus:  strings.SplitN(lines[ln+2], " ", 2)[1],
			ExitAddress: strings.Split(lines[ln+3], " ")[1],
		}
	}

	return nodes, nil
}
