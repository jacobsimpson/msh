package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jacobsimpson/msh/parser"
)

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
