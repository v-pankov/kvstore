package repository

import (
	deleteUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/delete"
	getUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/get"
	setUsecase "github.com/vdrpkv/kvstore/internal/core/usecase/item/set"
)

type Adapters interface {
	UseCaseItemDeleteAdapter() deleteUsecase.Repository
	UseCaseItemGetAdapter() getUsecase.Repository
	UseCaseItemSetAdapter() setUsecase.Repository
}
