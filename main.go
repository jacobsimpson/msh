//go:generate pigeon -o parser/grammar.go parser/grammar.peg
//go:generate goimports -w parser/grammar.go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/jacobsimpson/msh/color"
	"github.com/jacobsimpson/msh/interpreter"
	"github.com/jacobsimpson/msh/interpreter/builtin"
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
		r.SetPrompt(calculatePrompt())
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

		c := interpreter.Evaluate(ast.(*parser.Program))
		<-c
	}
}

func calculatePrompt() string {
	d, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("%s %s ", color.Red("<error>"), color.Yellow("%"))
	}
	return fmt.Sprintf("%s %s ", color.Cyan(d), color.Yellow("%"))
}
