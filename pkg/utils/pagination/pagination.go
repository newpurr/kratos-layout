package pagination

import "math"

type PaginateSearcher interface {
	GetCurrentPage() int64
	GetPageSize() int64
}
type Paginator interface {
	PaginateSearcher
	GetTotalRows() int64
	GetTotalPages() int64
}

type paginator struct {
	CurrentPage int64 `json:"current_page"`
	PageSize    int64 `json:"page_size"`
	TotalRows   int64 `json:"total_rows"`
	TotalPages  int64 `json:"total_pages"`
}

func (p paginator) GetCurrentPage() int64 {
	if p.TotalPages <= p.CurrentPage {
		return p.TotalPages
	}
	return p.CurrentPage
}

func (p paginator) GetPageSize() int64 {
	return p.PageSize
}

func (p paginator) GetTotalRows() int64 {
	return p.TotalRows
}

func (p paginator) GetTotalPages() int64 {
	return p.TotalPages
}

func NewPaginator(totalRows, currentPage, pageSize int64) Paginator {
	p := paginator{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		TotalRows:   totalRows,
		TotalPages:  int64(math.Ceil(float64(totalRows) / float64(pageSize))),
	}

	return p
}

func NewSearchPaginator(currentPage, pageSize int64) PaginateSearcher {
	p := paginator{
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}

	return p
}
