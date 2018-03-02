package urlook

import (
	"fmt"

	"github.com/olshevskiy87/urlook/urlook/status"
)

// Result contains an URL check info
type Result struct {
	URL     string
	Message string
	Status  *status.Status // unknown status (0) by default
}

// String returns string representation of the Result
func (r *Result) String() string {
	var statusText string
	if r.Status.Code != 0 {
		if r.Status.Text != "" {
			statusText = fmt.Sprintf("%d, %s", r.Status.Code, r.Status.Text)
		} else {
			statusText = string(r.Status.Code)
		}
	}
	outMessage := r.URL
	if statusText != "" {
		outMessage = fmt.Sprintf("%s [%s]", outMessage, statusText)
	}
	if status.IsRedirect(r.Status.Code) && r.Message != "" {
		return fmt.Sprintf("%s -> %s", outMessage, r.Message)
	}
	if r.Message == "" {
		return outMessage
	}
	return fmt.Sprintf("%s: %s", outMessage, r.Message)
}
