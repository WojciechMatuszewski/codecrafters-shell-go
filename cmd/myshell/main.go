package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

var builtins = map[string]Cmd{
	"exit": ExitCmd{},
	"echo": EchoCmd{writer: os.Stdout},
	"type": TypeCmd{writer: os.Stdout, execGlobFS: ExecGlobFS{}},
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
	execGlobFS fs.GlobFS
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

	matches, err := tc.execGlobFS.Glob(cmd)
	if err != nil {
		panic(err)
	}

	if len(matches) == 0 {
		fmt.Fprintf(tc.writer, "%s: not found\n", cmd)
		return
	}

	fmt.Fprintf(tc.writer, "%s is %s\n", cmd, matches[0])
}

type ExecGlobFS struct {
	fs.FS
}

func (eg ExecGlobFS) Glob(cmd string) ([]string, error) {
	p, found := os.LookupEnv("PATH")
	if !found {
		return []string{}, errors.New("PATH environment variable not found")
	}

	paths := strings.Split(p, ":")
	for _, path := range paths {
		fsys := os.DirFS(path)

		matches, err := fs.Glob(fsys, cmd)
		if err != nil {
			return []string{}, err
		}

		if len(matches) > 1 {
			panic(errors.New("found multiple executables"))
		}

		if len(matches) > 0 {
			match := fmt.Sprintf("%s/%s", path, matches[0])
			return []string{match}, nil
		}
	}

	return []string{}, nil
}
