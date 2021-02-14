package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jacobsimpson/msh/command"
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

		switch strings.TrimSpace(line) {
		case "exit":
			command.Exit()
		case "pwd":
			command.PWD()
		default:
			fmt.Printf("You said: %q\n", line)
		}
	}
}
