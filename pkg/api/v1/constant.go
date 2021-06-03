package v1

const (
	PageSize = "pageSize"
	PageNum  = "pageNum"
)

type Page struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
