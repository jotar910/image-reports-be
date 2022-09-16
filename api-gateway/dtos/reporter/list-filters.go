package dtos

type ListFilters struct {
	Page  int64 `form:"page"`
	Count int64 `form:"count"`
}
