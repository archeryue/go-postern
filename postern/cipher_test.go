package postern

import "testing"

func TestEncode(t *testing.T) {
	cipher := NewCipher("foobar!")
	raw := []byte("Hello World!")
	data := make([]byte, len(raw), len(raw))
	cipher.Encode(raw, data)
	msg := make([]byte, len(data), len(data))
	cipher.Decode(data, msg)
	for i, v := range raw {
		if v != msg[i] {
			t.Fatalf("%s: encode error at index %d\n", msg, i)
		}
	}
}
