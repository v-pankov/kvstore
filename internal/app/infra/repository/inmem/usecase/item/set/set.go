package set

import (
	"context"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/state"
	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
)

type Adapter struct {
	*state.State
}

var _ set.Repository = Adapter{}

func (a Adapter) CreateOrUpdateItem(ctx context.Context, key item.Key, val item.Val) error {
	a.Lock()
	err := a.createOrUpdateItem(ctx, key, val)
	a.Unlock()
	return err
}

func (a Adapter) createOrUpdateItem(ctx context.Context, key item.Key, val item.Val) error {
	for _, item := range a.Items {
		if !key.EqualsTo(item.Key) {
			continue
		}

		item.Val = val
		item.MarkUpdated()
		return nil
	}

	newItem := &item.Entity{
		Key: key,
		Val: val,
	}
	newItem.MarkCreated()

	a.Items = append(a.Items, newItem)

	return nil
}
