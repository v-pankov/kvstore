package gat

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/service"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/transport"
)

type Transport struct {
	CommandSender     transport.CommandSender
	ReplyReceiver     transport.ReplyReceiver
	DataBlockReceiver transport.DataBlockReceiver
}

type Args struct {
	ExpTime int
	Keys    []string
}

func Call(t Transport, args Args) ([]core.Item, error) {
	if err := t.CommandSender.SendCommand(&command.Gat{
		ExpTime: args.ExpTime,
		Keys:    args.Keys,
	}); err != nil {
		return nil, fmt.Errorf("send command: %w", err)
	}
	return service.ReadItems(t.ReplyReceiver, t.DataBlockReceiver, len(args.Keys))
}
