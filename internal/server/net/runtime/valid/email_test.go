package valid_test

import (
	"testing"

	"github.com/driif/echo-go-starter/internal/server/net/runtime/valid"
	"github.com/stretchr/testify/assert"
)

func TestValidEmail(t *testing.T) {
	err := valid.Email("tst")
	assert.Error(t, err)
	err = valid.Email("tst@")
	assert.Error(t, err)
	err = valid.Email("tst@tst.")
	assert.Error(t, err)
	err = valid.Email("tst@ts.go")
	assert.NoError(t, err)
}
