package core

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrUnexpectedReply = errors.New("unexpected reply")
)

type Item struct {
	Key   string
	Flags int16
	Value []byte
}
