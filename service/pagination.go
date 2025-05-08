package service

type Pagination struct {
	CurrentPage int
	TotalPages  int
	PerPage     int
}

func NewPagination(totalItems, perPage int) *Pagination {
	totalPages := totalItems / perPage
	if totalItems%perPage != 0 {
		totalPages++
	}
	return &Pagination{
		CurrentPage: 1,
		TotalPages:  totalPages,
		PerPage:     perPage,
	}
}
