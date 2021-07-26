package postern

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"sort"
)

type Cipher interface {
	Encode([]byte, []byte)	
	Decode([]byte, []byte)
}

type crTable [256]byte

type cipher struct {
	encTable *crTable
	decTable *crTable
}

func NewCipher(key string) Cipher {
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

func (c *cipher) Encode(in, out []byte) {
	for i, v := range in {
		out[i] = c.encTable[v]
	}
}

func (c *cipher) Decode(in, out []byte) {
	for i, v := range in {
		out[i] = c.decTable[v]
	}
}
