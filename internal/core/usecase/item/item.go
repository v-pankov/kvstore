package item

import (
	"context"
	"fmt"

	"github.com/vdrpkv/kvstore/internal/core/entity/item"
	"github.com/vdrpkv/kvstore/internal/core/usecase"
)

// KeyValidatingInteractor is a decorator for Interactor interface.
// KeyValidatingInteractor validates items key from request and then
// calls actual Interactor on success.
type KeyValidatingInteractor[RequestModel KeyValidatingRequest] struct {
	KeyValidator    item.KeyValidator
	InnerInteractor usecase.Interactor[RequestModel]
}

func (p KeyValidatingInteractor[RequestModel]) Process(
	ctx context.Context, req RequestModel,
) error {
	if err := req.ItemKey().Validate(p.KeyValidator); err != nil {
		return fmt.Errorf("validate item key: %w", err)
	}

	return p.InnerInteractor.Interact(ctx, req)
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
