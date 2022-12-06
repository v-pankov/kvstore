package set

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/reply"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/service"
)

type Transport struct {
	CommandSender   service.CommandSender
	DataBlockSender service.DataBlockSender
	ReplyReceiver   service.ReplyReceiver
}

type Args struct {
	Key     string
	Flags   int16
	ExpTime int
	Value   []byte
}

func Call(t Transport, args Args) error {
	if err := t.CommandSender.SendCommand(&command.Set{
		Key:     args.Key,
		Flags:   args.Flags,
		ExpTime: args.ExpTime,
		Bytes:   len(args.Value),
	}); err != nil {
		return fmt.Errorf("send command: %w", err)
	}

	if err := t.DataBlockSender.SendDataBlock(args.Value); err != nil {
		return fmt.Errorf("send data block: %w", err)
	}

	someReply, err := t.ReplyReceiver.ReceiveReply()
	if err != nil {
		return fmt.Errorf("receive reply: %w", err)
	}

	if _, ok := someReply.(*reply.Stored); !ok {
		return core.ErrUnexpectedReply
	}

	return nil
}
