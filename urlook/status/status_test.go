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

func TestIsStatus(t *testing.T) {
	s := New(0)

	type CodeTest []struct {
		Code    int
		IsValid bool
	}
	tests := []struct {
		Name      string
		Func      func() bool
		CodeTests CodeTest
	}{
		{"IsInfo", s.IsInfo, CodeTest{{100, true}, {200, false}}},
		{"IsSuccess", s.IsSuccess, CodeTest{{200, true}, {300, false}}},
		{"IsRedirect", s.IsRedirect, CodeTest{{300, true}, {400, false}}},
		{"IsClientError", s.IsClientError, CodeTest{{400, true}, {500, false}}},
		{"IsServerError", s.IsServerError, CodeTest{{500, true}, {600, false}}},
	}
	for _, test := range tests {
		for _, cT := range test.CodeTests {
			t.Run(
				fmt.Sprintf("%s: code %d", test.Name, cT.Code),
				func(t *testing.T) {
					s.Code = cT.Code
					isValid := test.Func()
					if isValid != cT.IsValid {
						t.Fatalf("expected status %s %v but got %v\n", test.Name, cT.IsValid, isValid)
					}
				},
			)
		}
	}
}
