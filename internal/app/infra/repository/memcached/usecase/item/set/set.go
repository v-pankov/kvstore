package set

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/types"
	"github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/usecase/item/get"
	"github.com/vdrpkv/kvstore/internal/core/entity"
	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
	"github.com/vdrpkv/kvstore/internal/pkg/memcached"
)

type Adapter struct {
	Client memcached.Client
}

var _ set.Repository = Adapter{}

func (a Adapter) CreateOrUpdateItem(ctx context.Context, key item.Key, val item.Val) error {
	item, err := get.Adapter{Client: a.Client}.FindItemByKey(ctx, key)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return err
	}

	var (
		timeNow = time.Now().Unix()
		data    = types.Value{
			Data:      val,
			CreatedAt: timeNow,
		}
	)

	if err == nil {
		data.CreatedAt = item.CreatedAt.Unix()
		data.UpdatedAt = timeNow
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
