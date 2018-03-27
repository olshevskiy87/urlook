package status

import (
	"net/http"
)

// Status contains http query status info
type Status struct {
	Code int
	Text string
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
		Code: code,
		Text: http.StatusText(code),
	}
}

// String returns corresponding status sign (icon)
func (s *Status) String() string {
	var internalCode int
	switch {
	case s.IsInfo():
		internalCode = Info
	case s.IsSuccess():
		internalCode = Success
	case s.IsRedirect():
		internalCode = Redirect
	case s.IsClientError():
		internalCode = ClientError
	case s.IsServerError():
		internalCode = ServerError
	default:
		internalCode = Unknown
	}
	sign, ok := signs[internalCode]
	if !ok {
		sign = signs[Unknown]
	}
	return sign
}

// IsInfo if response is "Informational"
func (s *Status) IsInfo() bool {
	return 100 <= s.Code && s.Code <= 199
}

// IsSuccess if response is "Success"
func (s *Status) IsSuccess() bool {
	return 200 <= s.Code && s.Code <= 299
}

// IsRedirect if response is "Redirection"
func (s *Status) IsRedirect() bool {
	return 300 <= s.Code && s.Code <= 399
}

// IsClientError if response is "Client errors"
func (s *Status) IsClientError() bool {
	return 400 <= s.Code && s.Code <= 499
}

// IsServerError if response is "Server errors"
func (s *Status) IsServerError() bool {
	return 500 <= s.Code && s.Code <= 599
}
