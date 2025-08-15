package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

// KeyPair holds a public and private key.
type KeyPair struct {
	PublicKey  *[32]byte
	PrivateKey *[32]byte
}

// GenerateKeyPair creates a new public/private key pair.
func GenerateKeyPair() (*KeyPair, error) {
	pubKey, privKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &KeyPair{PublicKey: pubKey, PrivateKey: privKey}, nil
}

// SaveKeyPair saves the private key to a file, base64 encoded.
func SaveKeyPair(kp *KeyPair, path string) error {
	data := base64.StdEncoding.EncodeToString(kp.PrivateKey[:])
	return os.WriteFile(path, []byte(data), 0600)
}

// LoadKeyPair loads a private key from a file and regenerates the key pair.
func LoadKeyPair(path string) (*KeyPair, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privKeyBytes, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	if len(privKeyBytes) != 32 {
		return nil, errors.New("invalid private key length")
	}

	var privKey [32]byte
	copy(privKey[:], privKeyBytes)

	// Re-derive the public key from the private key
	var pubKey [32]byte
	curve25519.ScalarBaseMult(&pubKey, &privKey)

	return &KeyPair{PublicKey: &pubKey, PrivateKey: &privKey}, nil
}
