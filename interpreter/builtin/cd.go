package builtin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jacobsimpson/msh/color"
	iio "github.com/jacobsimpson/msh/interpreter/io"
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

func (c *cd) Execute(stdio *iio.IOChannels, args []string) <-chan int {
	defer stdio.In.Close()
	defer stdio.Out.Close()
	defer stdio.Err.Close()

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
			fmt.Fprintf(stdio.Out.Writer, "%s %s\n", color.Blue(marker), d)
		}
		return done(1)
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
		fmt.Fprintf(stdio.Err.Writer, "no such file or directory: %s", dst)
		return done(1)
	}
	if updateHistory {
		directoryHistory.add(dst)
	}
	return done(0)
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
