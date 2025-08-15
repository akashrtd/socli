package config

// DefaultConfig returns the default configuration for the application.
func DefaultConfig() *Config {
	return &Config{
		Network: struct {
			ListenPort     int      `yaml:"listen_port"`
			BootstrapPeers []string `yaml:"bootstrap_peers"`
			EnableMDNS     bool     `yaml:"enable_mdns"`
			EnableDHT      bool     `yaml:"enable_dht"`
		}{
			ListenPort:     0, // 0 means a random port
			BootstrapPeers: []string{},
			EnableMDNS:     true,
			EnableDHT:      true,
		},
		UI: struct {
			Theme         string `yaml:"theme"`
			RefreshRate   int    `yaml:"refresh_rate_ms"`
			MaxPostLength int    `yaml:"max_post_length"`
		}{
			Theme:         "default",
			RefreshRate:   100,
			MaxPostLength: 280,
		},
		Privacy: struct {
			EncryptMessages bool   `yaml:"encrypt_messages"`
			KeyPath         string `yaml:"key_path"`
			AutoClear       bool   `yaml:"auto_clear_on_exit"`
		}{
			EncryptMessages: true,
			KeyPath:         "socli.key",
			AutoClear:       true,
		},
	}
}
