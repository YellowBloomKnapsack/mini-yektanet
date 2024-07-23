package tokeninterface

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"YellowBloomKnapsack/mini-yektanet/common/dto"
)

const testKey = "2f4f870d67560a73d39d5db780dfec6d" // Use a simple key for testing

func TestGenerateAndVerifyToken(t *testing.T) {
	// Define sample data
	interaction := dto.ClickType
	adID := uint(123)
	publisherUsername := "testPublisher"
	redirectPath := "/test/path"

	// Generate token
	token, err := GenerateToken(interaction, adID, publisherUsername, redirectPath, []byte(testKey))
	require.NoError(t, err, "expected no error while generating token")
	require.NotEmpty(t, token, "expected token to be generated")

	// Verify token
	verifiedToken, err := VerifyToken(token, []byte(testKey))
	require.NoError(t, err, "expected no error while verifying token")
	require.NotNil(t, verifiedToken, "expected token to be verified")

	// Assert token data
	assert.Equal(t, interaction, verifiedToken.Interaction, "expected interaction type to match")
	assert.Equal(t, adID, verifiedToken.AdID, "expected adID to match")
	assert.Equal(t, publisherUsername, verifiedToken.PublisherUsername, "expected publisherUsername to match")
	assert.Equal(t, redirectPath, verifiedToken.RedirectPath, "expected redirectPath to match")

	// Check created time is within reasonable range (e.g., within the last 5 seconds)
	assert.InDelta(t, time.Now().Unix(), verifiedToken.CreatedAt, 5, "expected createdAt to be within the last 5 seconds")
}
