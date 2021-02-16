package builtin

import (
	"fmt"
	"io"
	"os"
)

type pwd struct{}

func (*pwd) Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) int {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "Unable to get the current working directory: %+v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "%s\n", wd)
	return 0
}

func (*pwd) Name() string { return "pwd" }

func (*pwd) ShortHelp() string { return "pwd" }

func (*pwd) LongHelp() string {
	return "Print the current working directory."
}
