package transport

import (
	"fmt"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/api"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/command"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached/reply"
)

type commandSender struct {
	commandEncoder CommandEncoder
	commandWriter  CommandWriter
}

func NewCommandSender(
	commandEncoder CommandEncoder,
	commandWriter CommandWriter,
) api.CommandSender {
	return commandSender{
		commandEncoder: commandEncoder,
		commandWriter:  commandWriter,
	}
}

func (cs commandSender) SendCommand(cmd command.Command) error {
	commandBytes, err := cs.commandEncoder.EncodeCommand(cmd)
	if err != nil {
		return fmt.Errorf("encode command: %w", err)
	}

	if err := cs.commandWriter.WriteCommand(commandBytes); err != nil {
		return fmt.Errorf("write command: %w", err)
	}

	return nil
}

type replyReceiver struct {
	replyReader  ReplyReader
	replyDecoder ReplyDecoder
}

func NewReplyReceiver(
	replyReader ReplyReader,
	replyDecoder ReplyDecoder,
) api.ReplyReceiver {
	return replyReceiver{
		replyReader:  replyReader,
		replyDecoder: replyDecoder,
	}
}

func (rr replyReceiver) ReceiveReply() (reply.Reply, error) {
	replyBytes, err := rr.replyReader.ReadReply()
	if err != nil {
		return nil, fmt.Errorf("read reply: %w", err)
	}

	reply, err := rr.replyDecoder.DecodeReply(replyBytes)
	if err != nil {
		return nil, fmt.Errorf("decode reply: %w", err)
	}

	return reply, nil
}

type dataBlockSender struct {
	dataBlockWriter DataBlockWriter
}

func NewDataBlockSender(dataBlockWriter DataBlockWriter) api.DataBlockSender {
	return dataBlockSender{
		dataBlockWriter: dataBlockWriter,
	}
}

func (dbs dataBlockSender) SendDataBlock(dataBlock []byte) error {
	if err := dbs.dataBlockWriter.WriteDataBlock(dataBlock); err != nil {
		return fmt.Errorf("write data block: %w", err)
	}
	return nil
}

type dataBlockReceiver struct {
	dataBlockReader DataBlockReader
}

func NewDataBlockReceiver(dataBlockReader DataBlockReader) api.DataBlockReceiver {
	return dataBlockReceiver{
		dataBlockReader: dataBlockReader,
	}
}

func (dbr dataBlockReceiver) ReceiveDataBlock(bytes int) ([]byte, error) {
	dataBlock, err := dbr.dataBlockReader.ReadDataBlock(bytes)
	if err != nil {
		return nil, fmt.Errorf("read data block: %w", err)
	}
	return dataBlock, nil
}

type CommandEncoder interface {
	EncodeCommand(command.Command) ([]byte, error)
}

type ReplyDecoder interface {
	DecodeReply([]byte) (reply.Reply, error)
}

type CommandWriter interface {
	WriteCommand([]byte) error
}

type ReplyReader interface {
	ReadReply() ([]byte, error)
}

type DataBlockWriter interface {
	WriteDataBlock([]byte) error
}

type DataBlockReader interface {
	ReadDataBlock(int) ([]byte, error)
}
