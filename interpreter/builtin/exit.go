package builtin

import (
	"fmt"
	"os"
	"strconv"

	iio "github.com/jacobsimpson/msh/interpreter/io"
)

type exit struct{}

func (e *exit) Execute(stdio *iio.IOChannels, args []string) <-chan int {
	defer stdio.In.Close()
	defer stdio.Out.Close()
	defer stdio.Err.Close()

	if len(args) > 1 {
		fmt.Fprintf(stdio.Err.Writer, "msh: exit: too many arguments\n")
		return done(1)
	}
	if len(args) == 1 {
		status, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(stdio.Err.Writer, "msh: exit: %s: numeric argument required\n", args[0])
			os.Exit(255)
		}
		os.Exit(status)
	}
	os.Exit(0)
	return done(0)
}

func (*exit) Name() string { return "exit" }

func (*exit) ShortHelp() string { return "exit [status]" }

func (*exit) LongHelp() string {
	return "Exit the shell with a status of N. If N is omitted, the exit status is that of the last command executed."
}
