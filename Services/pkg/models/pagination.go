package models

type Pagination struct {
	Page    int    `query:"page" json:"page"`
	PerPage int    `query:"per_page" json:"per_page"`
	Total   int64  `query:"total" json:"total" swaggerignore:"true"`
	OrderBy string `query:"order_by" json:"order_by"`
	Order   string `query:"sort_by" json:"sort_by"`
}

func (p *Pagination) GetPaginationString() string {
	return "page:" + string(rune(p.Page)) + ":per_page:" + string(rune(p.PerPage)) + ":order_by:" + p.OrderBy + ":order:" + p.Order
}
