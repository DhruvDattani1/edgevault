package crypto

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)


func Encrypt(masterKey, plaintext []byte) (nonce []byte, ciphertext []byte, err error) {
	if len(masterKey) != chacha20poly1305.KeySize {
		return nil, nil, fmt.Errorf("invalid master key size: must be 32 bytes")
	}

	aead, err := chacha20poly1305.New(masterKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	nonce = make([]byte, chacha20poly1305.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext = aead.Seal(nil, nonce, plaintext, nil)
	return nonce, ciphertext, nil
}
