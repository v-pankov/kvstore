package set

import (
	"context"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/core/usecase"

	itemEntity "github.com/vdrpkv/kvstore/internal/core/entity/item"
)

type Request struct {
	Key string
	Val []byte
}

func (r Request) ItemKey() itemEntity.Key {
	return itemEntity.Key(r.Key)
}

type Response struct {
}

type Processor struct {
	Gateways Gateways
}

func (p Processor) Process(ctx context.Context, req *Request) (*Response, error) {
	if err := p.Gateways.Repository.CreateOrUpdateItem(ctx, itemEntity.Key(req.Key), req.Val); err != nil {
		return nil, fmt.Errorf("create or update item: %w", err)
	}

	return &Response{}, nil
}

type Gateways struct {
	Repository Repository
}

type Repository interface {
	CreateOrUpdateItem(ctx context.Context, key itemEntity.Key, val itemEntity.Val) error
}

var _ usecase.Processor[*Request, *Response] = Processor{}
