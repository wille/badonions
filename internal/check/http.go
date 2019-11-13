package check

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/wille/badonions/internal/nodetest"
)

// HTTPExecutableCheck downloads a file from a web server over HTTP or HTTPS without verifying the cert
// It checks the result with the provided file checksum
type HTTPFileCheck struct {
	// URL is the remote url of the file to download through each node
	URL string

	// Sum is the sha1 checksum in hex for the file
	Sum string
}

	return nil
func (e *HTTPFileCheck) Init() error {
}

	transport := &http.Transport{
		DialContext: t.DialContext,
func (e *HTTPFileCheck) Run(t *nodetest.T) error {
	}
	client := &http.Client{Transport: transport}

	resp, err := client.Get(e.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	h := sha1.New()
	_, err = io.Copy(h, resp.Body)
	if err != nil {
		return err
	}

	sum := hex.EncodeToString(h.Sum(nil))

	if sum != e.Sum {
		t.Fail(fmt.Errorf("File sum mismatch, expected %s, got %s", e.Sum, sum))
		return nil
	}

	return err
}
