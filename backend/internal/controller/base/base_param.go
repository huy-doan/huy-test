package base

type PaginationRequest struct {
	Page      int    `query:"page" json:"page" validate:"omitempty,min=1"`
	PageSize  int    `query:"page_size" json:"page_size" validate:"omitempty,min=1"`
	SortField string `query:"sort_field" json:"sort_field"`
	SortOrder string `query:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	Total      int64 `json:"total"`
}
