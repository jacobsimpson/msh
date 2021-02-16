package builtin

import (
	"fmt"
	"os"
)

type pwd struct{}

func (*pwd) Execute(args []string) int {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get the current working directory: %+v\n", err)
		return 1
	}
	fmt.Printf("%s\n", wd)
	return 0
}

func (*pwd) Name() string { return "pwd" }

func (*pwd) ShortHelp() string { return "pwd" }

func (*pwd) LongHelp() string {
	return "Print the current working directory."
}
