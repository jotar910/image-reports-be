package dtos

import (
	"errors"
	"net/http"
)

type ErrorOutbound struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func NewError(status int, err string) *ErrorOutbound {
	return &ErrorOutbound{
		Status: status,
		Error:  err,
	}
}

func NewInternalError(err string) *ErrorOutbound {
	return NewError(http.StatusInternalServerError, err)
}

func (e ErrorOutbound) String() string {
	return e.Error
}

func (e *ErrorOutbound) ToError() error {
	return errors.New(e.Error)
}
