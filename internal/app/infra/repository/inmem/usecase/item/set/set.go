package set

import (
	"context"
	"time"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/state"
	"github.com/vdrpkv/kvstore/internal/core/entity"
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
		if key != item.Key {
			continue
		}

		item.Val = val
		item.UpdatedAt = time.Now()
		return nil
	}

	newItem := &item.Entity{
		MetaData: entity.MetaData{
			CreatedAt: time.Now(),
		},
		Key: key,
		Val: val,
	}

	a.Items = append(a.Items, newItem)

	return nil
}
