package postern

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDest(t *testing.T) {
	dest := "www.google.com:443"
	buf := EncodeDest(dest)
	ret, _ := DecodeDest(buf[3:])
	assert.Equal(t, dest, ret)
}
