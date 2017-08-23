package main

import (
	"crypto/aes"
	"flag"
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	"os"
)

func main() {
	configPathPtr := flag.String("config", "", "Path to config file.")
	flag.Parse()

	if len(*configPathPtr) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Load and validate config
	config, err := LoadConfig(*configPathPtr)
	if err != nil {
		fmt.Println("Error loading config file:", err)
		os.Exit(1)
	}

	// Create local client to IPFS
	shell := ipfs.NewLocalShell()

	// Setup encryption key
	key := []byte(config.OFBSymmetricKey)
	block, _ := aes.NewCipher(key)

	// Create server context
	serverContext := &ServerContext{
		block:  block,
		config: config,
		shell:  shell,
	}

	// Let's go!
	serverContext.ListenAndServe()
}
