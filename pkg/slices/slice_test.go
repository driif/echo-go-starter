package slices_test

import (
	"testing"

	"github.com/driif/echo-go-starter/pkg/slices"
	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	test := []string{"a", "b", "d"}
	assert.True(t, slices.ContainsString(test, "a"))
	assert.True(t, slices.ContainsString(test, "b"))
	assert.False(t, slices.ContainsString(test, "c"))
	assert.True(t, slices.ContainsString(test, "d"))
}

func TestContainsAllString(t *testing.T) {
	test := []string{"a", "b", "d"}
	assert.True(t, slices.ContainsAllString(test, "a"))
	assert.True(t, slices.ContainsAllString(test, "b"))
	assert.False(t, slices.ContainsAllString(test, "c"))
	assert.True(t, slices.ContainsAllString(test, "d"))
	assert.True(t, slices.ContainsAllString(test, "a", "b"))
	assert.True(t, slices.ContainsAllString(test, "a", "d"))
	assert.True(t, slices.ContainsAllString(test, "b", "d"))
	assert.False(t, slices.ContainsAllString(test, "a", "c"))
	assert.False(t, slices.ContainsAllString(test, "b", "c"))
	assert.False(t, slices.ContainsAllString(test, "c", "d"))
	assert.True(t, slices.ContainsAllString(test, "a", "b", "d"))
	assert.False(t, slices.ContainsAllString(test, "a", "b", "c"))
	assert.False(t, slices.ContainsAllString(test, "a", "b", "c", "d"))
	assert.True(t, slices.ContainsAllString(test))
}

func TestUniqueString(t *testing.T) {
	test := []string{"a", "b", "d", "d", "a", "d"}
	assert.Equal(t, []string{"a", "b", "d"}, slices.UniqueString(test))
}
