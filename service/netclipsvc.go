package main

import (
	"flag"
	"fmt"
	"log"
	"netclip"
	"os"

	"github.com/kardianos/service"
)

func banner() {
	fmt.Println("netclip v" + netclip.AppVersion)
}

var logger service.Logger

type program struct {
	Port string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run(p.Port)
	return nil
}

func (p *program) run(port string) {
	netclip.Run(port)
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	var mode string
	var port string
	version := flag.Bool("v", false, "Prints current app version.")
	flag.StringVar(&mode, "mode", "", "install/restart/start/stop/uninstall")
	flag.StringVar(&port, "port", "9999", "Port to use")
	flag.Parse()

	if *version {
		banner()
		os.Exit(0)
	}

	svcConfig := &service.Config{
		Name:        "netclip",
		DisplayName: "netclip service",
		Description: "Tiny server for a text clipboard",
	}

	prg := &program{Port: port}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	switch mode {

	case "start":
		if err = s.Start(); err != nil {
			log.Fatal(err)
		}

	case "stop":
		if err = s.Stop(); err != nil {
			log.Fatal(err)
		}

	case "restart":
		if err = s.Restart(); err != nil {
			log.Fatal(err)
		}

	case "install":
		if err = s.Install(); err != nil {
			log.Fatal(err)
		}

	case "uninstall":
		s.Stop()

		if err = s.Uninstall(); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Printf("Starting the web server on port %s\n", port)
		if err = s.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
