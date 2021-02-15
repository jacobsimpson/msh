package command

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"

	"github.com/jacobsimpson/msh/parser"
)

func PWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get the current working directory: %+v\n", err)
		return
	}
	fmt.Printf("%s\n", wd)
}

func Exit() {
	os.Exit(0)
}

func CD(args []string) {
	dst := ""
	if len(args) == 0 {
		usr, err := user.Current()
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to determine home directory: %+v", err)
			return
		}
		dst = usr.HomeDir
	} else {
		dst = args[0]
	}
	if err := os.Chdir(dst); err != nil {
		fmt.Fprintf(os.Stderr, "no such file or directory: %s", dst)
	}
}

func ExecuteProgram(command *parser.Command) {
	cmd := exec.Command(command.Name, command.Arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintf(os.Stderr, "msh: command not found: %s\n", command.Name)
		}
	}
}

func Export() {
	for _, e := range os.Environ() {
		fmt.Printf("%s\n", e)
	}
}
