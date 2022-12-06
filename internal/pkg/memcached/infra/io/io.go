package io

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/vdrpkv/kvstore/internal/pkg/memcached/infra/transport"
)

type commandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) transport.CommandWriter {
	return commandWriter{
		writer: writer,
	}
}

func (cw commandWriter) WriteCommand(commandBytes []byte) error {
	if _, err := cw.writer.Write(commandBytes); err != nil {
		return fmt.Errorf("write command bytes: %w", err)
	}
	if _, err := cw.writer.Write(terminator); err != nil {
		return fmt.Errorf("write terminator: %w", err)
	}
	return nil
}

type replyReader struct {
	reader io.Reader
}

func NewReplyReader(reader io.Reader) transport.ReplyReader {
	return replyReader{
		reader: reader,
	}
}

func (rr replyReader) ReadReply() ([]byte, error) {
	return readBytesUntilTerminator(rr.reader, 256, 1024)
}

type dataBlockWriter struct {
	writer io.Writer
}

func NewDataBlockWriter(writer io.Writer) transport.DataBlockWriter {
	return dataBlockWriter{
		writer: writer,
	}
}

func (dbw dataBlockWriter) WriteDataBlock(dataBlock []byte) error {
	if _, err := dbw.writer.Write(dataBlock); err != nil {
		return fmt.Errorf("write data block bytes: %w", err)
	}
	if _, err := dbw.writer.Write(terminator); err != nil {
		return fmt.Errorf("write terminator: %w", err)
	}
	return nil
}

type dataBlockReader struct {
	reader io.Reader
}

func NewDataBlockReader(reader io.Reader) transport.DataBlockReader {
	return dataBlockReader{
		reader: reader,
	}
}

func (dbr dataBlockReader) ReadDataBlock(bytes int) ([]byte, error) {
	buf := make([]byte, bytes+2)
	n, err := dbr.reader.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	if n < bytes+2 {
		return nil, errIncompleteRead
	}
	if buf[bytes] != asciiCarriageR || buf[bytes+1] != asciiCarriageN {
		return nil, errProtocolViolation
	}
	return buf, nil
}

func readBytesUntilTerminator(r io.Reader, bufSize int, limit int) ([]byte, error) {
	var total int
	byteBuf := make([]byte, 1)
	outBuf := bytes.NewBuffer(nil) // todo: review this part
	for total < limit {
		b, err := readByte(r, byteBuf)
		if err != nil {
			return nil, fmt.Errorf("read byte: %w", err)
		}
		if b == asciiCarriageR {
			return outBuf.Bytes(), tryReadByte(
				r, byteBuf, asciiCarriageN, errProtocolViolation,
			)
		}
		outBuf.WriteByte(b)
		total++
	}
	return nil, errReadLimitExceeded
}

func tryReadByte(r io.Reader, buf []byte, wantByte byte, onMismatch error) error {
	b, err := readByte(r, buf)

	if err != nil {
		return fmt.Errorf("read byte: %w", err)
	}

	if b != wantByte {
		return onMismatch
	}

	return nil
}

func readByte(r io.Reader, buf []byte) (byte, error) {
	n, err := r.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("read: %w", err)
	}
	if n != 1 {
		return 0, errIncompleteRead
	}
	return buf[0], nil
}

var (
	errReadLimitExceeded = errors.New("read limit exceeded")
	errProtocolViolation = errors.New("protocol violation")
	errIncompleteRead    = errors.New("incomplete read")
)

const (
	asciiCarriageR = 13 // "\r"
	asciiCarriageN = 10 // "\n"
)

var (
	terminator = []byte{asciiCarriageR, asciiCarriageN}
)
