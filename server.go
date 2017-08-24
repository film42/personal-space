package main

import (
	"bufio"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"

	ipfs "github.com/ipfs/go-ipfs-api"
)

type ServerContext struct {
	shell  *ipfs.Shell
	stream *Stream
	config *Config
}

func (sc *ServerContext) ListenAndServe() {
	server := echo.New()
	server.POST("/set", sc.Set)
	server.GET("/get/:hash", sc.Get)

	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}, bytes_in=${bytes_in}, bytes_out=${bytes_out}\n",
	}))

	server.Logger.Fatal(server.Start(sc.config.Bind))
}

func (sc *ServerContext) isAuthenticated(context echo.Context) bool {
	apiKey := context.Request().Header.Get("X-Api-Key")
	if apiKey != sc.config.ApiKey {
		return false
	}
	return true
}

func (sc *ServerContext) Set(context echo.Context) error {
	if !sc.isAuthenticated(context) {
		return context.NoContent(http.StatusUnauthorized)
	}

	request := context.Request()
	if request.ContentLength < 1 {
		return context.String(http.StatusBadRequest, "Request payload is too small")
	}

	// Pull out the request body
	body := request.Body
	defer body.Close()

	encryptedBodyReader, err := sc.stream.EncryptReader(body)
	if err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}

	// Write the file to ipfs and get the hash back
	hash, err := sc.shell.Add(encryptedBodyReader)
	if err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}

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

	decryptedBodyReader, err := sc.stream.DecryptReader(readCloser)
	if err != nil {
		return context.String(http.StatusInternalServerError, err.Error())
	}

	// Add this heder should we detect an image type. Let the browser decide.
	context.Response().Header().Set("Content-Disposition", "inline")

	// Detect the mime-type. This bufio should only use the buffer to peek and the rest and
	// after its been read it should read directly into the callers []byte.
	decryptedBuffer := bufio.NewReaderSize(decryptedBodyReader, 128)
	peekedBytes, err := decryptedBuffer.Peek(128)
	if err != nil && peekedBytes == nil {
		// Just fall back for now since we have the data.
		peekedBytes = []byte{}
	}
	mimeType := http.DetectContentType(peekedBytes)

	// Let's send if over the wire.
	return context.Stream(http.StatusOK, mimeType, decryptedBuffer)
}
