package inmem

import (
	"github.com/vdrpkv/kvstore/internal/app/infra/repository"
	"github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/state"

	deleteAdapter "github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/usecase/item/delete"
	getAdapter "github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/usecase/item/get"
	setAdapter "github.com/vdrpkv/kvstore/internal/app/infra/repository/inmem/usecase/item/set"

	deleteUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
	getUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
	setUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
)

type repo struct {
	state state.State
}

func New() repository.Adapters {
	return &repo{}
}

func (r *repo) UseCaseItemDeleteAdapter() deleteUsecase.Repository {
	return deleteAdapter.Adapter{State: &r.state}
}

func (r *repo) UseCaseItemGetAdapter() getUsecase.Repository {
	return getAdapter.Adapter{State: &r.state}
}

func (r *repo) UseCaseItemSetAdapter() setUsecase.Repository {
	return setAdapter.Adapter{State: &r.state}
}
