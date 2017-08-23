package main

import (
	"flag"
	"fmt"
	ipfs "github.com/ipfs/go-ipfs-api"
	"os"
)

func main() {
	configPathPtr := flag.String("config", "", "Path to config file.")
	setFilePathPtr := flag.String("set", "", "Path to file to SET.")
	startServerPtr := flag.Bool("start-server", false, "Start a gateway server accepting POST / GET requests.")
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
	stream, err := NewStream(config.OFBSymmetricKey)
	if err != nil {
		fmt.Println("Error building encryption stream:", err)
		os.Exit(1)
	}

	switch {
	case *startServerPtr:
		// Create server context
		serverContext := &ServerContext{
			config: config,
			shell:  shell,
			stream: stream,
		}

		// Let's go!
		serverContext.ListenAndServe()

	case len(*setFilePathPtr) > 0:
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

	default:
		flag.Usage()
		os.Exit(1)
	}
}
