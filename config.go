package netclip

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string      `yaml:"port"`
	CertFile string      `yaml:"cert_file"`
	KeyFile  string      `yaml:"key_file"`
	Tailscale    TailscaleConfig `yaml:"tailscale"`
}

type TailscaleConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Hostname string `yaml:"hostname"`
	UseTLS   bool   `yaml:"use_tls"`
}

// LoadConfig loads the configuration file from the given path
func LoadConfig(configFile string) (Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}


	return config, nil
}

// ApplyFlags applies command line flag values to the config, with flags taking precedence
func ApplyFlags(config Config, port, certFile, keyFile, tailscaleHostname string, tailscaleEnabled, tailscaleTLS bool) Config {
	// Port handling
	if port != "" {
		config.Port = port
	} else if config.Port == "" {
		config.Port = "9999"
	}

	// SSL certificate handling
	if certFile != "" {
		config.CertFile = certFile
	}
	if keyFile != "" {
		config.KeyFile = keyFile
	}

	// Tailscale flag overrides
	if tailscaleEnabled {
		config.Tailscale.Enabled = true
	}
	if tailscaleHostname != "" {
		config.Tailscale.Hostname = tailscaleHostname
	} else if config.Tailscale.Enabled && config.Tailscale.Hostname == "" {
		config.Tailscale.Hostname = "netclip"
	}
	if tailscaleTLS {
		config.Tailscale.UseTLS = true
	}

	return config
}
