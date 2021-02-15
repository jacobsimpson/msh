package builtin

import (
	"fmt"
	"os"
	"strconv"
)

type exit struct{}

func (e *exit) Execute(args []string) int {
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "msh: exit: too many arguments\n")
		return 1
	}
	if len(args) == 1 {
		status, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "msh: exit: %s: numeric argument required\n", args[0])
			os.Exit(255)
		}
		os.Exit(status)
	}
	os.Exit(0)
	return 0
}

func (*exit) Name() string { return "exit" }

func (*exit) ShortHelp() string { return "exit [status]" }
