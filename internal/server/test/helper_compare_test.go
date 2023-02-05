package test_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/driif/echo-go-starter/internal/server/test"
	"github.com/driif/echo-go-starter/pkg/times"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareFileHashes(t *testing.T) {
	tmpDir := t.TempDir()
	newFilePath := tmpDir + "example2.jpg"
	// join relative path
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	filePath := filepath.Join(pwd, "testdata", "example.jpg")

	in, err := os.Open(filePath)
	require.NoError(t, err)
	defer in.Close()

	out, err := os.Create(newFilePath)
	require.NoError(t, err)
	defer out.Close()

	_, err = io.Copy(out, in)
	require.NoError(t, err)
	require.FileExists(t, newFilePath)

	test.CompareFileHashes(t, filePath, newFilePath)
}

func TestCompareAllPayload(t *testing.T) {
	uid, err := uuid.Parse("0862573e-6ccb-4684-847d-276d3364e91e")
	assert.NoError(t, err)

	payload := test.GenericPayload{
		"A":   1,
		"B":   "b",
		"C":   2.3,
		"D":   true,
		"E":   "2020-02-01",
		"F":   uid,
		"X_Y": "skiped",
	}
	response := map[string]string{
		"A": "1",
		"B": "b",
		"C": "2.3",
		"D": "true",
		"E": times.Date(2020, 2, 1, time.UTC).String(),
		"F": "0862573e-6ccb-4684-847d-276d3364e91e",
	}

	toSkip := map[string]bool{
		"X_Y": true,
	}
	test.CompareAllPayload(t, payload, response, toSkip)

	payload = test.GenericPayload{
		"a":   1,
		"B":   "b",
		"C":   2.3,
		"d":   true,
		"e":   "2020-02-01",
		"F":   uid,
		"X_Y": "skiped",
	}
	test.CompareAllPayload(t, payload, response, toSkip, func(s string) string {
		return strings.ToUpper(s)
	})
}

func TestCompareAll(t *testing.T) {
	uid, err := uuid.Parse("0862573e-6ccb-4684-847d-276d3364e91e")
	assert.NoError(t, err)

	payload := map[string]string{
		"A":   "1",
		"B":   "b",
		"C":   "2.3",
		"D":   "true",
		"E":   "2020-02-01",
		"F":   uid.String(),
		"X_Y": "skiped",
	}
	response := map[string]string{
		"A": "1",
		"B": "b",
		"C": "2.3",
		"D": "true",
		"E": times.Date(2020, 2, 1, time.UTC).String(),
		"F": "0862573e-6ccb-4684-847d-276d3364e91e",
	}

	toSkip := map[string]bool{
		"X_Y": true,
	}
	test.CompareAll(t, payload, response, toSkip)
}
