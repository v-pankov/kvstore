package delete

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/reply"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/service"
)

type Transport struct {
	CommandSender service.CommandSender
	ReplyReceiver service.ReplyReceiver
}

type Args struct {
	Key string
}

func Call(t Transport, args Args) error {
	if err := t.CommandSender.SendCommand(&command.Delete{
		Key: args.Key,
	}); err != nil {
		return fmt.Errorf("send command: %w", err)
	}

	someReply, err := t.ReplyReceiver.ReceiveReply()
	if err != nil {
		return fmt.Errorf("receive reply: %w", err)
	}

	if _, isNotFound := someReply.(*reply.NotFound); isNotFound {
		return core.ErrNotFound
	}

	if _, isDeleted := someReply.(*reply.Deleted); !isDeleted {
		return core.ErrUnexpectedReply
	}

	return nil
}
