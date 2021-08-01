package postern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDft(t *testing.T) {
	cipher := NewCipher("foobar!", Dft)
	raw := []byte("Hello World!")
	data := cipher.Encrypt(raw)
	msg := cipher.Decrypt(data)
	assert.Equal(t, raw, msg)
}

func TestRc4(t *testing.T) {
	cipher := NewCipher("foobar!", Rc4)
	raw := []byte("Hello World!")
	data := cipher.Encrypt(raw)
	msg := cipher.Decrypt(data)
	assert.Equal(t, raw, msg)
}
