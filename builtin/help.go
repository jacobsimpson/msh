package builtin

import (
	"fmt"
	"sort"
)

type help struct{}

func (*help) Execute([]string) int {
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
