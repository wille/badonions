package check

import (
	"math/rand"
	"net/http"

	"github.com/wille/badonions/internal/nodetest"
)

// HTTPBasicAuthCheck
type HTTPBasicAuthCheck struct {
	URL string
}

func (e HTTPBasicAuthCheck) Init() error {
	return nil
}

func (e HTTPBasicAuthCheck) Run(t *nodetest.T) error {
	transport := &http.Transport{
		DialContext: t.DialContext,
	}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("POST", e.URL, nil)
	if err != nil {
		return err
	}

	user := "admin1"
	pass := randomString(8)

	req.SetBasicAuth(user, pass)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
