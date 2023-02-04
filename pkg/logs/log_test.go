package logs_test

import (
	"testing"

	"github.com/driif/echo-go-starter/pkg/logs"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogLevelFromString(t *testing.T) {
	res := logs.LogLevelFromString("panic")
	assert.Equal(t, zerolog.PanicLevel, res)

	res = logs.LogLevelFromString("warn")
	assert.Equal(t, zerolog.WarnLevel, res)

	res = logs.LogLevelFromString("foo")
	assert.Equal(t, zerolog.DebugLevel, res)
}
