//go:generate pigeon -o parser/main.go parser/grammar.peg
//go:generate goimports -w parser/main.go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/jacobsimpson/msh/builtin"
	"github.com/jacobsimpson/msh/command"
	"github.com/jacobsimpson/msh/parser"
)

func main() {
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

		program := ast.(*parser.Program)

		if cmd := builtin.Get(program.Command.Name); cmd != nil {
			cmd.Execute(program.Command.Arguments)
		} else {
			switch program.Command.Name {
			case "pwd":
				builtin.PWD()
			case "export":
				builtin.Export()
			case "":
				break
			default:
				command.ExecuteProgram(program.Command)
			}
		}
	}
}
