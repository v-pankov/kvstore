package delete

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/types"
	"github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/usecase/item/get"
	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

type Adapter struct {
	Client memcached.Client
}

var _ delete.Repository = Adapter{}

func (a Adapter) DeleteItemByKey(ctx context.Context, key item.Key) error {
	item, err := get.Adapter{Client: a.Client}.FindItemByKey(ctx, key)
	if err != nil {
		return err
	}

	data := types.Value{
		Data:      item.Val,
		CreatedAt: item.CreatedAt.Unix(),
		UpdatedAt: item.UpdatedAt.Unix(),
		DeletedAt: time.Now().Unix(),
	}

	newVal, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	if err := a.Client.Set(string(key), 0, 0, string(newVal)); err != nil {
		return fmt.Errorf("memcached set: %w", err)
	}

	return nil
}
