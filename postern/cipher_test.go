package postern

import "testing"

func TestEncode(t *testing.T) {
	cipher := NewCipher("foobar!")
	raw := []byte("Hello World!")
	data := cipher.Encode(raw)
	msg := cipher.Decode(data)
	for i, v := range raw {
		if v != msg[i] {
			t.Fatalf("%s: encode error at index %d\n", msg, i)
		}
	}
}
