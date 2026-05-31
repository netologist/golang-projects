package main

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// AppError carries an HTTP status, a stable code, a message, and a cause.
type AppError struct {
	Code    string
	Message string
	Status  int
	stack   string
	cause   error
}

// Error renders the message and cause.
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

// Unwrap exposes the cause for errors.Is/As.
func (e *AppError) Unwrap() error { return e.cause }

// Stack returns the capture site (file:line).
func (e *AppError) Stack() string { return e.stack }

// MarshalJSON serialises the public fields for API responses.
func (e *AppError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{e.Code, e.Message, e.Status})
}

// New builds an AppError, capturing the caller's file:line.
func New(status int, code, msg string, cause error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Code:    code,
		Message: msg,
		Status:  status,
		stack:   fmt.Sprintf("%s:%d", file, line),
		cause:   cause,
	}
}
