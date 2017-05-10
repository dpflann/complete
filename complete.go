// Package complete provides a tool for bash writing bash completion in go.
//
// Writing bash completion scripts is a hard work. This package provides an easy way
// to create bash completion scripts for any command, and also an easy way to install/uninstall
// the completion of the command.
package complete

import (
	"fmt"
	"os"
	"strings"

	"github.com/posener/complete/cmd"
)

const (
	envComplete = "COMP_LINE"
	envDebug    = "COMP_DEBUG"
)

// Complete structs define completion for a command with CLI options
type Complete struct {
	Command Command
	cmd.CLI
}

// New creates a new complete command.
// name is the name of command we want to auto complete.
// IMPORTANT: it must be the same name - if the auto complete
// completes the 'go' command, name must be equal to "go".
// command is the struct of the command completion.
func New(name string, command Command) *Complete {
	return &Complete{
		Command: command,
		CLI:     cmd.CLI{Name: name},
	}
}

// Run get a command, get the typed arguments from environment
// variable, and print out the complete options
// returns success if the completion ran or if the cli matched
// any of the given flags, false otherwise
func (c *Complete) Run() bool {
	args, ok := getLine()
	if !ok {
		// make sure flags parsed,
		// in case they were not added in the main program
		return c.CLI.Run()
	}
	Log("Completing args: %s", args)

	options := complete(c.Command, args)

	Log("Completion: %s", options)
	output(options)
	return true
}

// complete get a command an command line arguments and returns
// matching completion options
func complete(c Command, args []string) (matching []string) {
	options, _ := c.options(args[:len(args)-1])

	// choose only matching options
	l := last(args)
	for _, option := range options {
		if option.Match(l) {
			matching = append(matching, option.String())
		}
	}
	return
}

func getLine() ([]string, bool) {
	line := os.Getenv(envComplete)
	if line == "" {
		return nil, false
	}
	return strings.Split(line, " "), true
}

func last(args []string) (last string) {
	if len(args) > 0 {
		last = args[len(args)-1]
	}
	return
}

func output(options []string) {
	// stdout of program defines the complete options
	for _, option := range options {
		fmt.Println(option)
	}
}