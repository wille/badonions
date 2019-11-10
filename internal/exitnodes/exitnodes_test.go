package exitnodes

import (
	"testing"
)

func TestParse(t *testing.T) {
	nodes, _ := parse([]string{
		"ExitNode AABBCCDDEEFF",
		"Published YYYY-mm-dd hh:mm:ss",
		"LastStatus YYYY-mm-dd hh:mm:ss",
		"ExitAddress 127.0.0.1 YYYY-mm-dd hh:mm:ss",

		"ExitNode 2",
		"Published YYYY-mm-dd hh:mm:ss",
		"LastStatus YYYY-mm-dd hh:mm:ss",
		"ExitAddress 127.0.0.1 YYYY-mm-dd hh:mm:ss",

		"ExitNode 3",
		"Published YYYY-mm-dd hh:mm:ss",
		"LastStatus YYYY-mm-dd hh:mm:ss",
		"ExitAddress 127.0.0.1 YYYY-mm-dd hh:mm:ss",
	})

	if nodes[0].Fingerprint != "AABBCCDDEEFF" {
		t.Errorf("mismatch on Fingerprint")
	}

	if nodes[1].Published != "YYYY-mm-dd hh:mm:ss" {
		t.Errorf("mismatch on Published")
	}

	if nodes[1].LastStatus != "YYYY-mm-dd hh:mm:ss" {
		t.Errorf("mismatch on LastStatus")
	}

	if nodes[2].ExitAddress != "127.0.0.1" {
		t.Errorf("mismatch on ExitAddress")
	}

	if nodes[2].Fingerprint != "3" {
		t.Errorf("mismatch on Fingerprint")
	}
}
