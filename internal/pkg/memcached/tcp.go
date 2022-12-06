package memcached

import (
	"errors"
	"fmt"
	"net"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api/delete"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api/gat"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api/get"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api/set"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infrastructure/codec"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infrastructure/io"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infrastructure/transport"
)

func OpenTCP(ip net.IP, port int) (Connection, error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		return nil, fmt.Errorf("dial tcp: %w", err)
	}

	return &tcpConn{
		tcpConn: conn,
	}, nil
}

type tcpConn struct {
	tcpConn *net.TCPConn
}

func (c *tcpConn) Set(key string, flags int16, exptime int, val string) error {
	return set.Call(
		set.Transport{
			CommandSender: transport.NewCommandSender(
				codec.NewCommandEncoder(),
				io.NewCommandWriter(c.tcpConn),
			),
			DataBlockSender: transport.NewDataBlockSender(
				io.NewDataBlockWriter(c.tcpConn),
			),
			ReplyReceiver: transport.NewReplyReceiver(
				io.NewReplyReader(c.tcpConn),
				codec.NewReplyDecoder(),
			),
		},
		set.Args{
			Key:     key,
			Flags:   flags,
			ExpTime: exptime,
			Value:   []byte(val),
		},
	)
}

func (c *tcpConn) Gat(exptime int, keys ...string) ([]Item, error) {
	apiItems, err := gat.Call(
		gat.Transport{
			CommandSender: transport.NewCommandSender(
				codec.NewCommandEncoder(),
				io.NewCommandWriter(c.tcpConn),
			),
			ReplyReceiver: transport.NewReplyReceiver(
				io.NewReplyReader(c.tcpConn),
				codec.NewReplyDecoder(),
			),
			DataBlockReceiver: transport.NewDataBlockReceiver(
				io.NewDataBlockReader(c.tcpConn),
			),
		},
		gat.Args{
			ExpTime: exptime,
			Keys:    keys,
		},
	)

	items := make([]Item, len(apiItems))
	for i := 0; i < len(items); i++ {
		items[0] = Item{
			Key:   apiItems[i].Key,
			Flags: apiItems[i].Flags,
			Value: string(apiItems[i].Value),
		}
	}

	return items, err
}

func (c *tcpConn) Get(keys ...string) ([]Item, error) {
	apiItems, err := get.Call(
		get.Transport{
			CommandSender: transport.NewCommandSender(
				codec.NewCommandEncoder(),
				io.NewCommandWriter(c.tcpConn),
			),
			ReplyReceiver: transport.NewReplyReceiver(
				io.NewReplyReader(c.tcpConn),
				codec.NewReplyDecoder(),
			),
			DataBlockReceiver: transport.NewDataBlockReceiver(
				io.NewDataBlockReader(c.tcpConn),
			),
		},
		get.Args{
			Keys: keys,
		},
	)

	items := make([]Item, len(apiItems))
	for i := 0; i < len(items); i++ {
		items[0] = Item{
			Key:   apiItems[i].Key,
			Flags: apiItems[i].Flags,
			Value: string(apiItems[i].Value),
		}
	}

	return items, err
}

func (c *tcpConn) Delete(key string) error {
	err := delete.Call(
		delete.Transport{
			CommandSender: transport.NewCommandSender(
				codec.NewCommandEncoder(),
				io.NewCommandWriter(c.tcpConn),
			),
			ReplyReceiver: transport.NewReplyReceiver(
				io.NewReplyReader(c.tcpConn),
				codec.NewReplyDecoder(),
			),
		},
		delete.Args{
			Key: key,
		},
	)

	if errors.Is(err, api.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (c *tcpConn) Close() error {
	return c.tcpConn.Close()
}
