package delete

import (
	"context"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/state"
	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
)

type Adapter struct {
	*state.State
}

var _ delete.Repository = Adapter{}

func (a Adapter) DeleteItemByKey(ctx context.Context, key item.Key) error {
	a.Lock()
	err := a.deleteItemByKey(ctx, key)
	a.Unlock()
	return err
}

func (a Adapter) deleteItemByKey(ctx context.Context, key item.Key) error {
	for _, item := range a.Items {
		if !key.EqualsTo(item.Key) {
			continue
		}

		if !item.IsDeleted() {
			item.MarkDeleted()
		}
	}

	return nil
}
