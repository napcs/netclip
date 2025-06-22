package main

import (
	"flag"
	"fmt"
	"log"
	"netclip"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

func banner() {
	fmt.Println("netclip v" + netclip.AppVersion)
}

var logger service.Logger

type program struct {
	Port          string
	CertFile      string
	KeyFile       string
	TailscaleEnabled  bool
	TailscaleHostname string
	TailscaleAuthKey  string
	TailscaleUseTLS   bool
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	if p.TailscaleEnabled {
		server := &netclip.TSNetServer{
			Hostname: p.TailscaleHostname,
			AuthKey:  p.TailscaleAuthKey,
			UseTLS:   p.TailscaleUseTLS,
		}
		netclip.Run(server)
	} else {
		server := &netclip.HTTPServer{
			Port:     p.Port,
			CertFile: p.CertFile,
			KeyFile:  p.KeyFile,
		}
		netclip.Run(server)
	}
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	var serviceMode string

	version := flag.Bool("v", false, "Prints current app version.")
	flag.StringVar(&serviceMode, "service", "", "install/restart/start/stop/uninstall")

	// Define flags with pointers to detect if they were set
	portFlag := flag.String("port", "", "Port to use (default: 9999)")
	certFlag := flag.String("cert", "", "Path to SSL certificate file")
	keyFlag := flag.String("key", "", "Path to SSL private key file")
	tailscaleFlag := flag.Bool("tailscale", false, "Enable Tailscale networking")
	tailscaleHostnameFlag := flag.String("tailscale-hostname", "", "Tailscale hostname (default: netclip)")
	tailscaleTLSFlag := flag.Bool("tailscale-tls", false, "Use HTTPS with Tailscale certificates")

	flag.Parse()

	if *version {
		banner()
		os.Exit(0)
	}

	// get path to current executable to load the config file
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)

	configPath := filepath.Join(exeDir, "netclip.yml")
	config, err := netclip.LoadConfig(configPath)

	if err != nil {
		// Try loading config.yaml from the current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		configPath = filepath.Join(cwd, "netclip.yml")
		config, err = netclip.LoadConfig(configPath)
		if err != nil {
			log.Printf("Failed to load config: %v", err)
			log.Printf("Using default options.")
		}
	}

	// Apply flag overrides - flags take precedence over config file
	config = netclip.ApplyFlags(config, *portFlag, *certFlag, *keyFlag, *tailscaleHostnameFlag, *tailscaleFlag, *tailscaleTLSFlag)

	// a mode was passed. Someone wants to do service things.

	svcConfig := &service.Config{
		Name:        "netclip",
		DisplayName: "netclip Clipboard Server",
		Description: "Tiny server for a text clipboard for your network",
	}

	prg := &program{
		Port:          config.Port,
		CertFile:      config.CertFile,
		KeyFile:       config.KeyFile,
		TailscaleEnabled:  config.Tailscale.Enabled,
		TailscaleHostname: config.Tailscale.Hostname,
		TailscaleAuthKey:  os.Getenv("TS_AUTHKEY"),
		TailscaleUseTLS:   config.Tailscale.UseTLS,
	}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	// what kind of service operation do they want to perform?
	switch serviceMode {

	case "start":
		if err = s.Start(); err != nil {
			log.Fatal(err)
		}

	case "run":
		if err = s.Run(); err != nil {
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
		// just run as a standalone server
		if config.Tailscale.Enabled {
			fmt.Printf("Starting netclip on Tailscale as %s.ts.net\n", config.Tailscale.Hostname)
		} else {
			fmt.Printf("Starting netclip on port %s\n", config.Port)
		}
		if err = s.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
