package main

import (
	"fmt"
	"log"
	"net"

	grpcGateway "github.com/vdrpkv/kvstore/internal/app/gateway/grpc"
	grpcHandlerDelete "github.com/vdrpkv/kvstore/internal/app/gateway/grpc/handler/delete"
	grpcHandlerGet "github.com/vdrpkv/kvstore/internal/app/gateway/grpc/handler/get"
	grpcHandlerSet "github.com/vdrpkv/kvstore/internal/app/gateway/grpc/handler/set"
	grpcServer "github.com/vdrpkv/kvstore/internal/app/gateway/grpc/server"
	"github.com/vdrpkv/kvstore/internal/pkg/memcachedpool"

	infraRepoInmem "github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem"
	infraRepoMemcached "github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached"

	itemEntity "github.com/vdrpkv/kvstore/internal/core/entity/item"
	itemUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item"

	usecaseItemDelete "github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
	usecaseItemGet "github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
	usecaseItemSet "github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
)

func main() {
	DoMain()
}

func DoMain() {

	memcachedPool := memcachedpool.NewDynamicClientPool(
		memcachedpool.TCPConnFactory{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 11211,
		},
		7,
	)

	var (
		itemKeyValidator = itemEntity.NewKeyValidator(10, " \r\n\t")

		memcachedRepo = infraRepoMemcached.New(memcachedpool.MemcachedClient{ClientPool: memcachedPool})
		inmemRepo     = infraRepoInmem.New()
		_             = inmemRepo

		repoAdapters = memcachedRepo

		grpcServer = grpcServer.Server{
			Handlers: grpcServer.Handlers{
				Delete: grpcHandlerDelete.Handler{
					Processor: itemUsecase.NewKeyValidationProcessor[
						*usecaseItemDelete.Request,
						usecaseItemDelete.Response,
					](
						itemKeyValidator, usecaseItemDelete.Processor{
							Gateways: usecaseItemDelete.Gateways{
								Repository: repoAdapters.UseCaseItemDeleteAdapter(),
							},
						},
					),
				},
				Get: grpcHandlerGet.Handler{
					Processor: itemUsecase.NewKeyValidationProcessor[
						*usecaseItemGet.Request,
						usecaseItemGet.Response,
					](
						itemKeyValidator, usecaseItemGet.Processor{
							Gateways: usecaseItemGet.Gateways{
								Repository: repoAdapters.UseCaseItemGetAdapter(),
							},
						},
					),
				},
				Set: grpcHandlerSet.Handler{
					Processor: itemUsecase.NewKeyValidationProcessor[
						*usecaseItemSet.Request,
						usecaseItemSet.Response,
					](
						itemKeyValidator, usecaseItemSet.Processor{
							Gateways: usecaseItemSet.Gateways{
								Repository: repoAdapters.UseCaseItemSetAdapter(),
							},
						},
					),
				},
			},
		}
	)

	const port = 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcGateway.Serve(lis, grpcServer); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
