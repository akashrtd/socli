package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the application's configuration.
type Config struct {
	Network struct {
		ListenPort     int      `yaml:"listen_port"`
		BootstrapPeers []string `yaml:"bootstrap_peers"`
		EnableMDNS     bool     `yaml:"enable_mdns"`
		EnableDHT      bool     `yaml:"enable_dht"`
	} `yaml:"network"`

	UI struct {
		Theme         string `yaml:"theme"`
		RefreshRate   int    `yaml:"refresh_rate_ms"`
		MaxPostLength int    `yaml:"max_post_length"`
	} `yaml:"ui"`

	Privacy struct {
		EncryptMessages bool   `yaml:"encrypt_messages"`
		KeyPath         string `yaml:"key_path"`
		AutoClear       bool   `yaml:"auto_clear_on_exit"`
	} `yaml:"privacy"`
}

// LoadConfig loads the configuration from the given path.
// If the file doesn't exist, it creates one with default values.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// File doesn't exist, so create it with default values
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		return cfg, nil
	} else if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
