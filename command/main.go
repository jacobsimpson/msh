package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jacobsimpson/msh/color"
	"github.com/jacobsimpson/msh/parser"
)

func PWD() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to get the current working directory: %+v\n", err)
		return
	}
	fmt.Printf("%s\n", wd)
}

func Exit() {
	os.Exit(0)
}

var directoryHistory *history

func init() {
	dh, err := newHistory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to find the current directory: %+v\n", err)
	}
	directoryHistory = dh
}

func CD(args []string) {
	dst := ""
	updateHistory := true
	if len(args) == 0 {
		dst = os.Getenv("HOME")
	} else if args[0] == "#" {
		for i, d := range directoryHistory.directories {
			marker := " "
			if i == directoryHistory.current {
				marker = "*"
			}
			fmt.Printf("%s %s\n", color.Blue(marker), d)
		}
		return
	} else if all(args[0], rune('-')) {
		for i := 0; i < len(args[0]); i++ {
			dst = directoryHistory.back()
		}
		updateHistory = false
	} else if all(args[0], rune('+')) {
		for i := 0; i < len(args[0]); i++ {
			dst = directoryHistory.forward()
		}
		updateHistory = false
	} else if args[0] == "." {
		dst = directoryHistory.get()
	} else if args[0] == ".." {
		dst = filepath.Dir(directoryHistory.get())
	} else {
		dst = args[0]
	}
	if err := os.Chdir(dst); err != nil {
		fmt.Fprintf(os.Stderr, "no such file or directory: %s", dst)
		return
	}
	if updateHistory {
		directoryHistory.add(dst)
	}
}

func all(s string, r rune) bool {
	for _, c := range s {
		if c != r {
			return false
		}
	}
	return true
}

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

func Export() {
	for _, e := range os.Environ() {
		fmt.Printf("%s\n", e)
	}
}
