package crypto

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"

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

const ChunkSize = 64 * 1024

func EncryptLargeFile(inPath, outPath string, masterKey []byte) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}

	defer func() {
		outFile.Sync()
		outFile.Close()
	}()

	if _, err := outFile.Write([]byte("EV1\x00")); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	aead, err := chacha20poly1305.New(masterKey)
	if err != nil {
		return err
	}

	buf := make([]byte, ChunkSize)
	for {
		n, readErr := inFile.Read(buf)
		if n > 0 {
			nonce := make([]byte, chacha20poly1305.NonceSize)
			if _, err := rand.Read(nonce); err != nil {
				return err
			}

			ciphertext := aead.Seal(nil, nonce, buf[:n], nil)

			chunkLen := make([]byte, 4)
			binary.LittleEndian.PutUint32(chunkLen, uint32(len(ciphertext)))
			if _, err := outFile.Write(chunkLen); err != nil {
				return fmt.Errorf("failed to write chunk length: %w", err)
			}

			if _, err := outFile.Write(nonce); err != nil {
				return err
			}
			if _, err := outFile.Write(ciphertext); err != nil {
				return err
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return nil
}
