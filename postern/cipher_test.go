package postern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	cipher := NewCipher("foobar!")
	raw := []byte("Hello World!")
	data := make([]byte, len(raw), len(raw))
	cipher.Encode(raw, data)
	msg := make([]byte, len(data), len(data))
	cipher.Decode(data, msg)
	assert.Equal(t, raw, msg)
}
