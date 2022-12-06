package memcachedpool

import (
	"net"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

type ConnectionFactory interface {
	CreateConnection() (memcached.Connection, error)
}

type TCPConnFactory struct {
	IP   net.IP
	Port int
}

var _ ConnectionFactory = TCPConnFactory{}

func (f TCPConnFactory) CreateConnection() (memcached.Connection, error) {
	return memcached.OpenTCP(f.IP, f.Port)
}
