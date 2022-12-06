package memcached

import (
	"errors"
	"fmt"
	"net"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/client/delete"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/client/gat"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/client/get"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/client/set"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infra/codec"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infra/io"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infra/transport"

	coreTransport "github.com/vdrpkv/kvstore/internal/pkg/memcached/core/transport"
)

func OpenTCP(ip net.IP, port int) (Connection, error) {
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		return nil, fmt.Errorf("dial tcp: %w", err)
	}

	var (
		commandEncoder = codec.NewCommandEncoder()
		replyDecoder   = codec.NewReplyDecoder()
	)

	return &tcpConn{
		conn: conn,

		commandSender: transport.NewCommandSender(
			commandEncoder,
			io.NewCommandWriter(conn),
		),
		replyReceiver: transport.NewReplyReceiver(
			io.NewReplyReader(conn),
			replyDecoder,
		),
		dataBlockSender: transport.NewDataBlockSender(
			io.NewDataBlockWriter(conn),
		),
		dataBlockReceiver: transport.NewDataBlockReceiver(
			io.NewDataBlockReader(conn),
		),
	}, nil
}

type tcpConn struct {
	conn *net.TCPConn

	commandSender     coreTransport.CommandSender
	replyReceiver     coreTransport.ReplyReceiver
	dataBlockSender   coreTransport.DataBlockSender
	dataBlockReceiver coreTransport.DataBlockReceiver
}

func (c *tcpConn) Set(key string, flags int16, exptime int, val string) error {
	return set.Call(
		set.Transport{
			CommandSender:   c.commandSender,
			DataBlockSender: c.dataBlockSender,
			ReplyReceiver:   c.replyReceiver,
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
	coreItems, err := gat.Call(
		gat.Transport{
			CommandSender:     c.commandSender,
			ReplyReceiver:     c.replyReceiver,
			DataBlockReceiver: c.dataBlockReceiver,
		},
		gat.Args{
			ExpTime: exptime,
			Keys:    keys,
		},
	)
	return mapCoreItems(coreItems), err
}

func (c *tcpConn) Get(keys ...string) ([]Item, error) {
	coreItems, err := get.Call(
		get.Transport{
			CommandSender:     c.commandSender,
			ReplyReceiver:     c.replyReceiver,
			DataBlockReceiver: c.dataBlockReceiver,
		},
		get.Args{
			Keys: keys,
		},
	)
	return mapCoreItems(coreItems), err
}

func (c *tcpConn) Delete(key string) error {
	err := delete.Call(
		delete.Transport{
			CommandSender: c.commandSender,
			ReplyReceiver: c.replyReceiver,
		},
		delete.Args{
			Key: key,
		},
	)

	if errors.Is(err, core.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (c *tcpConn) Close() error {
	return c.conn.Close()
}
