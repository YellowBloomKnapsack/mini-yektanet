package tokenhandler

import (
	"encoding/base64"
	"testing"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
	"github.com/stretchr/testify/assert"
)

const privateKey = "2f4f870d67560a73d39d5db780dfec6d" // Use a simple key for testing

func TestGenerateToken_Success(t *testing.T) {
	th := NewTokenHandlerService()
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	interaction := dto.ClickType
	adID := uint(1)
	publisherID := uint(100)
	bid := int64(20000)
	redirectPath := "http://example.com"

	token, err := th.GenerateToken(interaction, adID, publisherID, bid, redirectPath, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestVerifyToken_Success(t *testing.T) {
	th := NewTokenHandlerService()
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	interaction := dto.ClickType
	adID := uint(1)
	publisherID := uint(100)
	bid := int64(20000)
	redirectPath := "http://example.com"

	// Generate token
	token, err := th.GenerateToken(interaction, adID, publisherID, bid, redirectPath, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	verifiedToken, err := th.VerifyToken(token, key)
	assert.NoError(t, err)
	assert.NotNil(t, verifiedToken)
	assert.Equal(t, interaction, verifiedToken.Interaction)
	assert.Equal(t, adID, verifiedToken.AdID)
	assert.Equal(t, publisherID, verifiedToken.PublisherID)
	assert.Equal(t, bid, verifiedToken.Bid)
	assert.Equal(t, redirectPath, verifiedToken.RedirectPath)
}

func TestGenerateToken_EncryptError(t *testing.T) {
	th := NewTokenHandlerService()
	key := []byte("short key") // Short key to cause error in encryption

	interaction := dto.ClickType
	adID := uint(1)
	publisherID := uint(100)
	bid := int64(20000)
	redirectPath := "http://example.com"

	token, err := th.GenerateToken(interaction, adID, publisherID, bid, redirectPath, key)
	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestVerifyToken_DecryptError(t *testing.T) {
	th := NewTokenHandlerService()
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	invalidToken := "invalidTokenString"

	verifiedToken, err := th.VerifyToken(invalidToken, key)
	assert.Error(t, err)
	assert.Nil(t, verifiedToken)
}
