package codec

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infrastructure/transport"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/reply"
)

var (
	errUnknownCommand = errors.New("unknown command")
	errEmptyReply     = errors.New("empty reply")
	errUnknownReply   = errors.New("unknown reply")
	errNotEnoughArgs  = errors.New("not enough args")
	errTooManyArgs    = errors.New("too many args")
)

type commandEncoder struct{}

func NewCommandEncoder() transport.CommandEncoder {
	return commandEncoder{}
}

func (ce commandEncoder) EncodeCommand(cmd command.Command) ([]byte, error) {
	switch v := cmd.(type) {
	case *command.Delete:
		return []byte(
			"delete " + v.Key,
		), nil
	case *command.Gat:
		return []byte(fmt.Sprintf(
			"gat %d %s",
			v.ExpTime,
			strings.Join(v.Keys, " "),
		)), nil
	case *command.Get:
		return []byte(fmt.Sprintf(
			"get %s",
			strings.Join(v.Keys, " "),
		)), nil
	case *command.Set:
		return []byte(fmt.Sprintf(
			"set %s %d %d %d",
			v.Key, v.Flags, v.ExpTime, v.Bytes,
		)), nil
	default:
		return nil, errUnknownCommand
	}
}

type replyDecoder struct{}

func NewReplyDecoder() transport.ReplyDecoder {
	return replyDecoder{}
}

func (rd replyDecoder) DecodeReply(replyBytes []byte) (reply.Reply, error) {
	words := strings.Split(string(replyBytes), " ")

	if len(words) < 1 {
		return nil, errEmptyReply
	}

	replyType := words[0]
	words = words[1:]

	switch replyType {
	case "END":
		return createEnd(words)
	case "STORED":
		return createStored(words)
	case "DELETED":
		return createDeleted(words)
	case "NOT_FOUND":
		return createNotFound(words)
	case "ERROR":
		return createError(words)
	case "CLIENT_ERROR":
		return createClientError(words)
	case "SERVER_ERROR":
		return createServerError(words)
	case "VALUE":
		return createValue(words)
	default:
		return nil, fmt.Errorf("got reply [%s]: %w", replyType, errUnknownReply)
	}
}

func createEnd(args []string) (*reply.End, error) {
	return &reply.End{}, nil
}

func createStored(args []string) (*reply.Stored, error) {
	return &reply.Stored{}, nil
}

func createDeleted(args []string) (*reply.Deleted, error) {
	return &reply.Deleted{}, nil
}

func createNotFound(args []string) (*reply.NotFound, error) {
	return &reply.NotFound{}, nil
}

func createError(args []string) (*reply.Error, error) {
	return &reply.Error{}, nil
}

func createClientError(args []string) (*reply.ClientError, error) {
	if len(args) < 1 {
		return nil, errNotEnoughArgs
	}
	if len(args) > 1 {
		return nil, errTooManyArgs
	}
	return &reply.ClientError{Error: args[0]}, nil
}

func createServerError(args []string) (*reply.ServerError, error) {
	if len(args) < 1 {
		return nil, errNotEnoughArgs
	}
	if len(args) > 1 {
		return nil, errTooManyArgs
	}
	return &reply.ServerError{Error: args[0]}, nil
}

func createValue(args []string) (*reply.Value, error) {
	if len(args) < 3 {
		return nil, errNotEnoughArgs
	}

	if len(args) > 4 {
		return nil, errTooManyArgs
	}

	var (
		argKey   = args[0]
		argFlags = args[1]
		argBytes = args[2]
	)

	flags, err := strconv.Atoi(argFlags)
	if err != nil {
		return nil, fmt.Errorf("parse flags: %w", err)
	}

	bytes, err := strconv.Atoi(argBytes)
	if err != nil {
		return nil, fmt.Errorf("parse bytes: %w", err)
	}

	var cas int64
	if len(args) == 4 {
		cas, err = strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse cas: %w", err)
		}
	}

	return &reply.Value{
		Key:   argKey,
		Flags: int16(flags),
		Bytes: bytes,
		Cas:   cas,
	}, nil
}
