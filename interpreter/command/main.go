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
	// Not closing stdin here because wrapping os.Stdin in a noop close
	// implementation causes commands to halt after execution and wait for user
	// input.
	defer stdout.Close()
	defer stderr.Close()

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
	}

	cmd.Wait()

	// Stop listening for Ctrl-C signals.
	signal.Stop(signals)
	// Close the channel so the go routine handling these knows it's time to stop.
	close(signals)
	return done(0)
}

func done(status int) <-chan int {
	c := make(chan int)
	go func() {
		c <- status
		close(c)
	}()
	return c
}
