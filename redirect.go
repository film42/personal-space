package main

import (
	"io"
)

type NoOpWriteCloser struct {
	writer io.Writer
}

func (c *NoOpWriteCloser) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}

func (c *NoOpWriteCloser) Close() error {
	return nil
}

type Redirect struct {
	src      io.Reader
	srcEOF   bool
	redirect io.WriteCloser
	dest     io.Reader
}

func NewRedirect(src io.Reader, redirect io.WriteCloser, dest io.Reader) *Redirect {
	var writeCloser io.WriteCloser
	if wcloser, ok := redirect.(io.WriteCloser); ok {
		writeCloser = wcloser
	} else {
		writeCloser = &NoOpWriteCloser{writer: redirect}
	}

	return &Redirect{
		src:      io.TeeReader(src, redirect),
		srcEOF:   false,
		dest:     dest,
		redirect: writeCloser,
	}
}

func (r *Redirect) Read(p []byte) (int, error) {
	if !r.srcEOF {
		_, err := r.src.Read(p)
		if err == io.EOF {
			r.srcEOF = true
			r.redirect.Close()
		}
	}
	return r.dest.Read(p)
}
