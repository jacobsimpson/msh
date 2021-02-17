package interpreter

import (
	"fmt"
	"io"
	"os"

	"github.com/jacobsimpson/msh/interpreter/builtin"
	"github.com/jacobsimpson/msh/interpreter/command"
	"github.com/jacobsimpson/msh/parser"
)

type iochannels struct {
	in  io.ReadCloser
	out io.WriteCloser
	err io.WriteCloser
}

func Evaluate(program *parser.Program) <-chan int {
	stdio := &iochannels{
		// For some reason, using a noopReader around Stdin doesn't work. It
		// compiles, and runs, but in the shell, after a command exits, it
		// waits for the user to type something before ending. It is
		// reproducable with both the custom noopReadCloser and
		// ioutil.NopCloser. The intended use of Closers is to allow exiting
		// processes to close their own readers and writers, signalling to any
		// connected processes that there is no longer someone at the other
		// end. To start with, it should be sufficient to just close the
		// Writers, so I won't close in, so I don't need to wrap it for now.
		in:  os.Stdin,
		out: &noopWriteCloser{os.Stdout},
		err: &noopWriteCloser{os.Stderr},
	}
	return evaluate(stdio, program.Command)
}

func evaluate(stdio *iochannels, cmd parser.Command) <-chan int {
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

func evaluateExec(stdio *iochannels, cmd *parser.Exec) <-chan int {
	if cmd.Name == "" {
		// Do nothing.
		return done(0)
	} else if c := builtin.Get(cmd.Name); c != nil {
		return c.Execute(stdio.in, stdio.out, stdio.err, cmd.Arguments)
	}
	return command.ExecuteProgram(stdio.in, stdio.out, stdio.err, cmd)
}

func evaluateRedirection(stdio *iochannels, cmd *parser.Redirection) <-chan int {
	switch cmd.Type {
	case parser.Truncate:
		// This file will be closed by the execution process, when that process
		// completes. Processes execute async, so there is no place in the
		// stack that is safe to close the file.
		f, err := os.Create(cmd.Target)
		if err != nil {
			fmt.Fprintf(stdio.err, "msh: %+v", err)
			return done(1)
		}
		stdio.out = f
	case parser.TruncateAll:
		// This file will be closed by the execution process, when that process
		// completes. Processes execute async, so there is no place in the
		// stack that is safe to close the file.
		f, err := os.Create(cmd.Target)
		if err != nil {
			fmt.Fprintf(stdio.err, "msh: %+v", err)
			return done(1)
		}
		stdio.out = f
		stdio.err = f
	case parser.Append:
		// This file will be closed by the execution process, when that process
		// completes. Processes execute async, so there is no place in the
		// stack that is safe to close the file.
		f, err := os.OpenFile(cmd.Target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(stdio.err, "msh: %+v", err)
			return done(1)
		}
		stdio.out = f
	}

	return evaluate(stdio, cmd.Command)
}

func evaluatePipe(stdio *iochannels, cmd *parser.Pipe) <-chan int {
	r, w := io.Pipe()

	s := evaluate(&iochannels{in: stdio.in, out: w, err: stdio.err}, cmd.Src)
	d := evaluate(&iochannels{in: r, out: stdio.out, err: stdio.err}, cmd.Dst)
	c := make(chan int)
	go func() {
		c <- max(<-s, <-d)
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
	c := make(chan int)
	go func() {
		c <- status
		close(c)
	}()
	return c
}
