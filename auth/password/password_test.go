package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	raw     = "123456789"
	hashPwd string
)

func TestEncodePassword(t *testing.T) {
	hash, err := EncodePassword(raw)
	assert.NoError(t, err)
	hashPwd = hash
	t.Log(hash)
}

func TestValidatePassword(t *testing.T) {
	assert.Equal(t, ValidatePassword(hashPwd, raw), true)
}
