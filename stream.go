package main

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type NoOpWriteCloser struct {
	*io.Writer
}
func (c *NoOpWriteCloser) Write(p []byte) (int, error) {
	return c.Writer.Write(p)
}
func (c *NoOpWriteCloser) Close() error {
	return nil
}

type Redirect struct {
	src io.Reader
	srcEOF bool
	redirect io.WriteCloser
	dest io.Reader
}
func NewRedirect(src io.Reader, redirect io.WriteCloser, dest io.Reader) *Redirect {
	var writeCloser io.WriteCloser
	if wcloser, ok := redirect.(io.WriteCloser); ok {
		writeCloser =  wcloser
	} else {
		writeCloser = &NoOpWriteCloser{redirect}
	}

	return &Redirect{
		src: io.TeeReader(src, redirect),
		srcEOF: false,
		dest: dest,
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


type Stream struct {
	block cipher.Block
}

func NewStream(key string) (*Stream, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	return &Stream{block: block}, nil
}

func Read(b []byte) (int, error) {
	stream.Input.ReadUpTo(len(b))
	return stream.Output.Read(b)
}

func (s *Stream) EncryptReader(srcReader io.Reader) (io.Reader, error) {
	// Create a random IV for OFB
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		return nil, err
	}
	ivReader := bytes.NewReader(iv[:])

	bodyBuffer := &bytes.Buffer{}
	encryptWriter := &cipher.StreamWriter{S: stream, W: bodyBuffer}
	compressWriter := gzip.NewWriter(encryptWriter)

	ivAndEncryptedSrcReader := io.MultiReader(ivReader, bodyBuffer)
	return ivAndEncryptedSrcReader, nil
}

func (s *Stream) DecryptReader(src io.Reader) (io.Reader, error) {
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(src, iv[:]); err != nil {
		return nil, err
	}

	// Wrap the body in a streaming OFB encryption block
	stream := cipher.NewOFB(s.block, iv[:])
	reader := &cipher.StreamReader{S: stream, R: src}

	return s.newStream(src, iv[:]), nil
}
