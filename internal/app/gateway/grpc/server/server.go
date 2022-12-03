package server

import (
	"context"

	"github.com/vdrpkv/kvstore/generated/api/grpc"
)

type (
	Server struct {
		grpc.UnimplementedKVStoreServer

		Handlers Handlers
	}

	Handlers struct {
		Delete Handler[grpc.DeleteRequest, grpc.DeleteReply]
		Get    Handler[grpc.GetRequest, grpc.GetReply]
		Set    Handler[grpc.SetRequest, grpc.SetReply]
	}
)

func (s Server) Delete(ctx context.Context, req *grpc.DeleteRequest) (*grpc.DeleteReply, error) {
	return s.Handlers.Delete.Handle(ctx, req)
}

func (s Server) Get(ctx context.Context, req *grpc.GetRequest) (*grpc.GetReply, error) {
	return s.Handlers.Get.Handle(ctx, req)
}

func (s Server) Set(ctx context.Context, req *grpc.SetRequest) (*grpc.SetReply, error) {
	return s.Handlers.Set.Handle(ctx, req)
}

type Handler[Req any, Rep any] interface {
	Handle(context.Context, *Req) (*Rep, error)
}
