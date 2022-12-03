package item

import (
	"context"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase"
)

// KeyValidatingProcessor is a decorator for Processor interface.
// KeyValidatingProcessor validates items and then calls actual Processor.
type KeyValidatingProcessor[RequestModel KeyValidatingRequest, ResponseModel any] struct {
	KeyValidator   item.KeyValidator
	InnerProcessor usecase.Processor[RequestModel, *ResponseModel]
}

func _[RequestModel KeyValidatingRequest, ResponseModel any]() usecase.Processor[RequestModel, *ResponseModel] {
	return KeyValidatingProcessor[RequestModel, ResponseModel]{}
}

func (p KeyValidatingProcessor[RequestModel, ResponseModel]) Process(
	ctx context.Context, req RequestModel,
) (*ResponseModel, error) {
	if err := req.ItemKey().Validate(p.KeyValidator); err != nil {
		return nil, fmt.Errorf("validate item key: %w", err)
	}
	return p.InnerProcessor.Process(ctx, req)
}

func NewKeyValidationProcessor[
	RequestModel KeyValidatingRequest,
	ResponseModel any,
](
	keyValidator item.KeyValidator,
	innerProcessor usecase.Processor[
		RequestModel,
		*ResponseModel,
	],
) usecase.Processor[
	RequestModel,
	*ResponseModel,
] {
	return KeyValidatingProcessor[RequestModel, ResponseModel]{
		KeyValidator:   keyValidator,
		InnerProcessor: innerProcessor,
	}
}

// KeyValidatingRequest declares request model properties needed for KeyValidatingInteractor.
type KeyValidatingRequest interface {
	ItemKey() item.Key
}

// BasicRequest implements KeyValidatingRequest and should be used in most
// item related use case requests (because most of them contains item key).
type BasicRequest struct {
	Key []byte
}

var _ KeyValidatingRequest = BasicRequest{}

func (r BasicRequest) ItemKey() item.Key {
	return r.Key
}
