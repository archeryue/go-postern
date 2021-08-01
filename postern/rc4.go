package postern

import (
	"crypto/md5"
	"io"
)

type sbox [256]byte

type RC4 struct {
	sr *sbox
	sw *sbox
}

func (b *sbox) swap(i, j int) {
	b[i] ^= b[j]
	b[j] ^= b[i]
	b[i] ^= b[j]
}

func NewRC4(key string) *RC4 {
	hash := md5.New()
	io.WriteString(hash, key)
	k := hash.Sum(nil)

	var r sbox
	for i := range r {
		r[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(r[i]) + int(k[i%len(k)])) % 256
		r.swap(i, j)
	}
	var w sbox = r

	return &RC4{
		sr: &r,
		sw: &w,
	}
}

func rc4(s *sbox, data []byte) []byte {
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
	return rc4(c.sw, data)
}

func (c *RC4) Decrypt(data []byte) []byte {
	return rc4(c.sr, data)
}
