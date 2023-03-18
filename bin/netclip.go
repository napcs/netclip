package main

import (
	"flag"
	"fmt"
	"netclip"
	"os"
)

func banner() {
	fmt.Println("netclip v" + netclip.AppVersion)
}

func main() {

	version := flag.Bool("v", false, "Prints current app version.")

	flag.Parse()

	if *version {
		banner()
		os.Exit(0)
	}
	port := "4000"

	fmt.Printf("Starting the web server on port %s\n", port)
	netclip.Run(port)

}
