package get

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/types"
	"github.com/vdrpkv/kvstore/internal/core/entity"
	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

type Adapter struct {
	Client memcached.Client
}

var _ get.Repository = Adapter{}

func (a Adapter) FindItemByKey(ctx context.Context, key item.Key) (*item.Entity, error) {
	items, err := a.Client.Get(string(key))

	if err != nil {
		return nil, fmt.Errorf("memcached get: %w", err)
	}

	if len(items) == 0 {
		return nil, entity.ErrNotFound
	}

	firstItem := items[0]

	var data types.Value
	if err := json.Unmarshal([]byte(firstItem.Value), &data); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return &item.Entity{
		MetaData: entity.MetaData{
			CreatedAt: time.Unix(data.CreatedAt, 0),
			UpdatedAt: time.Unix(data.UpdatedAt, 0),
			DeletedAt: time.Unix(data.DeletedAt, 0),
		},
		Key: item.Key(firstItem.Key),
		Val: data.Data,
	}, nil
}
