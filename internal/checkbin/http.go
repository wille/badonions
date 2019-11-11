package checkbin

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/wille/badonions/internal/nodetest"
)

// HTTPExecutableCheck downloads a file from a web server over HTTP or HTTPS without verifying the cert
// It checks the result with the provided file checksum
type HTTPExecutableCheck struct {
	// URL is the remote url of the file to download through each node
	URL string

	// Sum is the sha1 checksum in hex for the file
	Sum string
}

func (e HTTPExecutableCheck) Run(t *nodetest.T) error {
	transport := &http.Transport{
		DialContext: t.DialContext,
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
		t.Fail()
		return nil
	}

	return err
}
