package main

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ErrInvalidParametersLength = errors.New("invalid parameters length")
	ErrInvalidParameters       = errors.New("invalid parameters")
)

type Cmd interface {
	Run([]string) error
}

type ExitCmd struct {
	processExitter func(int)
}

func (ec ExitCmd) Run(args []string) error {
	if len(args) != 1 {
		return ErrInvalidParametersLength
	}
	code := args[0]

	exitCode, err := strconv.Atoi(code)
	if err != nil {
		return ErrInvalidParameters
	}

	ec.processExitter(exitCode)
	return nil
}

type EchoCmd struct {
	stdout io.Writer
}

func (ec EchoCmd) Run(args []string) error {
	if len(args) == 0 {
		fmt.Fprint(ec.stdout, "\n")
		return nil
	}

	fmt.Fprintf(ec.stdout, "%s\n", strings.Join(args, ""))
	return nil
}

type TypeCmd struct {
	stdout io.Writer

	execPathLooker func(string) (string, error)

	builtinFinder func(string) bool
}

func (tc TypeCmd) Run(args []string) error {
	if len(args) != 1 {
		return ErrInvalidParametersLength
	}
	cmd := args[0]

	exists := tc.builtinFinder(cmd)
	if exists {
		fmt.Fprintf(tc.stdout, "%s is a shell builtin\n", cmd)

		return nil
	}

	execPath, err := tc.execPathLooker(cmd)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Fprintf(tc.stdout, "%s: not found\n", cmd)

			return nil
		}

		return err
	}

	fmt.Fprintf(tc.stdout, "%s is %s\n", cmd, execPath)

	return nil
}

type PwdCmd struct {
	stdout io.Writer

	wdGetter func() (string, error)
}

func (pc PwdCmd) Run(args []string) error {
	if len(args) != 0 {
		return ErrInvalidParameters
	}

	wd, err := pc.wdGetter()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(pc.stdout, "%s\n", wd)

	return nil
}

type CdCmd struct {
	dirChanger    func(string) error
	homeDirGetter func() string
}

func (cc CdCmd) Run(args []string) error {
	if len(args) != 1 {
		return ErrInvalidParametersLength
	}

	newDir := strings.ReplaceAll(args[0], "~", cc.homeDirGetter())
	err := cc.dirChanger(newDir)
	if err != nil {
		panic(err)
	}

	return nil
}
