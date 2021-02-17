package command

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	iio "github.com/jacobsimpson/msh/interpreter/io"
	"github.com/jacobsimpson/msh/parser"
)

func ExecuteProgram(stdio *iio.IOChannels, command *parser.Exec) <-chan int {
	cmd := exec.Command(command.Name, command.Arguments...)
	cmd.Stdin = stdio.In.Reader
	cmd.Stdout = stdio.Out.Writer
	cmd.Stderr = stdio.Err.Writer

	// Listen for Ctrl-C.
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		// If the channel returns nil, it was closed without a signal.
		if <-signals != nil {
			cmd.Process.Kill()
			// For some reason, this Println displays ^C twice, but if I take
			// it out, ^C isn't displayed at all. However, if I add something
			// extra to this string, the extra is only displayed once. So,
			// there is another ^C coming from somewhere else, but it is only
			// displayed if this Println is executed.
			fmt.Println("^C")
		}
	}()

	err := cmd.Start()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			fmt.Fprintf(stdio.Err.Writer, "msh: command not found: %s\n", command.Name)
		}
		return done(1)
	}

	c := make(chan int)
	go func() {
		status := 0
		err := cmd.Wait()
		if err != nil {
			status = 1
			if exiterr, ok := err.(*exec.ExitError); ok {
				// The program has exited with an exit code != 0

				// This works on both Unix and Windows. Although package
				// syscall is generally platform dependent, WaitStatus is
				// defined for both Unix and Windows and in both cases has
				// an ExitStatus() method with the same signature.
				if s, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					status = s.ExitStatus()
				}
			}
		}

		// Stop listening for Ctrl-C signals.
		signal.Stop(signals)
		// Close the channel so the go routine handling these knows it's time to stop.
		close(signals)

		stdio.In.Close()
		stdio.Out.Close()
		stdio.Err.Close()

		c <- status
		close(c)
	}()

	return c
}

func done(status int) <-chan int {
	c := make(chan int, 1)
	c <- status
	close(c)
	return c
}
