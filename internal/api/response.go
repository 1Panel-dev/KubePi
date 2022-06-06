package api

type Page struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
