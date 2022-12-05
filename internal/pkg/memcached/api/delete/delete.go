package delete

import (
	"context"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/reply"
)

type Transport struct {
	CommandSender api.CommandSender
	ReplyReceiver api.ReplyReceiver
}

type Args struct {
	Key string
}

func Call(ctx context.Context, t Transport, args Args) error {
	if err := t.CommandSender.SendCommand(ctx, &command.Delete{
		Key: args.Key,
	}); err != nil {
		return fmt.Errorf("send command: %w", err)
	}

	someReply, err := t.ReplyReceiver.ReceiveReply(ctx)
	if err != nil {
		return fmt.Errorf("receive reply: %w", err)
	}

	if _, isNotFound := someReply.(*reply.NotFound); isNotFound {
		return api.ErrNotFound
	}

	if _, isDeleted := someReply.(*reply.Deleted); !isDeleted {
		return api.ErrUnexpectedReply
	}

	return nil
}
