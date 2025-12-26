package models

import "fmt"

type Pagination struct {
	Page    int    `query:"page" json:"page"`
	PerPage int    `query:"per_page" json:"per_page"`
	Total   int64  `query:"total" json:"total" swaggerignore:"true"`
	OrderBy string `query:"order_by" json:"order_by"`
	Order   string `query:"sort_by" json:"sort_by"`
}

func (p *Pagination) GetPaginationString() string {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 || p.PerPage > 1000 {
		p.PerPage = 10
	}
	return fmt.Sprintf("page=%d&per_page=%d&order_by=%s&order=%s", p.Page, p.PerPage, p.OrderBy, p.Order)
}
