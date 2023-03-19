package netclip

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string `yaml:"port"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// LoadConfig loads the configuration file from the given path
func LoadConfig(configFile string) (Config, error) {
	data, err := ioutil.ReadFile(configFile)
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
