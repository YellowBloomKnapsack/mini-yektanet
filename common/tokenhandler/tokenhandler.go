package tokenhandler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"time"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"YellowBloomKnapsack/mini-yektanet/common/models"
)

type TokenHandlerInterface interface {
	GenerateToken(interaction models.AdsInteractionType, adID, publisherID uint, bid int64, key []byte) (string, error)
	VerifyToken(encryptedToken string, key []byte) (*dto.CustomToken, error)
}

type TokenHandlerService struct{}

func NewTokenHandlerService() TokenHandlerInterface {
	return &TokenHandlerService{}
}

func (th *TokenHandlerService) encrypt(data []byte, key []byte) (string, error) {
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

func (th *TokenHandlerService) GenerateToken(interaction models.AdsInteractionType, adID, publisherID uint, bid int64, key []byte) (string, error) {
	token := dto.CustomToken{
		Interaction: interaction,
		AdID:        adID,
		PublisherID: publisherID,
		Bid:         bid,
		CreatedAt:   time.Now().Unix(),
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	encryptedToken, err := th.encrypt(tokenBytes, key)
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

func (th *TokenHandlerService) decrypt(encryptedData string, key []byte) ([]byte, error) {
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

func (th *TokenHandlerService) VerifyToken(encryptedToken string, key []byte) (*dto.CustomToken, error) {
	tokenBytes, err := th.decrypt(encryptedToken, key)
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
