package user

import (
	"github.com/newpurr/stock/internal/biz"
	"github.com/newpurr/stock/pkg/utils/pagination"
)

func (x *CreateUserRequest) ToUserDO() (biz.UserDO, error) {

	return biz.UserDO{
		Name:      x.Name,
		Pwd:       x.Pwd,
		IsDeleted: x.IsDeleted,
	}, nil
}

func ToCreateUserReply(sc biz.UserDO) CreateUserReply {
	createdAt := sc.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := sc.UpdatedAt.Format("2006-01-02 15:04:05")

	return CreateUserReply{
		Id:        sc.Id,
		Name:      sc.Name,
		Pwd:       sc.Pwd,
		IsDeleted: sc.IsDeleted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (x *UpdateUserRequest) ToUserDO() biz.UserDO {

	return biz.UserDO{
		Id:        x.Id,
		Name:      x.Name,
		Pwd:       x.Pwd,
		IsDeleted: x.IsDeleted,
	}
}

func ToUpdateUserReply(sc biz.UserDO) UpdateUserReply {
	createdAt := sc.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := sc.UpdatedAt.Format("2006-01-02 15:04:05")

	return UpdateUserReply{
		Id:        sc.Id,
		Name:      sc.Name,
		Pwd:       sc.Pwd,
		IsDeleted: sc.IsDeleted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func ToGetUserReply(sc biz.UserDO) GetUserReply {
	createdAt := sc.CreatedAt.Format("2006-01-02 15:04:05")
	updatedAt := sc.UpdatedAt.Format("2006-01-02 15:04:05")

	return GetUserReply{
		Id:        sc.Id,
		Name:      sc.Name,
		Pwd:       sc.Pwd,
		IsDeleted: sc.IsDeleted,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (p *UserPagingSearchRequest) ToUserPagingSearch() biz.UserPagingSearch {
	return biz.UserPagingSearch{
		UserSearch: &biz.UserSearch{
			Ids: p.Ids,
		},
		PaginateReq: &biz.PaginateReq{
			CurrentPage: p.CurrentPage,
			PageSize:    p.PageSize,
		},
	}
}

func From(result []biz.UserDO, paginator pagination.Paginator) UserPaginationListReply {
	var (
		res []*GetUserReply
	)

	for _, do := range result {
		reply := ToGetUserReply(do)
		res = append(res, &reply)
	}

	return UserPaginationListReply{
		Data: res,
		PageMeta: &UserPaginationListReply_Page{
			TotalPages:  paginator.GetTotalPages(),
			TotalRows:   paginator.GetTotalRows(),
			PageSize:    paginator.GetPageSize(),
			CurrentPage: paginator.GetCurrentPage(),
		},
	}
}
