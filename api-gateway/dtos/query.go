package dtos

type QueryIds[T any] struct {
	Ids []T `json:"ids"`
}
