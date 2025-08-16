package crypto

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateKeyPair tests the GenerateKeyPair function.
func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair() error = %v, want nil", err)
	}

	if kp == nil {
		t.Fatal("GenerateKeyPair() returned nil, want a key pair")
	}

	if kp.PublicKey == nil {
		t.Error("GenerateKeyPair() returned a key pair with nil PublicKey")
	}

	if kp.PrivateKey == nil {
		t.Error("GenerateKeyPair() returned a key pair with nil PrivateKey")
	}

	// Test that the keys are 32 bytes
	if len(kp.PublicKey) != 32 {
		t.Errorf("GenerateKeyPair() returned PublicKey with length %d, want 32", len(kp.PublicKey))
	}

	if len(kp.PrivateKey) != 32 {
		t.Errorf("GenerateKeyPair() returned PrivateKey with length %d, want 32", len(kp.PrivateKey))
	}

	// Test that the keys are not zero
	allZero := true
	for _, b := range kp.PublicKey {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("GenerateKeyPair() returned a PublicKey that is all zeros")
	}

	allZero = true
	for _, b := range kp.PrivateKey {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("GenerateKeyPair() returned a PrivateKey that is all zeros")
	}
}

// TestSaveKeyPairAndLoadKeyPair tests the SaveKeyPair and LoadKeyPair functions.
func TestSaveKeyPairAndLoadKeyPair(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()

	// Generate a key pair
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Define a path for the key file
	keyPath := filepath.Join(tmpDir, "test_key")

	// Save the key pair
	err = SaveKeyPair(kp, keyPath)
	if err != nil {
		t.Fatalf("SaveKeyPair() error = %v, want nil", err)
	}

	// Check that the file was created
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		t.Fatal("SaveKeyPair() did not create the key file")
	}

	// Load the key pair
	loadedKp, err := LoadKeyPair(keyPath)
	if err != nil {
		t.Fatalf("LoadKeyPair() error = %v, want nil", err)
	}

	if loadedKp == nil {
		t.Fatal("LoadKeyPair() returned nil, want a key pair")
	}

	// Compare the keys
	if *kp.PrivateKey != *loadedKp.PrivateKey {
		t.Error("LoadKeyPair() returned a key pair with a different PrivateKey")
	}

	if *kp.PublicKey != *loadedKp.PublicKey {
		t.Error("LoadKeyPair() returned a key pair with a different PublicKey")
	}
}

// TestLoadKeyPairNonExistent tests LoadKeyPair with a non-existent file.
func TestLoadKeyPairNonExistent(t *testing.T) {
	// Define a path for a non-existent key file
	keyPath := "/tmp/non_existent_key_file"

	// Load the key pair
	_, err := LoadKeyPair(keyPath)
	if err == nil {
		t.Error("LoadKeyPair() error = nil, want an error")
	}
}

// TestLoadKeyPairInvalidData tests LoadKeyPair with invalid data.
func TestLoadKeyPairInvalidData(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()

	// Define a path for the key file
	keyPath := filepath.Join(tmpDir, "invalid_key")

	// Write invalid data to the key file
	err := os.WriteFile(keyPath, []byte("invalid data"), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid data to key file: %v", err)
	}

	// Load the key pair
	_, err = LoadKeyPair(keyPath)
	if err == nil {
		t.Error("LoadKeyPair() error = nil, want an error")
	}
}

// TestLoadKeyPairInvalidLength tests LoadKeyPair with data of invalid length.
func TestLoadKeyPairInvalidLength(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()

	// Define a path for the key file
	keyPath := filepath.Join(tmpDir, "invalid_length_key")

	// Write data of invalid length to the key file
	invalidData := make([]byte, 31) // 31 bytes instead of 32
	for i := range invalidData {
		invalidData[i] = 0xFF
	}
	encodedData := base64.StdEncoding.EncodeToString(invalidData)
	err := os.WriteFile(keyPath, []byte(encodedData), 0600)
	if err != nil {
		t.Fatalf("Failed to write invalid length data to key file: %v", err)
	}

	// Load the key pair
	_, err = LoadKeyPair(keyPath)
	if err == nil {
		t.Error("LoadKeyPair() error = nil, want an error")
	}
}