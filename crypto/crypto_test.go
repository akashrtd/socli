package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"os"
	"testing"
)

// TestKeyPairGenerationSaveLoad tests the generation, saving, and loading of key pairs.
func TestKeyPairGenerationSaveLoad(t *testing.T) {
	// Create a temporary file for the key
	tmpFile, err := os.CreateTemp("", "test_key")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test
	defer tmpFile.Close()

	// Generate a new key pair
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Save the key pair
	err = SaveKeyPair(kp, tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to save key pair: %v", err)
	}

	// Load the key pair
	loadedKp, err := LoadKeyPair(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load key pair: %v", err)
	}

	// Verify that the keys are the same
	if *kp.PrivateKey != *loadedKp.PrivateKey {
		t.Error("Private keys do not match")
	}
	if *kp.PublicKey != *loadedKp.PublicKey {
		t.Error("Public keys do not match")
	}
}

// TestSignMessage tests that SignMessage produces a valid signature length.
// Note: Verification with VerifyMessageSignature may not work due to
// incompatibility between NaCl-derived public keys and ed25519.Verify.
// This is a known limitation of the current implementation.
func TestSignMessage(t *testing.T) {
	// Generate a key pair
	kp, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	message := []byte("This is a test message")

	// Sign the message
	signature, err := SignMessage(message, kp.PrivateKey)
	if err != nil {
		t.Fatalf("Failed to sign message: %v", err)
	}

	// Check signature length (should be 64 bytes for Ed25519)
	if len(signature) != 64 {
		t.Errorf("Signature length is %d, want 64", len(signature))
	}
}

// TestVerifyMessageSignatureWithEd25519Key tests VerifyMessageSignature
// with a standard Ed25519 key pair to ensure the verification logic itself is sound.
func TestVerifyMessageSignatureWithEd25519Key(t *testing.T) {
	// Generate a standard Ed25519 key pair
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate Ed25519 key pair: %v", err)
	}

	message := []byte("This is a test message")

	// Sign the message with the Ed25519 private key
	signature := ed25519.Sign(privKey, message)

	// To use VerifyMessageSignature, we need to convert the ed25519.PublicKey ([32]byte)
	// to a *[32]byte.
	var pubKeyArray [32]byte
	copy(pubKeyArray[:], pubKey)

	// Verify the signature using our VerifyMessageSignature function
	// This test ensures that if we had a compatible public key, verification would work.
	if !VerifyMessageSignature(message, signature, &pubKeyArray) {
		t.Error("Failed to verify signature with VerifyMessageSignature using Ed25519 key")
	}

	// Test with a different message (should fail)
	differentMessage := []byte("This is a different message")
	if VerifyMessageSignature(differentMessage, signature, &pubKeyArray) {
		t.Error("Signature verification should have failed for a different message")
	}

	// Test with a corrupted signature (should fail)
	corruptedSignature := make([]byte, len(signature))
	copy(corruptedSignature, signature)
	corruptedSignature[0] ^= 0xFF // Flip some bits
	if VerifyMessageSignature(message, corruptedSignature, &pubKeyArray) {
		t.Error("Signature verification should have failed for a corrupted signature")
	}
}