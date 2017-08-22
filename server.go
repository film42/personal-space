package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/labstack/echo"
	"io"
	"net/http"

	ipfs "github.com/ipfs/go-ipfs-api"
)

type ServerContext struct {
	shell *ipfs.Shell
	block cipher.Block
}

func (sc *ServerContext) Upload(context echo.Context) error {
	request := context.Request()
	if request.ContentLength < 1 {
		return context.String(http.StatusInternalServerError, "Request payload is too small")
	}

	// Pull out the request body
	body := request.Body
	defer body.Close()

	// Create a random IV for OFB
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}
	ivReader := bytes.NewBuffer(iv[:])

	// Wrap the body in a streaming OFB encryption block
	stream := cipher.NewOFB(sc.block, iv[:])
	encryptedBody := &cipher.StreamReader{S: stream, R: body}
	ivAndEncryptedBodyReader := io.MultiReader(ivReader, encryptedBody)

	hash, err := sc.shell.Add(ivAndEncryptedBodyReader)
	if err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}

	// Write the file to ipfs and get the hash back
	return context.String(http.StatusOK, hash)
}

func (sc *ServerContext) Get(context echo.Context) error {
	hash := context.Param("hash")

	// Grab the io.ReadCloser from IPFS.
	readCloser, err := sc.shell.Cat(hash)
	if err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}

	// TODO: Let's double check that context.Stream does not leave this function until done.
	defer readCloser.Close()

	// Read the IV from the IPFS reader.
	var iv [aes.BlockSize]byte
	if _, err := io.ReadFull(readCloser, iv[:]); err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}

	// Wrap the body in a streaming OFB encryption block
	stream := cipher.NewOFB(sc.block, iv[:])
	decryptedBody := &cipher.StreamReader{S: stream, R: readCloser}

	// Add this heder should we detect an image type. Let the browser decide.
	context.Response().Header().Set("Content-Disposition", "inline")

	// Detect the mime-type.
	decryptedBuffer := bufio.NewReader(decryptedBody)
	peekedBytes, err := decryptedBuffer.Peek(128)
	if err != nil {
		// Just fall back for now since we have the data.
		peekedBytes = []byte{}
	}
	mimeType := http.DetectContentType(peekedBytes)

	// Let's send if over the wire.
	return context.Stream(http.StatusOK, mimeType, decryptedBuffer)
}

func startServer(bind string, serverContext *ServerContext) {
	server := echo.New()
	server.POST("/upload", serverContext.Upload)
	server.GET("/s/:hash", serverContext.Get)
	server.Logger.Fatal(server.Start(bind))
}
