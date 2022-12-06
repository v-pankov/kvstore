package gat

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/service"
)

type Transport struct {
	CommandSender     service.CommandSender
	ReplyReceiver     service.ReplyReceiver
	DataBlockReceiver service.DataBlockReceiver
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
