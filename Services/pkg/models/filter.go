package models

type Search struct {
	Keyword string `query:"keyword" json:"keyword"`
	Column  string `query:"column" json:"column"`
	Text    string `query:"text" json:"text"` // For Semantic Search
}

func (s *Search) GetSearchString() string {
	return "keyword=" + s.Keyword + "&column=" + s.Column
}
