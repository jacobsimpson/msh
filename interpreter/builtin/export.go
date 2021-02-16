package builtin

import (
	"fmt"
	"io"
	"os"
)

type export struct{}

func (*export) Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	for _, e := range os.Environ() {
		fmt.Fprintf(stdout, "%s\n", e)
	}
	return 0
}

func (*export) Name() string { return "export" }

func (*export) ShortHelp() string { return "export" }

func (*export) LongHelp() string {
	return "Lists all the environment variables that are defined."
}
