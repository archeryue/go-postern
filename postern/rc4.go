package postern

import (
	"crypto/md5"
	"io"
)

type box [256]byte

type RC4 struct {
	s *box
}

func (b *box) swap(i, j int) {
	b[i] ^= b[j]
	b[j] ^= b[i]
	b[i] ^= b[j]
}

func NewRC4(key string) *RC4 {
	hash := md5.New()
	io.WriteString(hash, key)
	k := hash.Sum(nil)

	var b box
	for i := range b {
		b[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(b[i]) + int(k[i%len(k)])) % 256
		b.swap(i, j)
	}

	return &RC4{
		s: &b,
	}
}

func (c *RC4) sbox() *box {
	var t box
	for i, v := range c.s {
		t[i] = v
	}
	return &t
}

func (c *RC4) rc4(data []byte) []byte {
	s := c.sbox()
	ret := make([]byte, len(data), len(data))
	i, j := 0, 0
	for k, v := range data {
		i = (i + 1) % 256
		j = (j + int(s[i])) % 256
		s.swap(i, j)
		ret[k] = v ^ s[int(s[i]+s[j])%256]
	}
	return ret
}

func (c *RC4) Encrypt(data []byte) []byte {
	return c.rc4(data)
}

func (c *RC4) Decrypt(data []byte) []byte {
	return c.rc4(data)
}
