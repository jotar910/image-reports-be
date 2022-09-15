package dtos

type PageableList[T any] struct {
	Content          []T   `json:"content"`
	Page             int64 `json:"page"`
	TotalPages       int64 `json:"totalPages"`
	TotalElements    int64 `json:"totalElements"`
	NumberOfElements int64 `json:"numberOfElements"`
}
