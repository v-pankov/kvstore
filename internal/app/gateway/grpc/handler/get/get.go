package get

import (
	"context"

	"github.com/vdrpkv/kvstore/generated/api/grpc"
	"github.com/vdrpkv/kvstore/internal/app/gateway/grpc/server"

	"github.com/vdrpkv/kvstore/internal/core/usecase"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
)

type Handler struct {
	Processor usecase.Processor[*get.Request, *get.Response]
}

var _ server.Handler[grpc.GetRequest, grpc.GetReply] = Handler{}

func (h Handler) Handle(ctx context.Context, req *grpc.GetRequest) (*grpc.GetReply, error) {
	rep, err := h.Processor.Process(ctx, &get.Request{
		Key: req.Key,
	})

	if err != nil {
		return &grpc.GetReply{
			Status:  1, // TODO: add error codes
			Message: err.Error(),
		}, nil
	}

	return &grpc.GetReply{
		Val: rep.Val,
	}, nil
}
