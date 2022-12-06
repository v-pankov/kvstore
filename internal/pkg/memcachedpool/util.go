package memcachedpool

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

func castClientToConnection(client memcached.Client) memcached.Connection {
	conn, ok := client.(memcached.Connection)
	if !ok {
		panic(fmt.Sprintf(
			"unexpected client type [%T]: it's not a derivative of connection",
			client,
		))
	}
	return conn
}
