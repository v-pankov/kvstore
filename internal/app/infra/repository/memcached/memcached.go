package memcached

import (
	"github.com/vdrpkv/kvstore/internal/pkg/memcached"

	"github.com/vdrpkv/kvstore/internal/app/infra/repository"
	deleteAdapter "github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/usecase/item/delete"
	getAdapter "github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/usecase/item/get"
	setAdapter "github.com/vdrpkv/kvstore/internal/app/infra/repository/memcached/usecase/item/set"

	deleteUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
	getUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
	setUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
)

type repo struct {
	client memcached.Client
}

func New(client memcached.Client) repository.Adapters {
	return &repo{
		client: client,
	}
}

func (r *repo) UseCaseItemDeleteAdapter() deleteUsecase.Repository {
	return deleteAdapter.Adapter{Client: r.client}
}

func (r *repo) UseCaseItemGetAdapter() getUsecase.Repository {
	return getAdapter.Adapter{Client: r.client}
}

func (r *repo) UseCaseItemSetAdapter() setUsecase.Repository {
	return setAdapter.Adapter{Client: r.client}
}
