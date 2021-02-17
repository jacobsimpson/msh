package builtin

import (
	"fmt"
	"io"
	"os"
)

type pwd struct{}

func (*pwd) Execute(stdin io.ReadCloser, stdout, stderr io.WriteCloser, args []string) <-chan int {
	defer stdout.Close()
	defer stderr.Close()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stderr, "Unable to get the current working directory: %+v\n", err)
		return done(1)
	}
	fmt.Fprintf(stdout, "%s\n", wd)
	return done(0)
}

func (*pwd) Name() string { return "pwd" }

func (*pwd) ShortHelp() string { return "pwd" }

func (*pwd) LongHelp() string {
	return "Print the current working directory."
}
