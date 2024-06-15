package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	commands := map[string]Cmd{
		"exit": ExitCmd{},
		"echo": EchoCmd{writer: os.Stdout},
	}

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
