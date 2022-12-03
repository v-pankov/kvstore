package set

import (
	"context"

	"github.com/vdrpkv/kvstore/generated/api/grpc"
	"github.com/vdrpkv/kvstore/internal/app/gateway/grpc/server"

	"github.com/vdrpkv/kvstore/internal/core/usecase"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
)

type Handler struct {
	Processor usecase.Processor[*set.Request, *set.Response]
}

var _ server.Handler[grpc.SetRequest, grpc.SetReply] = Handler{}

func (h Handler) Handle(ctx context.Context, req *grpc.SetRequest) (*grpc.SetReply, error) {
	_, err := h.Processor.Process(ctx, &set.Request{
		Key: req.Key,
		Val: req.Val,
	})

	if err != nil {
		return &grpc.SetReply{
			Status:  1, // TODO: add error codes
			Message: err.Error(),
		}, nil
	}

	return &grpc.SetReply{}, nil
}
