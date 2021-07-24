package postern

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"io"
	"sort"
)

type Table [256]byte

type Cipher struct {
	encTable *Table
	decTable *Table
}

func NewCipher(key string) *Cipher {
	hash := md5.New()
	io.WriteString(hash, key)
	buffer := bytes.NewBuffer(hash.Sum(nil))
	var a uint64
	binary.Read(buffer, binary.LittleEndian, &a)

	table := make([]byte, 256, 256)
	for i, _ := range table {
		table[i] = byte(i)
	}

	var t uint64
	for t = 1; t < 1024; t++ {
		sort.SliceStable(table, func(x, y int) bool {
			return (a % (uint64(x) + t)) < (a % (uint64(y) + t))
		})
	}

	enc := &Table{}
	for i, _ := range table {
		enc[i] = table[i]
	}

	dec := &Table{}
	for i, v := range table {
		dec[v] = byte(i)
	}

	return &Cipher {
		encTable: enc,
		decTable: dec,
	}
}

func (cipher *Cipher) Encode(data []byte) []byte {
	ret := make([]byte, len(data), len(data))
	for i, v := range data {
		ret[i] = cipher.encTable[v];
	}
	return ret
}

func (cipher *Cipher) Decode(data []byte) []byte {
	ret := make([]byte, len(data), len(data))
	for i, v := range data {
		ret[i] = cipher.decTable[v];
	}
	return ret
}
