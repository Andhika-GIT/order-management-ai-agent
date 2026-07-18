package model

type PaginationRequest struct {
	Page    int `json:"page" query:"page"`
	PerPage int `json:"per_page" query:"per_page"`
}

type Paginated[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}
