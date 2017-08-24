package piper

import (
	"bytes"
	"compress/gzip"
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

// Returns a new io.Reader that will compress and encrypt the contents of the source reader.
func (s *Stream) EncryptReader(srcReader io.Reader) (io.Reader, error) {
	// Create a random IV for OFB
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		return nil, err
	}
	ivReader := bytes.NewReader(iv[:])

	// Hold encrypted data.
	bodyBuffer := &bytes.Buffer{}
	// Neew encryptor that writes to the buffer.
	stream := cipher.NewOFB(s.block, iv[:])
	encryptWriter := &cipher.StreamWriter{S: stream, W: bodyBuffer}
	// New compressor that writes to the encryptor.
	compressWriter := gzip.NewWriter(encryptWriter)
	// Add IV as prefix to data.
	ivAndEncryptedSrcReader := io.MultiReader(ivReader, bodyBuffer)

	// Build a redirect so we can Read from the pipeline above.
	redirect := NewRedirect(srcReader, compressWriter, ivAndEncryptedSrcReader)
	return redirect, nil
}

// Returns an io.Reader that decrypts and then decompresses the source reader.
func (s *Stream) DecryptReader(src io.Reader) (io.Reader, error) {
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(src, iv[:]); err != nil {
		return nil, err
	}

	// Wrap the body in a streaming OFB encryption block.
	stream := cipher.NewOFB(s.block, iv[:])
	// First we decrypt.
	decryptReader := &cipher.StreamReader{S: stream, R: src}
	// Then we decompress.
	decompressReader, err := gzip.NewReader(decryptReader)
	if err != nil {
		return nil, err
	}

	// HACK: This is not a good way to ensure the gzip reader is closed.
	closeOnReaderEOF := newCloseOnReaderEOF(decompressReader)

	return closeOnReaderEOF, nil
}
