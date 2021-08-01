package postern

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"sort"
)

type Cipher interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) []byte
}

const (
	Dft = 0
	Rc4 = 1
	Des = 2
	Aes = 3
)

type crTable [256]byte

type cipher struct {
	encTable *crTable
	decTable *crTable
}

func NewCipher(key string, mtd int) Cipher {
	switch mtd {
		case Rc4:
			return NewRC4(key)
		default:
			return defaultCipher(key)
	}
}

func defaultCipher(key string) *cipher {
	hash := md5.New()
	io.WriteString(hash, key)
	buffer := bytes.NewBuffer(hash.Sum(nil))
	var a uint64
	binary.Read(buffer, binary.LittleEndian, &a)

	table := make([]byte, 256, 256)
	for i := range table {
		table[i] = byte(i)
	}

	var t uint64
	for t = 1; t < 1024; t++ {
		sort.SliceStable(table, func(x, y int) bool {
			return (a % (uint64(x) + t)) < (a % (uint64(y) + t))
		})
	}

	enc := &crTable{}
	for i := range table {
		enc[i] = table[i]
	}

	dec := &crTable{}
	for i, v := range table {
		dec[v] = byte(i)
	}

	return &cipher{
		encTable: enc,
		decTable: dec,
	}
}

func (c *cipher) Encrypt(data []byte) []byte {
	ret := make([]byte, len(data), len(data))
	for i, v := range data {
		ret[i] = c.encTable[v]
	}
	return ret
}

func (c *cipher) Decrypt(data []byte) []byte {
	ret := make([]byte, len(data), len(data))
	for i, v := range data {
		ret[i] = c.decTable[v]
	}
	return ret
}
