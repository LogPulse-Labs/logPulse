package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SortOptions struct {
	Column string
	Order  int
}

type Paginate struct {
	limit int64
	page  int64
	sort  *SortOptions
	total *int
}

type PaginateMeta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type PaginationResponse struct {
	Data any          `json:"data"`
	Meta PaginateMeta `json:"meta"`
}

func NewPaginate(limit, page int, sortOptions *SortOptions) *Paginate {
	return &Paginate{
		limit: int64(limit),
		page:  int64(page),
		sort:  sortOptions,
	}
}

func (p *Paginate) GetPaginatedOptions() *options.FindOptions {
	limit := p.limit
	skip := p.page*p.limit - p.limit

	return options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: p.sort.Column, Value: p.sort.Order}})
}

func (p *Paginate) GetMetaInfo() *PaginateMeta {
	return &PaginateMeta{
		Page:    int(p.page),
		PerPage: int(p.limit),
	}
}
