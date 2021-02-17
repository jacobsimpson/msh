package builtin

import (
	"fmt"
	"sort"

	iio "github.com/jacobsimpson/msh/interpreter/io"
)

type help struct{}

func (*help) Execute(stdio *iio.IOChannels, args []string) <-chan int {
	defer stdio.In.Close()
	defer stdio.Out.Close()
	defer stdio.Err.Close()

	if len(args) == 0 {
		return printHelpSummary(stdio)
	}
	status := 1
	for _, arg := range args {
		b := builtins[arg]
		if b == nil {
			fmt.Fprintf(stdio.Err.Writer, "msh: help: no help topics match '%s'.\n", arg)
			continue
		}
		fmt.Printf("%s: %s\n", arg, b.ShortHelp())
		for _, line := range wrap(b.LongHelp()) {
			fmt.Fprintf(stdio.Out.Writer, "    %s\n", line)
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

func printHelpSummary(stdio *iio.IOChannels) <-chan int {
	names := []string{}
	for _, b := range builtins {
		names = append(names, b.Name())
	}
	sort.Strings(names)

	fmt.Fprintf(stdio.Out.Writer, "msh, version %s\n", Version)
	fmt.Fprintf(stdio.Out.Writer, "These shell commands are defined internally. Type `help` to see this list.\n")
	fmt.Fprintf(stdio.Out.Writer, "\n")
	for _, name := range names {
		fmt.Fprintf(stdio.Out.Writer, "%-10s %s\n", name, builtins[name].ShortHelp())
	}
	fmt.Fprintf(stdio.Out.Writer, "\n")
	return done(0)
}
