package command

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/jacobsimpson/msh/parser"
)

func ExecuteProgram(stdin io.Reader, stdout, stderr io.Writer, command *parser.Exec) {
	cmd := exec.Command(command.Name, command.Arguments...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintf(stderr, "msh: command not found: %s\n", command.Name)
		}
	}
}
