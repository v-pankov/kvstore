package get

import (
	"context"
	"fmt"
	"time"

	"github.com/vdrpkv/kvstore/internal/core/usecase"

	itemEntity "github.com/vdrpkv/kvstore/internal/core/entity/item"
	itemUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item"
)

type Request itemUsecase.BasicRequest

func (r Request) ItemKey() itemEntity.Key {
	return r.Key
}

type Response struct {
	Val       []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Processor struct {
	Gateways Gateways
}

func (p Processor) Process(ctx context.Context, req *Request) (*Response, error) {
	item, err := p.Gateways.Repository.FindItemByKey(ctx, req.Key)
	if err != nil {
		return nil, fmt.Errorf("find item by key: %w", err)
	}

	// TODO: handle not-found case
	if item == nil {
		return &Response{}, nil
	}

	return &Response{
		Val:       item.Val,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
		DeletedAt: item.DeletedAt,
	}, nil
}

type Gateways struct {
	Repository Repository
}

type Repository interface {
	FindItemByKey(ctx context.Context, key itemEntity.Key) (*itemEntity.Entity, error)
}

var _ usecase.Processor[*Request, *Response] = Processor{}
