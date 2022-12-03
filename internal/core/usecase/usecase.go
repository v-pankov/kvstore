package usecase

import (
	"context"
)

// Processor handles requests to business logic and replies with its repsonse model.
type Processor[RequestModel any, ResponseModel any] interface {
	Process(ctx context.Context, requestModel RequestModel) (ResponseModel, error)
}
