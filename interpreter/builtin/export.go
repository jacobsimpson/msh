package builtin

import (
	"fmt"
	"io"
	"os"
)

type export struct{}

func (*export) Execute(stdin io.ReadCloser, stdout, stderr io.WriteCloser, args []string) <-chan int {
	defer stdout.Close()
	defer stderr.Close()

	for _, e := range os.Environ() {
		fmt.Fprintf(stdout, "%s\n", e)
	}
	return done(0)
}

func (*export) Name() string { return "export" }

func (*export) ShortHelp() string { return "export" }

func (*export) LongHelp() string {
	return "Lists all the environment variables that are defined."
}
