package gormutils

import (
	"github.com/newpurr/stock/pkg/utils/pagination"
	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, r pagination.PaginateSearcher, fn func(db *gorm.DB, paginator pagination.Paginator) error) error {
	var (
		totalRows int64
		page      int64
		pageSize  int64
	)

	page = r.GetCurrentPage()
	if page == 0 {
		page = 1
	}

	pageSize = r.GetPageSize()
	switch {
	case pageSize > 500:
		pageSize = 500
	case pageSize <= 0:
		pageSize = 10
	}

	err := db.Count(&totalRows).Error
	if err != nil {
		return err
	}

	paginator := pagination.NewPaginator(totalRows, page, pageSize)

	offset := (paginator.GetCurrentPage() - 1) * pageSize
	db.Offset(int(offset)).Limit(int(pageSize))
	return fn(db, paginator)
}
