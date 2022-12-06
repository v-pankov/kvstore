package transport

import (
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/core/reply"
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
