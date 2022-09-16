package dtos

type QuerySearch struct {
	Ids []uint `json:"ids" binding:"required"`
}
