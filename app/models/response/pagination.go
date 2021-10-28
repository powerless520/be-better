package response

type Pagination struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int64       `json:"page"`
	PerPage  int64       `json:"per_page"`
	LastPage int64       `json:"last_page"`
}
