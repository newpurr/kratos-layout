package biz

import (
	"context"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewUserUsecase,
)

type TransactionManager interface {
	Transaction(context.Context, func(ctx context.Context) error) error
}

type PaginateReq struct {
	CurrentPage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
}
