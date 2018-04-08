package urlook

import (
	"testing"
	"time"
)

func TestSetTimeout(t *testing.T) {
	b := New([]string{})
	t.Run(
		"check errors",
		func(t *testing.T) {
			tests := []struct {
				timeout int
				ok      bool
			}{
				{-1, false},
				{0, true},
				{10, true},
			}
			for _, test := range tests {
				err := b.SetTimeout(test.timeout)
				if test.ok && err != nil {
					t.Fatalf("expected no errors but got error \"%v\"", err)
				} else if !test.ok && err == nil {
					t.Fatalf("expected error but got nil")
				}
			}
		},
	)
	t.Run(
		"check timeout duration",
		func(t *testing.T) {
			for _, timeout := range []int{0, 10} {
				err := b.SetTimeout(timeout)
				if err != nil {
					t.Fatalf("expected no errors but got error \"%v\"", err)
				}
				expectedDuration := time.Duration(time.Duration(timeout) * time.Second)
				if b.clientHTTP.Timeout != expectedDuration {
					t.Fatalf(
						"expected timeout duration %v but got %v",
						expectedDuration,
						b.clientHTTP.Timeout,
					)
				}
			}
		},
	)
}

func TestSetWorkersCount(t *testing.T) {
	b := New([]string{})
	tests := []struct {
		workers int
		ok      bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}
	for _, test := range tests {
		b.SetWorkersCount(test.workers)
		if test.ok && b.workersCount != test.workers {
			t.Fatalf("expected workers count %v but got %v", test.workers, b.workersCount)
		} else if !test.ok && b.workersCount == test.workers {
			t.Fatalf("expected workers count differs but got equal values: %v", b.workersCount)
		}
	}
}

func TestSetIsFailOnDuplicates(t *testing.T) {
	b := New([]string{})
	for _, expected := range []bool{true, false} {
		b.SetIsFailOnDuplicates(expected)
		if b.isFailOnDuplicates != expected {
			t.Fatalf("expected isFailOnDuplicates %v but got %v", expected, b.isFailOnDuplicates)
		}
	}
}
