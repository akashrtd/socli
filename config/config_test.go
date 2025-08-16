package config

import (
	"os"
	"testing"
)

// TestLoadConfigDefault tests that a default config is created if the file doesn't exist.
func TestLoadConfigDefault(t *testing.T) {
	// Create a temporary file name that shouldn't exist
	tmpFile := "test_config_default.yaml"
	// Ensure the file is deleted after the test
	defer os.Remove(tmpFile)

	// Make sure the file doesn't exist before the test
	_, err := os.Stat(tmpFile)
	if !os.IsNotExist(err) {
		// If the file exists, delete it
		os.Remove(tmpFile)
	}

	// Load the config (this should create a default file)
	cfg, err := LoadConfig(tmpFile)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil, want a config")
	}

	// Check that some default values are set
	if cfg.Network.ListenPort != 0 {
		t.Errorf("cfg.Network.ListenPort = %d, want 0 (default)", cfg.Network.ListenPort)
	}
	if cfg.UI.Theme != "default" {
		t.Errorf("cfg.UI.Theme = %s, want 'default'", cfg.UI.Theme)
	}
	if cfg.Privacy.KeyPath != "socli.key" {
		t.Errorf("cfg.Privacy.KeyPath = %s, want 'socli.key'", cfg.Privacy.KeyPath)
	}

	// Verify that the file was created
	_, err = os.Stat(tmpFile)
	if os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
}

// TestLoadConfigExisting tests that an existing config file is loaded correctly.
func TestLoadConfigExisting(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_config_existing_*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test
	defer tmpFile.Close()

	// Write a known config to the file
	configContent := `network:
  listen_port: 12345
  bootstrap_peers: []
  enable_mdns: false
  enable_dht: true
ui:
  theme: "dark"
  refresh_rate_ms: 200
  max_post_length: 140
privacy:
  encrypt_messages: true
  key_path: "mykey.key"
  auto_clear_on_exit: false
`
	_, err = tmpFile.Write([]byte(configContent))
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file to ensure it's flushed to disk
	tmpFile.Close()

	// Load the config
	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig() error = %v, want nil", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil, want a config")
	}

	// Check that the values from the file are loaded
	if cfg.Network.ListenPort != 12345 {
		t.Errorf("cfg.Network.ListenPort = %d, want 12345", cfg.Network.ListenPort)
	}
	if cfg.Network.EnableMDNS != false {
		t.Errorf("cfg.Network.EnableMDNS = %v, want false", cfg.Network.EnableMDNS)
	}
	if cfg.Network.EnableDHT != true {
		t.Errorf("cfg.Network.EnableDHT = %v, want true", cfg.Network.EnableDHT)
	}
	if cfg.UI.Theme != "dark" {
		t.Errorf("cfg.UI.Theme = %s, want 'dark'", cfg.UI.Theme)
	}
	if cfg.UI.RefreshRate != 200 {
		t.Errorf("cfg.UI.RefreshRate = %d, want 200", cfg.UI.RefreshRate)
	}
	if cfg.UI.MaxPostLength != 140 {
		t.Errorf("cfg.UI.MaxPostLength = %d, want 140", cfg.UI.MaxPostLength)
	}
	if cfg.Privacy.EncryptMessages != true {
		t.Errorf("cfg.Privacy.EncryptMessages = %v, want true", cfg.Privacy.EncryptMessages)
	}
	if cfg.Privacy.KeyPath != "mykey.key" {
		t.Errorf("cfg.Privacy.KeyPath = %s, want 'mykey.key'", cfg.Privacy.KeyPath)
	}
	if cfg.Privacy.AutoClear != false {
		t.Errorf("cfg.Privacy.AutoClear = %v, want false", cfg.Privacy.AutoClear)
	}
}