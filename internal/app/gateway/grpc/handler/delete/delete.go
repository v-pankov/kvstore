package delete

import (
	"context"

	"github.com/vdrpkv/kvstore/generated/api/grpc"
	"github.com/vdrpkv/kvstore/internal/app/gateway/grpc/server"

	"github.com/vdrpkv/kvstore/internal/core/usecase"
	"github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
)

type Handler struct {
	Processor usecase.Processor[*delete.Request, *delete.Response]
}

var _ server.Handler[grpc.DeleteRequest, grpc.DeleteReply] = Handler{}

func (h Handler) Handle(ctx context.Context, req *grpc.DeleteRequest) (*grpc.DeleteReply, error) {
	_, err := h.Processor.Process(ctx, &delete.Request{
		Key: req.Key,
	})

	if err != nil {
		return &grpc.DeleteReply{
			Status:  1, // TODO: add error codes
			Message: err.Error(),
		}, nil
	}

	return &grpc.DeleteReply{}, nil
}
