package usecase

import (
	"context"
	"fmt"
)

// Processor handles requests to business logic and replies with its repsonse model.
type Processor[RequestModel any, ResponseModel any] interface {
	Process(ctx context.Context, requestModel RequestModel) (ResponseModel, error)
}

// Interactor sends request to business logic but doesn't know about its response model.
// In other words, Interactor hides response model information from the caller.
type Interactor[RequestModel any] interface {
	Interact(ctx context.Context, requestModel RequestModel) error
}

// Presenter handles business logic response model and doesn't know about its request model.
// In other words, Presenter hides request model information from the caller.
type Presenter[ResponseModel any] interface {
	Present(ctx context.Context, responseModel ResponseModel) error
}

// NewInteractor creates Interactor by combining Processor and Presenter.
func NewInteractor[RequestModel any, ResponseModel any](
	processor Processor[RequestModel, ResponseModel],
	presenter Presenter[ResponseModel],
) Interactor[RequestModel] {
	return interactor[RequestModel, ResponseModel]{
		processor: processor,
		presenter: presenter,
	}
}

// interactor implements Interactor interface
type interactor[RequestModel any, ResponseModel any] struct {
	processor Processor[RequestModel, ResponseModel]
	presenter Presenter[ResponseModel]
}

func (i interactor[RequestModel, ResponseModel]) Interact(ctx context.Context, requestModel RequestModel) error {
	rsp, err := i.processor.Process(ctx, requestModel)
	if err != nil {
		return fmt.Errorf("process request: %w", err)
	}

	if err := i.presenter.Present(ctx, rsp); err != nil {
		return fmt.Errorf("present response: %w", err)
	}

	return nil
}
