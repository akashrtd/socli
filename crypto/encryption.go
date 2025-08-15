package crypto

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/nacl/box"
)

// Encrypt secures a message for a recipient using their public key.
func Encrypt(msg []byte, recipientPubKey *[32]byte, senderPrivKey *[32]byte) ([]byte, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	encrypted := box.Seal(nonce[:], msg, &nonce, recipientPubKey, senderPrivKey)
	return encrypted, nil
}

// Decrypt opens a sealed message using the recipient's private key.
func Decrypt(encrypted []byte, senderPubKey *[32]byte, recipientPrivKey *[32]byte) ([]byte, bool) {
	var nonce [24]byte
	copy(nonce[:], encrypted[:24])
	decrypted, ok := box.Open(nil, encrypted[24:], &nonce, senderPubKey, recipientPrivKey)
	return decrypted, ok
}
