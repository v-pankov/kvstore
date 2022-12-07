package item

import (
	"github.com/vdrpkv/kvstore/internal/core/entity"
)

type (
	Entity struct {
		entity.MetaData
		Key Key
		Val Val
	}

	Key string
	Val []byte
)

func (k Key) EqualsTo(thatKey Key) bool {
	return k == thatKey
}
