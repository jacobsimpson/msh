package interpreter

import (
	"fmt"
	"io"
	"os"

	"github.com/jacobsimpson/msh/interpreter/builtin"
	"github.com/jacobsimpson/msh/interpreter/command"
	iio "github.com/jacobsimpson/msh/interpreter/io"
	"github.com/jacobsimpson/msh/parser"
)

func Evaluate(program *parser.Program) <-chan int {
	stdio := &iio.IOChannels{
		In:  &iio.Reader{os.Stdin, false},
		Out: &iio.Writer{os.Stdout, false},
		Err: &iio.Writer{os.Stderr, false},
	}
	return evaluate(stdio, program.Command)
}

func evaluate(stdio *iio.IOChannels, cmd parser.Command) <-chan int {
	switch c := cmd.(type) {
	case *parser.Exec:
		return evaluateExec(stdio, c)
	case *parser.Redirection:
		return evaluateRedirection(stdio, c)
	case *parser.Pipe:
		return evaluatePipe(stdio, c)
	}
	return done(1)
}

func evaluateExec(stdio *iio.IOChannels, cmd *parser.Exec) <-chan int {
	if cmd.Name == "" {
		// Do nothing.
		return done(0)
	} else if c := builtin.Get(cmd.Name); c != nil {
		return c.Execute(stdio, cmd.Arguments)
	}
	return command.ExecuteProgram(stdio, cmd)
}

func evaluateRedirection(stdio *iio.IOChannels, cmd *parser.Redirection) <-chan int {
	switch cmd.Type {
	case parser.Truncate:
		// This file will be closed by the execution process, when that process
		// completes. Processes execute async, so there is no place in the
		// stack that is safe to close the file.
		f, err := os.Create(cmd.Target)
		if err != nil {
			fmt.Fprintf(stdio.Err.Writer, "msh: %+v", err)
			return done(1)
		}
		stdio.Out = &iio.Writer{f, true}
	case parser.TruncateAll:
		// This file will be closed by the execution process, when that process
		// completes. Processes execute async, so there is no place in the
		// stack that is safe to close the file.
		f, err := os.Create(cmd.Target)
		if err != nil {
			fmt.Fprintf(stdio.Err.Writer, "msh: %+v", err)
			return done(1)
		}
		stdio.Out = &iio.Writer{f, true}
		stdio.Err = &iio.Writer{f, true}
	case parser.Append:
		// This file will be closed by the execution process, when that process
		// completes. Processes execute async, so there is no place in the
		// stack that is safe to close the file.
		f, err := os.OpenFile(cmd.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(stdio.Err.Writer, "msh: %+v", err)
			return done(1)
		}
		stdio.Out = &iio.Writer{f, true}
	}

	return evaluate(stdio, cmd.Command)
}

func evaluatePipe(stdio *iio.IOChannels, cmd *parser.Pipe) <-chan int {
	r, w := io.Pipe()

	s := evaluate(&iio.IOChannels{In: stdio.In, Out: &iio.Writer{w, true}, Err: stdio.Err}, cmd.Src)
	d := evaluate(&iio.IOChannels{In: &iio.Reader{r, true}, Out: stdio.Out, Err: stdio.Err}, cmd.Dst)
	c := make(chan int)
	go func() {
		c <- max(<-s, <-d)
		close(c)
	}()
	return c
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

type noopReadCloser struct{ reader io.Reader }

func (r *noopReadCloser) Read(b []byte) (int, error) {
	return r.reader.Read(b)
}

func (r *noopReadCloser) Close() error { return nil }

type noopWriteCloser struct{ writer io.Writer }

func (w *noopWriteCloser) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *noopWriteCloser) Close() error { return nil }

func done(status int) <-chan int {
	c := make(chan int, 1)
	c <- status
	close(c)
	return c
}
