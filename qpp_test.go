package qpp

import (
	"crypto/rand"
	"io"
	mathrand "math/rand"
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

	sender := NewQPP(seed, 1024, 8)
	receiver := NewQPP(seed, 1024, 8)

	original := make([]byte, 65536)
	io.ReadFull(rand.Reader, original)
	msg := make([]byte, len(original))
	copy(msg, original)
	sender.Encrypt(msg)
	assert.NotEqual(t, original, msg, "not encrypted")
	receiver.Decrypt(msg)
	assert.Equal(t, original, msg, "not equal")
}

func TestEncryption2(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	qpp := NewQPP(seed, 1024, 8)

	original := make([]byte, 65536)
	io.ReadFull(rand.Reader, original)
	msg := make([]byte, len(original))
	copy(msg, original)
	qpp.Encrypt(msg)
	assert.NotEqual(t, original, msg, "not encrypted")
	qpp.Decrypt(msg)
	assert.Equal(t, original, msg, "not equal")
}

func TestEncryptionMixedPRNG(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	qpp := NewQPP(seed, 1024, 8)

	original := make([]byte, 65536)
	io.ReadFull(rand.Reader, original)
	msg := make([]byte, len(original))
	copy(msg, original)

	rand_enc := qpp.CreatePRNG(seed)
	qpp.EncryptWithPRNG(msg, rand_enc)
	assert.NotEqual(t, original, msg, "not encrypted")

	rand_dec := qpp.CreatePRNG(seed)
	qpp.DecryptWithPRNG(msg, rand_dec)
	assert.Equal(t, original, msg, "not equal")
}

func TestSeedToChunk(t *testing.T) {
	seed := "hello quantum world, hello quantum world, hello quantum world"
	t.Log("long seed, 8 qubit:", seedToChunks([]byte(seed), 8))
	shortSeed := "hello"
	t.Log("short seed, 8 qubit", seedToChunks([]byte(shortSeed), 8))
	t.Log("chunk size:", len(seedToChunks([]byte(seed), 8)))
}

func TestQPPMinimumSeedLength(t *testing.T) {
	for i := 1; i < 16; i++ {
		t.Log(i, "qubit -> minimum seed length(bytes):", QPPMinimumSeedLength(uint8(i)), "Min Pads:", QPPMinimumPads(uint8(i)))
	}
}

func BenchmarkEncryption(b *testing.B) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 64, 8)
	msg := make([]byte, 512)
	io.ReadFull(rand.Reader, msg)

	b.ResetTimer()
	b.SetBytes(int64(len(msg)))
	for i := 0; i < b.N; i++ {
		sender.Encrypt(msg)
	}
}

func BenchmarkRand(b *testing.B) {
	encSource := mathrand.NewSource(0xff00ff00)
	encRand := mathrand.New(encSource)
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		_ = encRand.Uint32()
	}
}
