package exitcode_test

import (
	"errors"
	"testing"

	"github.com/synchestra-io/specscore/pkg/exitcode"
)

func TestErrorSatisfiesInterface(t *testing.T) {
	type exitCoder interface{ ExitCode() int }
	var err error = exitcode.New(1, "test")
	var ec exitCoder
	if !errors.As(err, &ec) {
		t.Fatal("exitcode.Error does not satisfy exitCoder interface")
	}
	if ec.ExitCode() != 1 {
		t.Fatalf("expected exit code 1, got %d", ec.ExitCode())
	}
}

func TestErrorMessage(t *testing.T) {
	err := exitcode.New(2, "bad args")
	if err.Error() != "bad args" {
		t.Fatalf("expected 'bad args', got %q", err.Error())
	}
}

func TestNewf(t *testing.T) {
	err := exitcode.Newf(10, "failed: %s", "disk full")
	if err.Error() != "failed: disk full" {
		t.Fatalf("expected 'failed: disk full', got %q", err.Error())
	}
	if err.ExitCode() != 10 {
		t.Fatalf("expected exit code 10, got %d", err.ExitCode())
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *exitcode.Error
		code int
	}{
		{"Conflict", exitcode.ConflictError("c"), exitcode.Conflict},
		{"ConflictF", exitcode.ConflictErrorf("c %d", 1), exitcode.Conflict},
		{"InvalidArgs", exitcode.InvalidArgsError("a"), exitcode.InvalidArgs},
		{"InvalidArgsF", exitcode.InvalidArgsErrorf("a %d", 2), exitcode.InvalidArgs},
		{"NotFound", exitcode.NotFoundError("n"), exitcode.NotFound},
		{"NotFoundF", exitcode.NotFoundErrorf("n %d", 3), exitcode.NotFound},
		{"InvalidState", exitcode.InvalidStateError("s"), exitcode.InvalidState},
		{"InvalidStateF", exitcode.InvalidStateErrorf("s %d", 4), exitcode.InvalidState},
		{"Unexpected", exitcode.UnexpectedError("u"), exitcode.Unexpected},
		{"UnexpectedF", exitcode.UnexpectedErrorf("u %d", 10), exitcode.Unexpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.ExitCode() != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, tt.err.ExitCode())
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if exitcode.Success != 0 {
		t.Error("Success should be 0")
	}
	if exitcode.Conflict != 1 {
		t.Error("Conflict should be 1")
	}
	if exitcode.InvalidArgs != 2 {
		t.Error("InvalidArgs should be 2")
	}
	if exitcode.NotFound != 3 {
		t.Error("NotFound should be 3")
	}
	if exitcode.InvalidState != 4 {
		t.Error("InvalidState should be 4")
	}
	if exitcode.Unexpected != 10 {
		t.Error("Unexpected should be 10")
	}
}
