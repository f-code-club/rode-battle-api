package pagination

type Response[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

func NewResponse[T any](data []T, total int64, page, limit int) *Response[T] {
	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	if data == nil {
		data = []T{}
	}

	return &Response[T]{
		Items:      data,
		Page:       page,
		PageSize:   limit,
		TotalItems: total,
		TotalPages: totalPages,
	}
}
