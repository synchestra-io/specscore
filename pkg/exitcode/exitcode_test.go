package exitcode

import (
	"testing"
)

func TestError(t *testing.T) {
	err := New(InvalidArgs, "bad input")
	if err.Error() != "bad input" {
		t.Errorf("expected 'bad input', got %q", err.Error())
	}
	if err.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d", err.ExitCode())
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		code int
	}{
		{"ConflictError", ConflictError("c"), Conflict},
		{"InvalidArgsError", InvalidArgsError("i"), InvalidArgs},
		{"NotFoundError", NotFoundError("n"), NotFound},
		{"InvalidStateError", InvalidStateError("s"), InvalidState},
		{"UnexpectedError", UnexpectedError("u"), Unexpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.ExitCode() != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, tt.err.ExitCode())
			}
		})
	}
}
