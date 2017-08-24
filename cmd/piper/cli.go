package main

import (
	"github.com/film42/piper"
	ipfs "github.com/ipfs/go-ipfs-api"
	"os"
)

type CliContext struct {
	shell  *ipfs.Shell
	stream *piper.Stream
}

func (c *CliContext) Set(path string) (string, error) {
	// Open file
	file, err := os.Open(path)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	// Encrypt stream
	encryptedFile, err := c.stream.EncryptReader(file)
	if err != nil {
		return "", err
	}

	hash, err := c.shell.Add(encryptedFile)
	return hash, err
}
