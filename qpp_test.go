package qpp

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQPP(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)
	qpp := NewQPP(seed, 4)
	t.Log(qpp.pads[0])
	t.Log(qpp.rpads[0])
	assert.Equal(t, len(qpp.pads), len(qpp.rpads), "pads not equal")
	for i := 0; i < len(qpp.pads); i++ {
		for j := 0; j < len(qpp.pads[i]); j++ {
			assert.Equal(t, qpp.rpads[i][qpp.pads[i][j]], byte(j), "not reservable")
		}
	}
}
