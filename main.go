package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
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

		fmt.Printf("You said: %q\n", line)
	}
}
