package check

import (
	"bytes"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/wille/badonions/internal/nodetest"
)

// HTTPFileCheck verifies the integrity of resources downloaded from a remote HTTP server
// It detects tampered SSL/TLS certificates as well
type HTTPFileCheck struct {
	// URL is the remote url of the file to download through each node
	URL string

	// Sum is the sha1 checksum in hex for the file
	sum string

	originalFp string
}

// Init sends a request to the target URL checking the sum of the file and fingerprint etc
func (e *HTTPFileCheck) Init() error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Get(e.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.TLS != nil {
		certfp := sha1.Sum(resp.TLS.PeerCertificates[0].Raw)
		e.originalFp = hex.EncodeToString([]byte(certfp[:]))

		log.Printf("HTTPS detected with certificate fingerprint sha1(%s)\n", e.originalFp)
	}

	sum, err := readSum(resp.Body)
	e.sum = sum

	return err
}

// Run runs the check against URL through exit nodes
// It skips certificate validation but detects tampering with the certificate
// If the sha1sum of the downloaded resource does not match the initial sum, store the file
func (e *HTTPFileCheck) Run(t *nodetest.T) error {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: t.DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get(e.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.TLS != nil {
		sha1fp := sha1.Sum(resp.TLS.PeerCertificates[0].Raw)
		s := hex.EncodeToString([]byte(sha1fp[:]))

		if s != e.originalFp {
			log.Printf("Certificate mismatch, expected %s, got %s\n", e.originalFp, hex.EncodeToString([]byte(sha1fp[:])))
		}
	}

	var buffer bytes.Buffer
	io.Copy(&buffer, resp.Body)

	sum, err := readSum(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return err
	}

	if sum != e.sum {
		// dump buffer
		ioutil.WriteFile(fmt.Sprintf("%s-%s.dat", t.ExitNode.Fingerprint, path.Base(e.URL)), buffer.Bytes(), 0644)
		t.Fail(fmt.Errorf("File sum mismatch, expected %s, got %s", e.sum, sum))
		return nil
	}

	return err
}

func readSum(r io.Reader) (string, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
