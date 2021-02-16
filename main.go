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

	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open file: %+v\n", err)
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

		if program.Command.Name == "" {
			// Do nothing.
		} else if cmd := builtin.Get(program.Command.Name); cmd != nil {
			cmd.Execute(os.Stdin, os.Stdout, os.Stderr, program.Command.Arguments)
		} else {
			command.ExecuteProgram(program.Command)
		}
	}
}
