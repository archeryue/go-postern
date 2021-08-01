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

func TestRc4MultiWrite(t *testing.T) {
	cipher := NewCipher("foobar!", Rc4)
	raw := []byte("Hello World!")
	buf := make([]byte, len(raw), len(raw))
	data := cipher.Encrypt(raw[:5])
	copy(buf, data)
	data = cipher.Encrypt(raw[5:])
	copy(buf[5:], data)
	data = cipher.Decrypt(buf[:3])
	copy(buf, data)
	data = cipher.Decrypt(buf[3:7])
	copy(buf[3:], data)
	data = cipher.Decrypt(buf[7:])
	copy(buf[7:], data)
	assert.Equal(t, raw, buf)
}
