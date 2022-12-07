package get

import (
	"context"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/state"
	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
)

type Adapter struct {
	*state.State
}

var _ get.Repository = Adapter{}

func (a Adapter) FindItemByKey(ctx context.Context, key item.Key) (*item.Entity, error) {
	a.Lock()
	val, err := a.findItemByKey(ctx, key)
	a.Unlock()
	return val, err
}

func (a Adapter) findItemByKey(ctx context.Context, key item.Key) (*item.Entity, error) {
	for _, item := range a.Items {
		if key != item.Key {
			continue
		}
		return item, nil
	}
	return nil, nil
}
