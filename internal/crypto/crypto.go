package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"io"
)

const nonceSize = 24

func GenerateKeyPair() (*ecdh.PrivateKey, error) {
	return ecdh.X25519().GenerateKey(rand.Reader)
}

func ParsePublicKey(data []byte) (*ecdh.PublicKey, error) {
	return ecdh.X25519().NewPublicKey(data)
}

func Encrypt(data []byte, pub *ecdh.PublicKey) ([]byte, error) {
	ephemeral, err := GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	shared, err := ephemeral.ECDH(pub)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(shared)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return append(ephemeral.PublicKey().Bytes(), ciphertext...), nil
}

func Decrypt(encrypted []byte, priv *ecdh.PrivateKey) ([]byte, error) {
	if len(encrypted) < 32 {
		return nil, errors.New("encrypted data too short")
	}

	ephemeralPub, err := ecdh.X25519().NewPublicKey(encrypted[:32])
	if err != nil {
		return nil, err
	}

	shared, err := priv.ECDH(ephemeralPub)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(shared)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	ciphertext := encrypted[32:]
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)

}
