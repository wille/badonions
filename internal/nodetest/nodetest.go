package nodetest

import (
	"context"
	"net"

	"github.com/wille/badonions/internal/exitnodes"
)

type Test interface {
	Run(*T) error
}

type T struct {
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
	ExitNode    exitnodes.ExitNode
}

func (t *T) Fail() {

}
