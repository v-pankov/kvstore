package memcached

import (
	"errors"
)

type Item struct {
	Key   string
	Flags int16
	Value string
}

var ErrNotFound = errors.New("key is not found")

type Connection interface {
	Set(key string, flags int16, exptime int, val string) error
	Gat(exptime int, key ...string) ([]Item, error)
	Get(key ...string) ([]Item, error)
	Delete(key string) error
	Close() error
}
