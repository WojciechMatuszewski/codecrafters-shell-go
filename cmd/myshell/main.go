package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	builtins := map[string]Cmd{
		"exit": ExitCmd{
			processExitter: os.Exit,
		},
		"echo": EchoCmd{
			stdout: os.Stdout,
		},
		"type": TypeCmd{
			stdout: os.Stdout,

			execPathLooker: exec.LookPath,

			builtinFinder: func(s string) bool {
				switch {
				case s == "exit":
					return true
				case s == "echo":
					return true
				case s == "type":
					return true
				default:
					return false
				}
			},
		},
		"pwd": PwdCmd{
			stdout:   os.Stdout,
			wdGetter: os.Getwd,
		},
		"cd": CdCmd{
			dirChanger:    os.Chdir,
			homeDirGetter: func() string { return os.Getenv("HOME") },
		},
	}

	executor := Executor{
		stdin:          os.Stdin,
		stdout:         os.Stdout,
		stderr:         os.Stderr,
		execPathLooker: exec.LookPath,
	}

	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, args := parseLine(os.Stdin)
		builtin, foundBuiltIn := builtins[cmd]
		if !foundBuiltIn {
			executor.Run(cmd, args)
		} else {
			err := builtin.Run(args)
			if err != nil {
				panic(err)
			}
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

type Executor struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	execPathLooker func(string) (string, error)
}

func (ex Executor) Run(program string, args []string) {
	pPath, err := ex.execPathLooker(program)
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
