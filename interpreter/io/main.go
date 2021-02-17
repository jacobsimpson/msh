package io

import "io"

type Reader struct {
	Reader      io.ReadCloser
	ShouldClose bool
}

func (r *Reader) Close() {
	if r.ShouldClose {
		r.Reader.Close()
	}
}

type Writer struct {
	Writer      io.WriteCloser
	ShouldClose bool
}

func (w *Writer) Close() {
	if w.ShouldClose {
		w.Writer.Close()
	}
}

// IOChannels contains the IO stream triple that describes how a process can
// read and write.
//
// For some reason, wrapping something directly around the Readers and Writers
// available from os.Stdin, os.Stdout, and os.Stderr leads to unexpected
// behavior. Processes that are started don't recognize the wrapped Readers and
// Writers as belonging to the termina. It compiles, and runs, but in the
// shell, after a command exits, if Stdin is wrapped, it waits for the user to
// type something before ending the subprocess. If Stdout is wrapped, ls prints
// in a single column, and less and more won't page.
//
// It is reproducable with both the custom noopReadCloser and ioutil.NopCloser.
//
// The intended use of the wrappers is to allow exiting processes to close
// their own readers and writers, signalling to any connected processes that
// there is no longer someone at the other end. If the readers and writers are
// the actual streams connected to the terminal, I don't want to close them, so
// the close request should be ignored.
type IOChannels struct {
	In  *Reader
	Out *Writer
	Err *Writer
}
