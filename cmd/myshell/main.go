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

var builtins = map[string]Cmd{
	"exit": ExitCmd{},
	"echo": EchoCmd{writer: os.Stdout},
	"type": TypeCmd{writer: os.Stdout, pathLooker: exec.LookPath},
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, args := parseLine(os.Stdin)
		toRun, found := builtins[cmd]
		if !found {
			fmt.Fprintf(os.Stdout, "%s: command not found\n", strings.TrimSuffix(cmd, "\n"))
		} else {
			toRun.Run(args)
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

type Cmd interface {
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
	writer io.Writer
}

func (ec EchoCmd) Run(args []string) {
	if len(args) == 0 {
		fmt.Fprint(ec.writer, "\n")
		return
	}

	fmt.Fprintf(ec.writer, "%s\n", strings.Join(args, ""))
}

type TypeCmd struct {
	writer     io.Writer
	pathLooker func(string) (string, error)
}

func (tc TypeCmd) Run(args []string) {
	if len(args) != 1 {
		panic(errors.New("invalid length of parameters"))
	}
	cmd := args[0]

	_, exists := builtins[cmd]
	if exists {
		fmt.Fprintf(tc.writer, "%s is a shell builtin\n", cmd)
		return
	}

	execPath, err := tc.pathLooker(cmd)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Fprintf(tc.writer, "%s: not found\n", cmd)
			return
		}

		panic(err)
	}

	fmt.Fprintf(tc.writer, "%s is %s\n", cmd, execPath)
}

type Executor struct {
	writer     io.Writer
	pathLooker func(string) (string, error)
}

func (ex Executor) Run(program string, args []string) {
	// matches, err := ex.execGlobFS.Glob(program)
	// if err != nil {
	// 	panic(err)
	// }

	// cmd := exec.LookPath()
}
