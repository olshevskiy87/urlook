package status

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		Code int
		Sign string
	}{
		{0, signs[Unknown]},
		{100, signs[Info]},
		{200, signs[Success]},
		{300, signs[Redirect]},
		{400, signs[ClientError]},
		{500, signs[ServerError]},
		{600, signs[Unknown]},
	}
	for _, test := range tests {
		t.Run(
			fmt.Sprintf("code %d", test.Code),
			func(t *testing.T) {
				s := New(test.Code)
				sign := s.String()
				if sign != test.Sign {
					t.Fatalf("expected status sign \"%s\" but got \"%s\"\n", test.Sign, sign)
				}
			},
		)
	}
}
