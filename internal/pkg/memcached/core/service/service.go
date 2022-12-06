package service

import (
	"errors"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/reply"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/transport"
)

func ReadItems(
	replyReceiver transport.ReplyReceiver,
	dataBlockReceiver transport.DataBlockReceiver,
	itemsCount int,
) (
	[]core.Item,
	error,
) {
	items := make([]core.Item, 0, itemsCount+1) // one extra iteration to read END reply
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
			return items, core.ErrUnexpectedReply
		}

		dataBlock, err := dataBlockReceiver.ReceiveDataBlock(valueReply.Bytes)
		if err != nil {
			return items, fmt.Errorf("receive data block: %w", err)
		}

		items = append(items, core.Item{
			Key:   valueReply.Key,
			Flags: valueReply.Flags,
			Value: dataBlock,
		})
	}
	// getting here means END reply was not received
	return nil, errEndReplyWasNotReceived
}

var errEndReplyWasNotReceived = errors.New("END reply was not received")
