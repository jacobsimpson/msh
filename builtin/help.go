package builtin

import (
	"fmt"
	"sort"
)

type help struct{}

func (*help) Execute(args []string) int {
	if len(args) == 0 {
		return printHelpSummary()
	}
	status := 1
	for _, arg := range args {
		b := builtins[arg]
		if b == nil {
			fmt.Printf("msh: help: no help topics match '%s'.\n", arg)
			continue
		}
		fmt.Printf("%s: %s\n", arg, b.ShortHelp())
		status = 0
	}
	return status
}

func printHelpSummary() int {
	names := []string{}
	for _, b := range builtins {
		names = append(names, b.Name())
	}
	sort.Strings(names)

	fmt.Printf("These shell commands are defined internally. Type `help` to see this list.\n")
	fmt.Println()
	for _, name := range names {
		fmt.Printf("%-10s %s\n", name, builtins[name].ShortHelp())
	}
	fmt.Println()
	return 0
}

func (*help) Name() string { return "help" }

func (*help) ShortHelp() string { return "help" }
