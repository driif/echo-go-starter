package db_test

import (
	"testing"

	"github.com/driif/echo-go-starter/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestDBTypeConversions(t *testing.T) {
	i := int64(19)
	res := db.NullIntFromInt64Ptr(&i)
	assert.Equal(t, 19, res.Int)
	assert.True(t, res.Valid)

	res = db.NullIntFromInt64Ptr(nil)
	assert.False(t, res.Valid)

	f := 19.9999
	res2 := db.NullFloat32FromFloat64Ptr(&f)
	assert.Equal(t, float32(19.9999), res2.Float32)
	assert.True(t, res2.Valid)

	res2 = db.NullFloat32FromFloat64Ptr(nil)
	assert.False(t, res2.Valid)
}
