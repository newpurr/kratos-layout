package repositories

import (
	"github.com/google/wire"
	"github.com/newpurr/stock/internal/data/repositories/user"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	user.NewUserCacheRepo,
)
