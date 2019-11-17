package check

import (
	"fmt"
	"log"
	"net"

	"github.com/wille/badonions/internal/nodetest"
	"golang.org/x/crypto/ssh"
)

type SSHFingerprintCheck struct {
	Host string

	publicKey ssh.PublicKey
}

func (e *SSHFingerprintCheck) Init() error {
	config := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.Password("test"),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fmt.Println("HostKeyCallback", hostname, remote, key.Type(), ssh.FingerprintSHA256(key))
			e.publicKey = key
			return nil
		},
	}

	_, err := ssh.Dial("tcp", e.Host, config)

	if err != nil {
		log.Printf("WARN: ssh init failed: %s\n", err.Error())
	}

	return nil
}

func (e *SSHFingerprintCheck) Run(t *nodetest.T) error {
	config := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.Password("test"),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			fmt.Println("HostKeyCallback", hostname, remote, key.Type(), ssh.FingerprintSHA256(key))
			e.publicKey = key
			return nil
		},
	}

	conn, err := t.Dial("tcp", e.Host)
	c, chans, reqs, err := ssh.NewClientConn(conn, e.Host, config)
	if err != nil {
		conn.Close()
		return err
	}
	client := ssh.NewClient(c, chans, reqs)
	sess, _ := client.NewSession()
	sess.Run("ls")
	sess.Close()
	client.Close()
	conn.Close()

	return err
}
