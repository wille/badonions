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

func (t *T) Fail(err error) {
	log.Printf("%s fail: %s\n", t.ExitNode.Fingerprint, err.Error())
}
