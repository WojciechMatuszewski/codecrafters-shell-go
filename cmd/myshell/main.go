package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"strconv"
	"strings"
)

func main() {
	builtin := map[string]Cmd{
		"exit": ExitCmd{},
		"echo": EchoCmd{writer: os.Stdout},
	}

	meta := map[string]Cmd{
		"type": TypeCmd{
			writer: os.Stdout,
			cmdExists: func(s string) bool {
				_, found := builtin[s]
				return found
			},
		},
	}

	commands := make(map[string]Cmd)
	maps.Copy(commands, builtin)
	maps.Copy(commands, meta)

	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, args := parseLine(os.Stdin)
		toRun, found := commands[cmd]
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
	writer    io.Writer
	cmdExists func(string) bool
}

func (tc TypeCmd) Run(args []string) {
	if len(args) != 1 {
		panic(errors.New("invalid length of parameters"))
	}

	cmd := args[0]
	if cmd == "type" {
		fmt.Fprintf(tc.writer, "%s is a shell builtin\n", cmd)
		return
	}

	if tc.cmdExists(cmd) {
		fmt.Fprintf(tc.writer, "%s is a shell builtin\n", cmd)
		return
	}

	fmt.Fprintf(tc.writer, "%s: not found \n", args[0])
}
