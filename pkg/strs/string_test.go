package strs_test

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/driif/echo-go-starter/pkg/strs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandom(t *testing.T) {
	res, err := strs.GenerateRandomBytes(13)
	require.NoError(t, err)
	assert.Len(t, res, 13)

	randString, err := strs.GenerateRandomBase64String(17)
	require.NoError(t, err)
	res, err = base64.StdEncoding.DecodeString(randString)
	require.NoError(t, err)
	assert.Len(t, res, 17)

	randString, err = strs.GenerateRandomHexString(19)
	require.NoError(t, err)
	res, err = hex.DecodeString(randString)
	require.NoError(t, err)
	assert.Len(t, res, 19)

	randString, err = strs.GenerateRandomString(19, []strs.CharRange{strs.CharRangeAlphaLowerCase}, "/%$")
	require.NoError(t, err)
	assert.Len(t, randString, 19)
	for _, r := range randString {
		assert.True(t, (r >= 'a' && r <= 'z') || r == '/' || r == '%' || r == '$')
	}

	randString, err = strs.GenerateRandomString(19, []strs.CharRange{strs.CharRangeAlphaUpperCase}, "^\"")
	require.NoError(t, err)
	assert.Len(t, randString, 19)
	for _, r := range randString {
		assert.True(t, (r >= 'A' && r <= 'Z') || r == '^' || r == '"')
	}

	randString, err = strs.GenerateRandomString(19, []strs.CharRange{strs.CharRangeNumeric}, "")
	require.NoError(t, err)
	assert.Len(t, randString, 19)
	for _, r := range randString {
		assert.True(t, (r >= '0' && r <= '9'))
	}

	_, err = strs.GenerateRandomString(1, nil, "")
	require.Error(t, err)

	randString, err = strs.GenerateRandomString(8, nil, "a")
	require.NoError(t, err)
	assert.Len(t, randString, 8)
	assert.Equal(t, "aaaaaaaa", randString)

}
