package postern

import (
	"crypto/md5"
	"io"
)

type sbox struct {
	s	[256]byte
	i	int
	j	int
}

type RC4 struct {
	es *sbox
	ds *sbox
}

func (b *sbox) swap(i, j int) {
	b.s[i] ^= b.s[j]
	b.s[j] ^= b.s[i]
	b.s[i] ^= b.s[j]
}

func NewRC4(key string) *RC4 {
	hash := md5.New()
	io.WriteString(hash, key)
	k := hash.Sum(nil)

	var e sbox
	for i := range e.s {
		e.s[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(e.s[i]) + int(k[i%len(k)])) % 256
		e.swap(i, j)
	}
	var d sbox = e

	return &RC4 {
		es: &e,
		ds: &d,
	}
}

func rc4(b *sbox, data []byte) []byte {
	ret := make([]byte, len(data), len(data))
	for k, v := range data {
		b.i = (b.i + 1) % 256
		b.j = (b.j + int(b.s[b.i])) % 256
		b.swap(b.i, b.j)
		ret[k] = v ^ b.s[int(b.s[b.i] + b.s[b.j]) % 256]
	}
	return ret
}

func (c *RC4) Encrypt(data []byte) []byte {
	return rc4(c.es, data)
}

func (c *RC4) Decrypt(data []byte) []byte {
	return rc4(c.ds, data)
}
