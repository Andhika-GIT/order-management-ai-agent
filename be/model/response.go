package model

import "fmt"

type APIResponse struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// implement error interface
func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// helper constructor
func WriteError(code int, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}
