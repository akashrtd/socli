package crypto

import (
	"crypto/ed25519"
)

// SignMessage signs a message using the private key.
// This uses Ed25519 signing, which is standard for libp2p peer IDs.
func SignMessage(message []byte, privKey *[32]byte) ([]byte, error) {
	// Convert [32]byte private key to ed25519.PrivateKey
	// box.PrivateKey is the private key, we need to derive the public key part for Ed25519
	// However, for signing, we typically use the full 64-byte ed25519 private key.
	// If we only have the 32-byte seed, we can generate the full key.
	// For this example, let's assume we have a way to get the full key or use a simple derivation.
	// A more robust solution would involve managing Ed25519 keys separately if signing is distinct from NaCl Box.
	
	// Simple approach: Use the NaCl private key seed to generate an Ed25519 key pair
	// This is generally safe as both use the same underlying Curve25519/Ed25519.
	ed25519PrivKey := ed25519.NewKeyFromSeed(privKey[:])
	
	// Sign the message
	signature := ed25519.Sign(ed25519PrivKey, message)
	return signature, nil
}

// VerifyMessageSignature verifies a message signature using the public key.
func VerifyMessageSignature(message, signature []byte, pubKey *[32]byte) bool {
	// Convert [32]byte public key to ed25519.PublicKey
	ed25519PubKey := ed25519.PublicKey(pubKey[:])
	
	// Verify the signature
	return ed25519.Verify(ed25519PubKey, message, signature)
}