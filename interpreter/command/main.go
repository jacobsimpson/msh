package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/jacobsimpson/msh/parser"
)

func ExecuteProgram(stdin io.ReadCloser, stdout, stderr io.WriteCloser, command *parser.Exec) <-chan int {
	cmd := exec.Command(command.Name, command.Arguments...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

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
			fmt.Fprintf(stderr, "msh: command not found: %s\n", command.Name)
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

		// Not closing stdin here because wrapping os.Stdin in a noop close
		// implementation causes commands to halt after execution and wait for user
		// input.
		stdout.Close()
		stderr.Close()

		c <- status
		close(c)
	}()

	return c
}

func done(status int) <-chan int {
	c := make(chan int)
	go func() {
		c <- status
		close(c)
	}()
	return c
}
