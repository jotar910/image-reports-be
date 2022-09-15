package dtos

type ReportCreation struct {
	Name     string `json:"name" binding:"required,max=255"`
	Callback string `json:"callback" binding:"required,max=2048"`
	ImageID  string `json:"imageId" binding:"required,max=2048"`
}
