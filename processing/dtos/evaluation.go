package dtos

type Evaluation struct {
	ReportID   uint     `json:"reportId"`
	ImageID    string   `json:"imageId"`
	Grade      int      `json:"grade"`
	Categories []string `json:"categories"`
}
