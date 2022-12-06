package memcachedpool

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

var ErrNoAvailableConnections = errors.New("no available connections")

func NewStaticClientPool(
	connFactory ConnectionFactory,
	connCount int,
	connWaitTimeout time.Duration,
) (ClientPool, error) {
	conns, err := createConnections(connFactory, connCount)
	if err != nil {
		return nil, err
	}

	idleConnections := make(chan memcached.Connection, connCount)
	for _, conn := range conns {
		idleConnections <- conn
	}

	return staticClientPool{
		idleConnections:     idleConnections,
		idleConnWaitTimeout: connWaitTimeout,
	}, nil
}

type staticClientPool struct {
	idleConnections     chan memcached.Connection
	idleConnWaitTimeout time.Duration
}

func (p staticClientPool) AcquireClient() (memcached.Client, error) {
	select {
	case conn := <-p.idleConnections:
		return conn, nil
	case <-time.After(p.idleConnWaitTimeout):
		return nil, ErrNoAvailableConnections
	}
}

func (p staticClientPool) ReleaseClient(client memcached.Client) {
	p.idleConnections <- castClientToConnection(client)
}

func createConnections(
	connFactory ConnectionFactory,
	connCount int,
) (
	conns []memcached.Connection,
	err error,
) {
	defer func() {
		if err == nil {
			return
		}

		err = fmt.Errorf("create connection: %w", err)

		closeErrMsgs := make([]string, 0, len(conns))
		for _, conn := range conns {
			if err := conn.Close(); err != nil {
				closeErrMsgs = append(closeErrMsgs, err.Error())
			}
		}
		conns = nil

		if len(closeErrMsgs) > 0 {
			err = fmt.Errorf(
				"close connections: [%s]: %w",
				strings.Join(closeErrMsgs, ", "), err,
			)
		}
	}()

	conns = make([]memcached.Connection, 0, connCount)
	for connCount > 0 {
		var conn memcached.Connection

		conn, err = connFactory.CreateConnection()
		if err != nil {
			return
		}

		conns = append(conns, conn)
		connCount--
	}

	return
}
