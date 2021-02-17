package builtin

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jacobsimpson/msh/color"
)

var directoryHistory *history

func init() {
	dh, err := newHistory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to find the current directory: %+v\n", err)
	}
	directoryHistory = dh
}

type cd struct{}

func (c *cd) Execute(stdin io.ReadCloser, stdout, stderr io.WriteCloser, args []string) int {
	defer stdout.Close()
	defer stderr.Close()

	dst := ""
	updateHistory := true
	if len(args) == 0 || args[0] == "~" {
		dst = os.Getenv("HOME")
	} else if args[0] == "#" {
		for i, d := range directoryHistory.directories {
			marker := " "
			if i == directoryHistory.current {
				marker = "*"
			}
			fmt.Fprintf(stdout, "%s %s\n", color.Blue(marker), d)
		}
		return 1
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
		fmt.Fprintf(stderr, "no such file or directory: %s", dst)
		return 1
	}
	if updateHistory {
		directoryHistory.add(dst)
	}
	return 0
}

func (*cd) Name() string { return "cd" }

func (*cd) ShortHelp() string { return "cd [dir]" }

func (*cd) LongHelp() string {
	return "Change the current directory to DIR. The variable $HOME is the default DIR. There are a few shortcuts available. `cd ~` changes the current directory to the value of $HOME, `cd -` moves down the previous directory stack. `cd --` moves down the previous directory stack by two. `cd +` moves up the previous directory stack. `cd #` shows the previous directory stack."
}

func all(s string, r rune) bool {
	for _, c := range s {
		if c != r {
			return false
		}
	}
	return true
}
