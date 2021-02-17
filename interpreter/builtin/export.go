package builtin

import (
	"fmt"
	"os"

	iio "github.com/jacobsimpson/msh/interpreter/io"
)

type export struct{}

func (*export) Execute(stdio *iio.IOChannels, args []string) <-chan int {
	defer stdio.In.Close()
	defer stdio.Out.Close()
	defer stdio.Err.Close()

	for _, e := range os.Environ() {
		fmt.Fprintf(stdio.Out.Writer, "%s\n", e)
	}
	return done(0)
}

func (*export) Name() string { return "export" }

func (*export) ShortHelp() string { return "export" }

func (*export) LongHelp() string {
	return "Lists all the environment variables that are defined."
}
