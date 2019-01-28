package batch

import (
	"encoding/json"
)

// Error is a json serializable http error
type Error struct {
	Err      string `json:"err"`
	HTTPCode int    `json:"code"`
	IsError  bool   `json:"-"`
}

// NewError creates a new batch.Error
func NewError(err string, code int) *Error {
	return &Error{
		Err:      err,
		HTTPCode: code,
		IsError:  true,
	}
}

// Code returns the http code
func (e *Error) Code() int {
	return e.HTTPCode
}

// IsErr returns a boolean flag
func (e *Error) IsErr() bool {
	return e.IsError
}

// Error will marshal the error to json string
func (e *Error) Error() string {
	errMsg, _ := json.Marshal(e)

	return string(errMsg)
}
