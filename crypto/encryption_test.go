package crypto

import (
	"testing"
)

// TestEncryptDecrypt tests the Encrypt and Decrypt functions.
func TestEncryptDecrypt(t *testing.T) {
	// Generate two key pairs: one for the sender and one for the recipient
	senderKeyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate sender key pair: %v", err)
	}

	recipientKeyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate recipient key pair: %v", err)
	}

	// The message to encrypt
	plaintext := []byte("This is a secret message!")

	// Encrypt the message
	// The sender uses their private key and the recipient's public key
	ciphertext, err := Encrypt(plaintext, recipientKeyPair.PublicKey, senderKeyPair.PrivateKey)
	if err != nil {
		t.Fatalf("Encrypt() error = %v, want nil", err)
	}

	// The ciphertext should not be the same as the plaintext
	if string(ciphertext) == string(plaintext) {
		t.Error("Ciphertext is the same as plaintext, encryption failed")
	}

	// Decrypt the message
	// The recipient uses their private key and the sender's public key
	decrypted, ok := Decrypt(ciphertext, senderKeyPair.PublicKey, recipientKeyPair.PrivateKey)
	if !ok {
		t.Fatal("Decrypt() returned false, decryption failed")
	}

	// The decrypted plaintext should match the original plaintext
	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypt() = %q, want %q", string(decrypted), string(plaintext))
	}

	// Test decryption with the wrong key (should fail)
	// Generate another key pair to act as the wrong recipient
	wrongRecipientKeyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate wrong recipient key pair: %v", err)
	}

	// Try to decrypt with the wrong recipient's private key
	_, ok = Decrypt(ciphertext, senderKeyPair.PublicKey, wrongRecipientKeyPair.PrivateKey)
	if ok {
		t.Error("Decrypt() with wrong key returned true, want false")
	}
}