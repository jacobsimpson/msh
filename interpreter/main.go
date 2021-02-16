package interpreter

import (
	"fmt"
	"os"

	"github.com/jacobsimpson/msh/builtin"
	"github.com/jacobsimpson/msh/command"
	"github.com/jacobsimpson/msh/parser"
)

func Evaluate(program *parser.Program) {
	stdin, stdout, stderr := os.Stdin, os.Stdout, os.Stderr

	cmd := program.Command.(*parser.Exec)

	if cmd.Redirection != nil {
		switch cmd.Redirection.Type {
		case parser.Truncate:
			f, err := os.Create(cmd.Redirection.Target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "msh: %+v", err)
				return
			}
			stdout = f
			defer f.Close()
		case parser.TruncateAll:
			f, err := os.Create(cmd.Redirection.Target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "msh: %+v", err)
				return
			}
			stdout = f
			stderr = f
			defer f.Close()
		case parser.Append:
			f, err := os.OpenFile(cmd.Redirection.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "msh: %+v", err)
				return
			}
			stdout = f
			defer f.Close()
		}
	}

	if cmd.Name == "" {
		// Do nothing.
	} else if c := builtin.Get(cmd.Name); c != nil {
		c.Execute(stdin, stdout, stderr, cmd.Arguments)
	} else {
		command.ExecuteProgram(stdin, stdout, stderr, cmd)
	}
}
