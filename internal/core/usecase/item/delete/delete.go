package delete

import (
	"context"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/core/usecase"

	itemEntity "github.com/vdrpkv/kvstore/internal/core/entity/item"
	itemUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item"
)

type Request struct {
	itemUsecase.BasicRequest
}

type Response struct {
}

type Processor struct {
	Gateways Gateways
}

func (p Processor) Process(ctx context.Context, req *Request) (*Response, error) {
	if err := p.Gateways.Repository.DeleteItemByKey(ctx, req.Key); err != nil {
		return nil, fmt.Errorf("delete item by key: %w", err)
	}
	return &Response{}, nil
}

var _ usecase.Processor[*Request, *Response] = Processor{}

type Gateways struct {
	Repository Repository
}

type Repository interface {
	DeleteItemByKey(ctx context.Context, key itemEntity.Key) error
}
