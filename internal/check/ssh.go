package check

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/wille/badonions/internal/nodetest"
	"golang.org/x/crypto/ssh"
)

type SSHFingerprintCheck struct {
	Host string

	publicKey ssh.PublicKey
}

// Init connects to the host and stores the public key fingerprint
func (e *SSHFingerprintCheck) Init() error {
	config := &ssh.ClientConfig{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			log.Printf("%s host key %s\n", remote.String(), ssh.FingerprintSHA256(key))
			e.publicKey = key
			return nil
		},
	}

	ssh.Dial("tcp", e.Host, config)

	if e.publicKey == nil {
		return fmt.Errorf("Could not get host key for %s", e.Host)
	}

	return nil
}

func (e *SSHFingerprintCheck) Run(t *nodetest.T) error {
	config := &ssh.ClientConfig{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fp := ssh.FingerprintSHA256(key)
			if fp != ssh.FingerprintSHA256(e.publicKey) {
				t.Fail(fmt.Errorf("Fingerprint mismatch! %s", fp))
			}
			return nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	conn, err := t.DialContext(ctx, "tcp", e.Host)

	// failed to etablish connection
	if err != nil {
		return err
	}
	defer conn.Close()

	c, _, _, err := ssh.NewClientConn(conn, e.Host, config)

	// failed to authenticate
	// this is ok if we don't want to keep the connection open
	if err != nil {
		return nil
	}
	c.Close()

	return err
}
