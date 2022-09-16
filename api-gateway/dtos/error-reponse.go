package dtos

type ErrorResponse struct {
	Error string `json:"error"`
}

func (e ErrorResponse) String() string {
	return e.Error
}
