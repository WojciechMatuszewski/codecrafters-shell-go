package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var builtins = map[string]BuiltinCmd{
	"exit": ExitCmd{},
	"echo": EchoCmd{stdout: os.Stdout},
	"type": TypeCmd{stdout: os.Stdout, pathLooker: exec.LookPath},
}

func main() {
	executor := Executor{
		stdin:      os.Stdin,
		stdout:     os.Stdout,
		stderr:     os.Stderr,
		pathLooker: exec.LookPath,
	}

	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, args := parseLine(os.Stdin)
		builtin, foundBuiltIn := builtins[cmd]
		if !foundBuiltIn {
			executor.Run(cmd, args)
		} else {
			builtin.Run(args)
		}

	}
}

func parseLine(reader io.Reader) (string, []string) {
	line, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		panic(err)
	}
	line = strings.TrimSuffix(line, "\n")

	cmd := strings.SplitAfter(line, " ")
	return strings.TrimSuffix(cmd[0], " "), cmd[1:]
}

type BuiltinCmd interface {
	Run([]string)
}

type ExitCmd struct{}

func (ec ExitCmd) Run(args []string) {
	if len(args) != 1 {
		panic(errors.New("invalid length of parameters"))
	}
	code := args[0]

	exitCode, err := strconv.Atoi(code)
	if err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

type EchoCmd struct {
	stdout io.Writer
}

func (ec EchoCmd) Run(args []string) {
	if len(args) == 0 {
		fmt.Fprint(ec.stdout, "\n")
		return
	}

	fmt.Fprintf(ec.stdout, "%s\n", strings.Join(args, ""))
}

type TypeCmd struct {
	stdout     io.Writer
	pathLooker func(string) (string, error)
}

func (tc TypeCmd) Run(args []string) {
	if len(args) != 1 {
		panic(errors.New("invalid length of parameters"))
	}
	cmd := args[0]

	_, exists := builtins[cmd]
	if exists {
		fmt.Fprintf(tc.stdout, "%s is a shell builtin\n", cmd)
		return
	}

	execPath, err := tc.pathLooker(cmd)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Fprintf(tc.stdout, "%s: not found\n", cmd)
			return
		}

		panic(err)
	}

	fmt.Fprintf(tc.stdout, "%s is %s\n", cmd, execPath)
}

type Executor struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	pathLooker func(string) (string, error)
}

func (ex Executor) Run(program string, args []string) {
	pPath, err := ex.pathLooker(program)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Fprintf(os.Stdout, "%s: command not found\n", strings.TrimSuffix(program, "\n"))

			return
		}

		panic(err)
	}

	process := exec.Command(pPath, args...)
	process.Stdin = ex.stdin
	process.Stdout = ex.stdout
	process.Stderr = ex.stderr

	err = process.Run()
	if err != nil {
		panic(err)
	}

}
