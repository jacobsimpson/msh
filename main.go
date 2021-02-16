//go:generate pigeon -o parser/grammar.go parser/grammar.peg
//go:generate goimports -w parser/grammar.go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/jacobsimpson/msh/builtin"
	"github.com/jacobsimpson/msh/command"
	"github.com/jacobsimpson/msh/parser"
	flag "github.com/spf13/pflag"
)

var showVersion = flag.BoolP("version", "v", false, "show the version")

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version %s\n", builtin.Version)
		os.Exit(0)
	}

	r, err := readline.New("msh> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize newline library: %+v\n", err)
		os.Exit(1)
	}

	for {
		line, err := r.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		ast, err := parser.Parse("shell", []byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Don't understand: %+v\n", err)
			continue
		}

		evaluate(ast.(*parser.Program))
	}
}

func evaluate(program *parser.Program) {
	stdin, stdout, stderr := os.Stdin, os.Stdout, os.Stderr

	if program.Command.Redirection != nil {
		switch program.Command.Redirection.Type {
		case parser.Truncate:
			f, err := os.Create(program.Command.Redirection.Target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "msh: %+v", err)
				return
			}
			stdout = f
			defer f.Close()
		case parser.TruncateAll:
			f, err := os.Create(program.Command.Redirection.Target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "msh: %+v", err)
				return
			}
			stdout = f
			stderr = f
			defer f.Close()
		case parser.Append:
			f, err := os.OpenFile(program.Command.Redirection.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "msh: %+v", err)
				return
			}
			stdout = f
			defer f.Close()
		}
	}

	if program.Command.Name == "" {
		// Do nothing.
	} else if cmd := builtin.Get(program.Command.Name); cmd != nil {
		cmd.Execute(stdin, stdout, stderr, program.Command.Arguments)
	} else {
		command.ExecuteProgram(stdin, stdout, stderr, program.Command)
	}
}
