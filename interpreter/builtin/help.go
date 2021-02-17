package builtin

import (
	"fmt"
	"io"
	"sort"
)

type help struct{}

func (*help) Execute(stdin io.ReadCloser, stdout, stderr io.WriteCloser, args []string) <-chan int {
	if len(args) == 0 {
		return printHelpSummary(stdin, stdout, stderr)
	}
	status := 1
	for _, arg := range args {
		b := builtins[arg]
		if b == nil {
			fmt.Fprintf(stderr, "msh: help: no help topics match '%s'.\n", arg)
			continue
		}
		fmt.Printf("%s: %s\n", arg, b.ShortHelp())
		for _, line := range wrap(b.LongHelp()) {
			fmt.Fprintf(stdout, "    %s\n", line)
		}
		status = 0
	}
	return done(status)
}

func (*help) Name() string { return "help" }

func (*help) ShortHelp() string { return "help [pattern]" }

func (*help) LongHelp() string {
	return "Display helpful information about builtin commands. If PATTERN is specified, gives detailed help on all commands matching PATTERN, otherwise a list of the builtins is printed."
}

func printHelpSummary(stdin io.Reader, stdout, stderr io.Writer) <-chan int {
	names := []string{}
	for _, b := range builtins {
		names = append(names, b.Name())
	}
	sort.Strings(names)

	fmt.Fprintf(stdout, "msh, version %s\n", Version)
	fmt.Fprintf(stdout, "These shell commands are defined internally. Type `help` to see this list.\n")
	fmt.Fprintf(stdout, "\n")
	for _, name := range names {
		fmt.Fprintf(stdout, "%-10s %s\n", name, builtins[name].ShortHelp())
	}
	fmt.Fprintf(stdout, "\n")
	return done(0)
}
