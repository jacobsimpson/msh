package builtin

import (
	"fmt"
	"os"

	iio "github.com/jacobsimpson/msh/interpreter/io"
)

type pwd struct{}

func (*pwd) Execute(stdio *iio.IOChannels, args []string) <-chan int {
	defer stdio.In.Close()
	defer stdio.Out.Close()
	defer stdio.Err.Close()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(stdio.Err.Writer, "Unable to get the current working directory: %+v\n", err)
		return done(1)
	}
	fmt.Fprintf(stdio.Out.Writer, "%s\n", wd)
	return done(0)
}

func (*pwd) Name() string { return "pwd" }

func (*pwd) ShortHelp() string { return "pwd" }

func (*pwd) LongHelp() string {
	return "Print the current working directory."
}
