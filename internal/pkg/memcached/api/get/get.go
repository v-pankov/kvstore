package get

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/command"
)

type Transport struct {
	CommandSender     api.CommandSender
	ReplyReceiver     api.ReplyReceiver
	DataBlockReceiver api.DataBlockReceiver
}

type Args struct {
	Keys []string
}

func Call(t Transport, args Args) ([]api.Item, error) {
	if err := t.CommandSender.SendCommand(&command.Get{
		Keys: args.Keys,
	}); err != nil {
		return nil, fmt.Errorf("send command: %w", err)
	}
	return api.ReadItems(t.ReplyReceiver, t.DataBlockReceiver, len(args.Keys))
}
