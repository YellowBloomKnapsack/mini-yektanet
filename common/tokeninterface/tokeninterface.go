package tokeninterface

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"
	"errors"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

func encrypt(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func GenerateToken(interaction dto.InteractionType, adID uint, publisherUsername, redirectPath string, key []byte) (string, error) {
	token := dto.CustomToken{
		Interaction:       interaction,
		AdID:              adID,
		PublisherUsername: publisherUsername,
		RedirectPath:      redirectPath,
		CreatedAt:         time.Now().Unix(),
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	encryptedToken, err := encrypt(tokenBytes, key)
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

func decrypt(encryptedData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func VerifyToken(encryptedToken string, key []byte) (*dto.CustomToken, error) {
	tokenBytes, err := decrypt(encryptedToken, key)
	if err != nil {
		return nil, err
	}

	var token dto.CustomToken
	err = json.Unmarshal(tokenBytes, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
