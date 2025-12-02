package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// EncryptFile encrypts data using AES-256-GCM
func EncryptFile(data []byte, key string) ([]byte, error) {
	// Convert hex key to bytes
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key format: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data (nonce is prepended to ciphertext)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptFile decrypts data using AES-256-GCM
func DecryptFile(encryptedData []byte, key string) ([]byte, error) {
	// Convert hex key to bytes
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key format: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GenerateEncryptionKey generates a random 32-byte (256-bit) key in hex format
func GenerateEncryptionKey() string {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return hex.EncodeToString(key)
}
