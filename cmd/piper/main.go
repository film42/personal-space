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
	setFilePathPtr := flag.String("set", "", "Path to file to SET.")
	flag.Parse()

	if len(*configPathPtr) == 0 || len(*setFilePathPtr) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Load and validate config
	config, err := piper.LoadConfig(*configPathPtr)
	if config == nil && len(config.OFBSymmetricKey) > 0 {
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

	// Upload the file to ipfs
	cli := &CliContext{
		shell:  shell,
		stream: stream,
	}
	hash, err := cli.Set(*setFilePathPtr)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	} else {
		// Avoiding a newline here.
		fmt.Print(hash)
	}
}
