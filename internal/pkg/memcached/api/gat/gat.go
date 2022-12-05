package gat

import (
	"context"
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
	ExpTime int
	Keys    []string
}

func Call(ctx context.Context, t Transport, args Args) ([]api.Item, error) {
	if err := t.CommandSender.SendCommand(ctx, &command.Gat{
		ExpTime: args.ExpTime,
		Keys:    args.Keys,
	}); err != nil {
		return nil, fmt.Errorf("send command: %w", err)
	}
	return api.ReadItems(ctx, t.ReplyReceiver, t.DataBlockReceiver, len(args.Keys))
}
