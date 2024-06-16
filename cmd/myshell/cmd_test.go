package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"testing"
)

func TestExit(t *testing.T) {
	t.Run("Fails with multiple parameters", func(t *testing.T) {
		cmd := ExitCmd{processExitter: func(i int) {}}

		err := cmd.Run([]string{"1", "2"})
		if !errors.Is(err, ErrInvalidParametersLength) {
			t.Errorf("got %v, wanted %v", err, ErrInvalidParametersLength)
		}
	})

	t.Run("Fails with malformed parameter", func(t *testing.T) {
		cmd := ExitCmd{processExitter: func(i int) {}}

		err := cmd.Run([]string{"a"})
		if !errors.Is(err, ErrInvalidParameters) {
			t.Errorf("got %v, wanted %v", err, ErrInvalidParameters)
		}
	})

	t.Run("Succeeds", func(t *testing.T) {
		var calledWithCode int
		providedCode := 10

		cmd := ExitCmd{processExitter: func(i int) { calledWithCode = i }}

		err := cmd.Run([]string{fmt.Sprintf("%v", providedCode)})
		if err != nil {
			t.Error("did not expect an error")
		}

		if calledWithCode != providedCode {
			t.Errorf("got %v, wanted %v", calledWithCode, providedCode)
		}
	})
}

func TestEcho(t *testing.T) {
	t.Run("Handles single parameter", func(t *testing.T) {
		writer := bytes.NewBuffer([]byte{})

		cmd := EchoCmd{stdout: writer}

		err := cmd.Run([]string{"Hello"})
		if err != nil {
			t.Error("did not expect an error")
		}

		got := writer.String()
		want := "Hello\n"
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})

	t.Run("Handles multiple parameters", func(t *testing.T) {
		stdout := bytes.NewBuffer([]byte{})

		cmd := EchoCmd{stdout: stdout}

		err := cmd.Run([]string{"Hello ", "To ", "You"})
		if err != nil {
			t.Error("did not expect an error")
		}

		got := stdout.String()
		want := "Hello To You\n"
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})
}

func TestType(t *testing.T) {
	t.Run("Fails with multiple parameters", func(t *testing.T) {
		stdout := bytes.NewBuffer([]byte{})

		cmd := TypeCmd{stdout: stdout}

		err := cmd.Run([]string{"1", "2"})
		if !errors.Is(err, ErrInvalidParametersLength) {
			t.Errorf("got %v, wanted %v", err, ErrInvalidParametersLength)
		}
	})

	t.Run("Prints when builtin is found", func(t *testing.T) {
		stdout := bytes.NewBuffer([]byte{})

		providedArg := "existingBuiltin"

		cmd := TypeCmd{
			stdout: stdout,
			builtinFinder: func(s string) bool {
				return s == providedArg
			},
		}

		err := cmd.Run([]string{providedArg})
		if err != nil {
			t.Error("did not expect an error")
		}

		got := stdout.String()
		want := fmt.Sprintf("%s is a shell builtin\n", providedArg)
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})

	t.Run("Fails if looking up executable fails", func(t *testing.T) {
		stdout := bytes.NewBuffer([]byte{})

		providedArg := "existingBuiltin"
		expectedError := errors.New("test error")

		cmd := TypeCmd{
			stdout: stdout,
			builtinFinder: func(s string) bool {
				return false
			},
			execPathLooker: func(s string) (string, error) {
				return "", expectedError
			},
		}

		err := cmd.Run([]string{providedArg})
		if err == nil {
			t.Error("expected an error")
		}

		if !errors.Is(err, expectedError) {
			t.Errorf("got %v, expected %v", err, expectedError)
		}
	})

	t.Run("Prints when builtin and executable is not found", func(t *testing.T) {
		stdout := bytes.NewBuffer([]byte{})

		providedArg := "existingBuiltin"

		cmd := TypeCmd{
			stdout: stdout,
			builtinFinder: func(s string) bool {
				return false
			},
			execPathLooker: func(s string) (string, error) {
				return "", exec.ErrNotFound
			},
		}

		err := cmd.Run([]string{providedArg})
		if err != nil {
			t.Error("did not expect an error")
		}

		got := stdout.String()
		want := fmt.Sprintf("%s: not found\n", providedArg)
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})

	t.Run("Prints when executable is found", func(t *testing.T) {
		stdout := bytes.NewBuffer([]byte{})

		providedArg := "existingBuiltin"
		pathToProvided := "/some/imaginary/path"

		cmd := TypeCmd{
			stdout: stdout,
			builtinFinder: func(s string) bool {
				return false
			},
			execPathLooker: func(s string) (string, error) {
				if s == providedArg {
					return pathToProvided, nil
				}

				return "", exec.ErrNotFound
			},
		}

		err := cmd.Run([]string{providedArg})
		if err != nil {
			t.Error("did not expect an error")
		}

		got := stdout.String()
		want := fmt.Sprintf("%s is %s\n", providedArg, pathToProvided)
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})
}
