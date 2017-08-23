package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

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

func (s *Stream) newStream(src io.Reader, iv []byte) io.Reader {
	// Wrap the body in a streaming OFB encryption block
	stream := cipher.NewOFB(s.block, iv[:])
	return &cipher.StreamReader{S: stream, R: src}
}

func (s *Stream) EncryptReader(srcReader io.Reader) (io.Reader, error) {
	// Create a random IV for OFB
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		return nil, err
	}

	ivReader := bytes.NewReader(iv[:])
	encryptedSrc := s.newStream(srcReader, iv[:])
	ivAndEncryptedSrcReader := io.MultiReader(ivReader, encryptedSrc)
	return ivAndEncryptedSrcReader, nil
}

func (s *Stream) DecryptReader(src io.Reader) (io.Reader, error) {
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(src, iv[:]); err != nil {
		return nil, err
	}
	return s.newStream(src, iv[:]), nil
}
