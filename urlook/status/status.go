package status

import (
	"net/http"
)

// Status contains http query status info
type Status struct {
	Code         int
	internalCode int
	Text         string
}

// internal http status codes
const (
	Unknown = iota
	Info
	Success
	Redirect
	ClientError
	ServerError
)

// signs is a map with statuses icons
var signs = map[int]string{
	Info:        "i",
	Success:     "✓",
	Redirect:    "→",
	ClientError: "x",
	ServerError: "X",
	Unknown:     "?",
}

// New returns new Status object
func New(code int) *Status {
	return &Status{
		Code:         code,
		internalCode: getInternalCode(code),
		Text:         http.StatusText(code),
	}
}

// GetSign returns corresponding status sign (icon)
func (s *Status) GetSign() string {
	var res string
	if sign, ok := signs[s.internalCode]; ok {
		res = sign
	} else {
		res = signs[Unknown]
	}
	return res
}

// getInternalCode returns corresponding internal
// status code by http code
func getInternalCode(code int) int {
	var res int
	switch {
	case IsInfo(code):
		res = Info
	case IsSuccess(code):
		res = Success
	case IsRedirect(code):
		res = Redirect
	case IsClientError(code):
		res = ClientError
	case IsServerError(code):
		res = ServerError
	default:
		res = Unknown
	}
	return res
}

// TODO: change follow funcs as Status' methods

// IsInfo if response is "Informational"
func IsInfo(code int) bool {
	return 100 <= code && code <= 199
}

// IsSuccess if response is "Success"
func IsSuccess(code int) bool {
	return 200 <= code && code <= 299
}

// IsRedirect if response is "Redirection"
func IsRedirect(code int) bool {
	return 300 <= code && code <= 399
}

// IsClientError if response is "Client errors"
func IsClientError(code int) bool {
	return 400 <= code && code <= 499
}

// IsServerError if response is "Server errors"
func IsServerError(code int) bool {
	return 500 <= code && code <= 599
}
