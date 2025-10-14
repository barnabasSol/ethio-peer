package pagination

import (
	"strconv"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Pagination struct {
	Page     int64
	PageSize int64
}

func New(page, page_size string) *Pagination {
	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}
	ps, err := strconv.Atoi(page_size)
	if err != nil {
		ps = 20
	}
	if p < 1 {
		p = 1
	}
	if ps < 1 || ps > 20 {
		ps = 20
	}
	return &Pagination{
		PageSize: int64(ps),
		Page:     int64(p),
	}
}

func (p *Pagination) GetOptions() *options.FindOptionsBuilder {
	l := p.PageSize
	skip := p.Page*p.PageSize - p.PageSize
	builder := &options.FindOptionsBuilder{}
	return builder.SetLimit(l).SetSkip(skip)
}
