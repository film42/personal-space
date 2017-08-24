package piper

import (
	"io"
)

type closeOnReaderEOF struct {
	reader    io.Reader
	closer    io.Closer
	wasClosed bool
}

func newCloseOnReaderEOF(src io.Reader) *closeOnReaderEOF {
	var closer io.Closer
	if c, ok := src.(io.Closer); ok {
		closer = c
	}

	return &closeOnReaderEOF{
		closer:    closer,
		wasClosed: false,
		reader:    src,
	}
}

func (core *closeOnReaderEOF) Read(p []byte) (int, error) {
	n, err := core.reader.Read(p)
	if err == io.EOF && core.closer != nil && !core.wasClosed {
		core.wasClosed = true
		core.closer.Close()
	}
	return n, err
}
