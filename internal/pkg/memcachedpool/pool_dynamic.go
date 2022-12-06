package memcachedpool

import (
	"container/list"
	"errors"
	"fmt"
	"sync"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

var ErrMaxConnCount = errors.New("maximum connection count is reached")

func NewDynamicClientPool(connFactory ConnectionFactory, maxConnCount int) ClientPool {
	return &dynamicClientPool{
		idleConnections:     list.New(),
		maxConnectionsCount: maxConnCount,
		connectionFactory:   connFactory,
	}
}

type dynamicClientPool struct {
	sync.Mutex
	idleConnections          *list.List
	maxConnectionsCount      int
	acquiredConnectionsCount int
	connectionFactory        ConnectionFactory
}

func (p *dynamicClientPool) AcquireClient() (memcached.Client, error) {
	p.Lock()
	client, err := p.acquireClient()
	p.Unlock()
	return client, err
}

func (p *dynamicClientPool) ReleaseClient(client memcached.Client) {
	p.Lock()
	p.releaseClient(client)
	p.Unlock()
}

func (p *dynamicClientPool) acquireClient() (memcached.Client, error) {
	if p.idleConnections.Len() > 0 {
		return p.getIdleConnection()
	}

	if p.acquiredConnectionsCount == p.maxConnectionsCount {
		return nil, ErrMaxConnCount
	}

	return p.createNewConnection()
}

func (p *dynamicClientPool) releaseClient(client memcached.Client) {
	p.idleConnections.PushBack(castClientToConnection(client))
	p.acquiredConnectionsCount--
}

func (p *dynamicClientPool) getIdleConnection() (memcached.Connection, error) {
	idleClientListEntry := p.idleConnections.Front()

	p.idleConnections.Remove(idleClientListEntry)
	p.acquiredConnectionsCount++

	conn, ok := idleClientListEntry.Value.(memcached.Connection)
	if !ok {
		panic(fmt.Sprintf(
			"connection value stored in idle connections list has unexpected type [%T]",
			idleClientListEntry.Value,
		))
	}

	return conn, nil
}

func (p *dynamicClientPool) createNewConnection() (memcached.Connection, error) {
	conn, err := p.connectionFactory.CreateConnection()
	if err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}

	p.acquiredConnectionsCount++

	return conn, nil
}
