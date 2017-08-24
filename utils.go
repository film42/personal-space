package piper

import (
	"io"
)

type CloseOnReaderEOF struct {
	reader    io.Reader
	closer    io.Closer
	wasClosed bool
}

func NewCloseOnReaderEOF(src io.Reader) *CloseOnReaderEOF {
	var closer io.Closer
	if c, ok := src.(io.Closer); ok {
		closer = c
	}

	return &CloseOnReaderEOF{
		closer:    closer,
		wasClosed: false,
		reader:    src,
	}
}

func (core *CloseOnReaderEOF) Read(p []byte) (int, error) {
	n, err := core.reader.Read(p)
	if err == io.EOF && core.closer != nil && !core.wasClosed {
		core.wasClosed = true
		core.closer.Close()
	}
	return n, err
}
