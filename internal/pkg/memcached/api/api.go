package api

import (
	"errors"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/reply"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrUnexpectedReply = errors.New("unexpected reply")

	errEndReplyWasNotReceived = errors.New("END reply was not received")
)

type (
	CommandSender interface {
		SendCommand(command.Command) error
	}

	ReplyReceiver interface {
		ReceiveReply() (reply.Reply, error)
	}

	DataBlockSender interface {
		SendDataBlock([]byte) error
	}

	DataBlockReceiver interface {
		ReceiveDataBlock(int) ([]byte, error)
	}
)

type Item struct {
	Key   string
	Flags int16
	Value []byte
}

func ReadItems(
	replyReceiver ReplyReceiver,
	dataBlockReceiver DataBlockReceiver,
	itemsCount int,
) (
	[]Item,
	error,
) {
	items := make([]Item, 0, itemsCount+1) // one extra iteration to read END reply
	for i := 0; i < itemsCount+1; i++ {
		someReply, err := replyReceiver.ReceiveReply()
		if err != nil {
			return items, fmt.Errorf("receive reply: %w", err)
		}

		if _, isEnd := someReply.(*reply.End); isEnd {
			return items, nil // success
		}

		valueReply, isValueReply := someReply.(*reply.Value)
		if !isValueReply {
			return items, ErrUnexpectedReply
		}

		dataBlock, err := dataBlockReceiver.ReceiveDataBlock(valueReply.Bytes)
		if err != nil {
			return items, fmt.Errorf("receive data block: %w", err)
		}

		items = append(items, Item{
			Key:   valueReply.Key,
			Flags: valueReply.Flags,
			Value: dataBlock,
		})
	}
	// getting here means END reply was not received
	return nil, errEndReplyWasNotReceived
}
