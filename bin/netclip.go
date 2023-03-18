package main

import (
	"flag"
	"fmt"
	"netclip"
	"os"
)

func main() {

	var port string
	version := flag.Bool("v", false, "Prints current app version.")
	flag.StringVar(&port, "port", "9999", "Port to use")

	flag.Parse()

	if *version {
		netclip.Banner()
		os.Exit(0)
	}
	fmt.Printf("Starting the web server on port %s\n", port)
	netclip.Run(port)

}
