package check

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/wille/badonions/internal/exitnodes"
	"github.com/wille/badonions/internal/nodetest"
)

// HTTPBasicAuthCheck
type HTTPBasicAuthCheck struct {
	URL string
}

func (e *HTTPBasicAuthCheck) Init() error {
	return nil
}

func (e *HTTPBasicAuthCheck) Run(t *nodetest.T) error {
	transport := &http.Transport{
		DialContext: t.DialContext,
	}
	client := &http.Client{Transport: transport}

	user := "admin"
	pass := randomString(16)

	url := fmt.Sprintf("%s/%s/%s", e.URL, user, pass)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	storeFingerprintCredentials(e, t.ExitNode, user, pass)

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fail(fmt.Errorf("Failed to login with provided credentials, HTTP %s", resp.Status))
	}

	return nil
}

func storeFingerprintCredentials(e *HTTPBasicAuthCheck, exit exitnodes.ExitNode, user, pass string) {
	log.Printf("%s: %s:%s\t\t%s", exit.Fingerprint, user, pass, e.URL)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
