package interpreter

import (
	"fmt"
	"io"
	"os"

	"github.com/jacobsimpson/msh/interpreter/builtin"
	"github.com/jacobsimpson/msh/interpreter/command"
	"github.com/jacobsimpson/msh/parser"
)

type stdio struct {
	in  io.ReadCloser
	out io.WriteCloser
	err io.WriteCloser
}

func Evaluate(program *parser.Program) {
	stdio := &stdio{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
	evaluate(stdio, program.Command)
}

func evaluate(stdio *stdio, cmd parser.Command) {
	switch c := cmd.(type) {
	case *parser.Exec:
		evaluateExec(stdio, c)
	case *parser.Redirection:
		evaluateRedirection(stdio, c)
	}
}

func evaluateExec(stdio *stdio, cmd *parser.Exec) {
	if cmd.Name == "" {
		// Do nothing.
	} else if c := builtin.Get(cmd.Name); c != nil {
		c.Execute(stdio.in, stdio.out, stdio.err, cmd.Arguments)
	} else {
		command.ExecuteProgram(stdio.in, stdio.out, stdio.err, cmd)
	}
}

func evaluateRedirection(stdio *stdio, cmd *parser.Redirection) {
	switch cmd.Type {
	case parser.Truncate:
		f, err := os.Create(cmd.Target)
		if err != nil {
			fmt.Fprintf(stdio.err, "msh: %+v", err)
			return
		}
		stdio.out = f
		defer f.Close()
	case parser.TruncateAll:
		f, err := os.Create(cmd.Target)
		if err != nil {
			fmt.Fprintf(stdio.err, "msh: %+v", err)
			return
		}
		stdio.out = f
		stdio.err = f
		defer f.Close()
	case parser.Append:
		f, err := os.OpenFile(cmd.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(stdio.err, "msh: %+v", err)
			return
		}
		stdio.out = f
		defer f.Close()
	}

	evaluate(stdio, cmd.Command)
}
