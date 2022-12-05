package set

import (
	"context"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/reply"
)

type Transport struct {
	CommandSender   api.CommandSender
	DataBlockSender api.DataBlockSender
	ReplyReceiver   api.ReplyReceiver
}

type Args struct {
	Key     string
	Flags   int16
	ExpTime int
	Value   []byte
}

func Call(ctx context.Context, t Transport, args Args) error {
	if err := t.CommandSender.SendCommand(ctx, &command.Set{
		Key:     args.Key,
		Flags:   args.Flags,
		ExpTime: args.ExpTime,
		Bytes:   len(args.Value),
	}); err != nil {
		return fmt.Errorf("send command: %w", err)
	}

	if err := t.DataBlockSender.SendDataBlock(ctx, args.Value); err != nil {
		return fmt.Errorf("send data block: %w", err)
	}

	someReply, err := t.ReplyReceiver.ReceiveReply(ctx)
	if err != nil {
		return fmt.Errorf("receive reply: %w", err)
	}

	if _, ok := someReply.(*reply.Stored); !ok {
		return api.ErrUnexpectedReply
	}

	return nil
}
