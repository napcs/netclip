package netclip

import (
	"os"
	"path/filepath"
	"runtime"

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

// GetConfigPaths returns a list of paths to search for config files in priority order
func GetConfigPaths() []string {
	var paths []string
	
	// 1. Executable directory (current behavior)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		paths = append(paths, filepath.Join(exeDir, "netclip.yml"))
	}
	
	// 2. Current working directory (current fallback behavior)
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, "netclip.yml"))
	}
	
	// 3. User home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".netclip.yml"))
		paths = append(paths, filepath.Join(homeDir, "netclip.yml"))
	}
	
	// 4. XDG_CONFIG_HOME or ~/.config (Linux and macOS)
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		configDir := os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			if homeDir, err := os.UserHomeDir(); err == nil {
				configDir = filepath.Join(homeDir, ".config")
			}
		}
		if configDir != "" {
			paths = append(paths, filepath.Join(configDir, "netclip", "netclip.yml"))
			paths = append(paths, filepath.Join(configDir, "netclip.yml"))
		}
	}
	
	// 5. Linux system paths
	if runtime.GOOS == "linux" {
		paths = append(paths, "/etc/netclip/netclip.yml")
		paths = append(paths, "/usr/local/etc/netclip/netclip.yml")
		paths = append(paths, "/opt/netclip/netclip.yml")
	}
	
	// 6. macOS system paths
	if runtime.GOOS == "darwin" {
		paths = append(paths, "/etc/netclip/netclip.yml")
		paths = append(paths, "/usr/local/etc/netclip/netclip.yml")
		paths = append(paths, "/opt/netclip/netclip.yml")
		paths = append(paths, "/Library/Application Support/netclip/netclip.yml")
	}
	
	// 7. Windows system paths
	if runtime.GOOS == "windows" {
		if programData := os.Getenv("PROGRAMDATA"); programData != "" {
			paths = append(paths, filepath.Join(programData, "netclip", "netclip.yml"))
		}
		if programFiles := os.Getenv("PROGRAMFILES"); programFiles != "" {
			paths = append(paths, filepath.Join(programFiles, "netclip", "netclip.yml"))
		}
	}
	
	return paths
}

// LoadConfigFromPaths tries to load config from standard locations
func LoadConfigFromPaths() (Config, error) {
	var config Config
	var err error
	
	configPaths := GetConfigPaths()
	
	for _, path := range configPaths {
		config, err = LoadConfig(path)
		if err == nil {
			return config, nil
		}
	}
	
	return Config{}, err
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
