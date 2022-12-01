package state

import (
	"sync"

	"github.com/vdrpkv/kvstore/internal/core/entity/item"
)

type State struct {
	sync.Mutex
	Items []*item.Entity
}
