package crypto

import (
	"crypto/ed25519"
	"testing"
)

// TestSignMessageIdentity tests the SignMessage function.
func TestSignMessageIdentity(t *testing.T) {
	// Generate a key pair for testing
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Convert ed25519 private key to *[32]byte
	// ed25519.PrivateKey is 64 bytes, the first 32 are the seed
	var privKey32 [32]byte
	copy(privKey32[:], privKey.Seed())

	message := []byte("This is a test message")

	// Sign the message
	signature, err := SignMessage(message, &privKey32)
	if err != nil {
		t.Fatalf("SignMessage() error = %v, want nil", err)
	}

	// Verify the signature using standard ed25519.Verify
	if !ed25519.Verify(pubKey, message, signature) {
		t.Error("SignMessage() produced an invalid signature")
	}
}

// TestVerifyMessageSignature tests the VerifyMessageSignature function.
func TestVerifyMessageSignature(t *testing.T) {
	// Generate a key pair for testing
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Convert ed25519 public key to *[32]byte
	// ed25519.PublicKey is 32 bytes
	var pubKey32 [32]byte
	copy(pubKey32[:], pubKey)

	// Convert ed25519 private key to *[32]byte
	// ed25519.PrivateKey is 64 bytes, the first 32 are the seed
	var privKey32 [32]byte
	copy(privKey32[:], privKey.Seed())

	message := []byte("This is a test message")

	// Sign the message using standard ed25519.Sign
	signature := ed25519.Sign(privKey, message)

	// Verify the signature using our VerifyMessageSignature function
	if !VerifyMessageSignature(message, signature, &pubKey32) {
		t.Error("VerifyMessageSignature() failed to verify a valid signature")
	}

	// Test with a different message (should fail)
	differentMessage := []byte("This is a different message")
	if VerifyMessageSignature(differentMessage, signature, &pubKey32) {
		t.Error("VerifyMessageSignature() should have failed for a different message")
	}

	// Test with a corrupted signature (should fail)
	corruptedSignature := make([]byte, len(signature))
	copy(corruptedSignature, signature)
	corruptedSignature[0] ^= 0xFF // Flip some bits
	if VerifyMessageSignature(message, corruptedSignature, &pubKey32) {
		t.Error("VerifyMessageSignature() should have failed for a corrupted signature")
	}
}