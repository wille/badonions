package nodetest

import (
	"context"
	"log"
	"net"

	"github.com/wille/badonions/internal/exitnodes"
)

type Test interface {
	Init() error
	Run(*T) error
}

type T struct {
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
	ExitNode    exitnodes.ExitNode
}

// Dial is a wrapper for DialContext since be only have a underlying DialContext function
func (t *T) Dial(network, addr string) (net.Conn, error) {
	return t.DialContext(context.Background(), network, addr)
}

func (t *T) Fail(err error) {
	log.Printf("%s fail: %s\n", t.ExitNode.Fingerprint, err.Error())
}
