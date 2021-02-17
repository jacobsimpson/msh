package builtin

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type exit struct{}

func (e *exit) Execute(stdin io.ReadCloser, stdout, stderr io.WriteCloser, args []string) int {
	defer stdout.Close()
	defer stderr.Close()

	if len(args) > 1 {
		fmt.Fprintf(stderr, "msh: exit: too many arguments\n")
		return 1
	}
	if len(args) == 1 {
		status, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(stderr, "msh: exit: %s: numeric argument required\n", args[0])
			os.Exit(255)
		}
		os.Exit(status)
	}
	os.Exit(0)
	return 0
}

func (*exit) Name() string { return "exit" }

func (*exit) ShortHelp() string { return "exit [status]" }

func (*exit) LongHelp() string {
	return "Exit the shell with a status of N. If N is omitted, the exit status is that of the last command executed."
}
