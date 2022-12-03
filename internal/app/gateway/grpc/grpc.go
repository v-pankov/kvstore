package grpc

import (
	"net"

	"google.golang.org/grpc"

	"github.com/vdrpkv/kvstore/internal/app/gateway/grpc/server"

	apiGrpc "github.com/vdrpkv/kvstore/generated/api/grpc"
)

func Serve(listener net.Listener, server server.Server, opts ...grpc.ServerOption) error {
	grpcServer := grpc.NewServer(opts...)
	apiGrpc.RegisterKVStoreServer(grpcServer, server)
	return grpcServer.Serve(listener)
}
