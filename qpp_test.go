package qpp

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPads(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)
	qpp := NewQPP(seed, 8, 8)
	t.Log(qpp.pads[0])
	t.Log(qpp.rpads[0])
	assert.Equal(t, len(qpp.pads), len(qpp.rpads), "pads not equal")
	for i := 0; i < len(qpp.pads); i++ {
		for j := 0; j < len(qpp.pads[i]); j++ {
			assert.Equal(t, qpp.rpads[i][qpp.pads[i][j]], byte(j), "not reservable")
		}
	}
}

func TestEncryption(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 8, 8)
	receiver := NewQPP(seed, 8, 8)

	original := []byte("hello world")
	msg := []byte("hello world")
	sender.Encrypt(msg)
	assert.NotEqual(t, original, msg, "not encrypted")
	receiver.Decrypt(msg)
	t.Log(msg)
	assert.Equal(t, original, msg, "not equal")
}

func BenchmarkEncryption(b *testing.B) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 8, 8)
	msg := []byte("hello world")
	b.ResetTimer()
	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		sender.Encrypt(msg)
	}
}
