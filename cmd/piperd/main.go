package main

import (
	"flag"
	"fmt"
	"github.com/film42/piper"
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
	config, err := piper.LoadConfig(*configPathPtr)
	if err != nil {
		fmt.Println("Error loading config file:", err)
		os.Exit(1)
	}

	// Create local client to IPFS
	shell := ipfs.NewLocalShell()

	// Setup encryption key
	stream, err := piper.NewStream(config.OFBSymmetricKey)
	if err != nil {
		fmt.Println("Error building encryption stream:", err)
		os.Exit(1)
	}

	// Create server context
	serverContext := &ServerContext{
		config: config,
		shell:  shell,
		stream: stream,
	}

	// Let's go!
	serverContext.ListenAndServe()
}
