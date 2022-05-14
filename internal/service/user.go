package service

import (
	"context"
	"github.com/asaskevich/EventBus"
	"github.com/go-kratos/kratos/v2/log"
	api "github.com/newpurr/stock/api/user"
	"github.com/newpurr/stock/internal/biz"
)

type (
	UserService struct {
		api.UnimplementedUserServer
		cu  *biz.UserUsecase
		log *log.Helper
		bus EventBus.BusPublisher
	}
)

// NewUserService .
func NewUserService(cu *biz.UserUsecase, logger log.Logger, bus EventBus.BusPublisher) *UserService {
	return &UserService{cu: cu, log: log.NewHelper(logger), bus: bus}
}

func (h UserService) CreateUser(
	ctx context.Context,
	request *api.CreateUserRequest,
) (*api.CreateUserReply, error) {
	do, err := request.ToUserDO()
	if err != nil {
		return nil, err
	}
	res, err := h.cu.Create(ctx, do)
	if err != nil {
		return nil, err
	}

	reply := api.ToCreateUserReply(res)
	return &reply, err
}

func (h UserService) UpdateUser(
	ctx context.Context,
	request *api.UpdateUserRequest,
) (*api.UpdateUserReply, error) {
	do, err := h.cu.Update(ctx, request.ToUserDO())
	if err != nil {
		return nil, err
	}

	reply := api.ToUpdateUserReply(do)
	return &reply, err
}

func (h UserService) DeleteUser(
	ctx context.Context,
	request *api.DeleteUserRequest,
) (*api.DeleteUserReply, error) {
	err := h.cu.Delete(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	return &api.DeleteUserReply{}, err
}

func (h UserService) FindUser(
	ctx context.Context,
	request *api.GetUserRequest,
) (*api.GetUserReply, error) {
	do, err := h.cu.Get(ctx, request.GetId())
	if err != nil {
		return nil, err
	}

	reply := api.ToGetUserReply(do)
	return &reply, err
}

func (h UserService) PagingSearchUser(
	ctx context.Context,
	request *api.UserPagingSearchRequest,
) (*api.UserPaginationListReply, error) {
	res, paginator, err := h.cu.PagingQueryUser(ctx, request.ToUserPagingSearch())
	if err != nil {
		return nil, err
	}
	var (
		reply = api.From(res, paginator)
	)

	return &reply, err
}
