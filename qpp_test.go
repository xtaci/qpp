package qpp

import (
	"crypto/rand"
	"io"
	"testing"

	mathrand "math/rand"
	mathrand2 "math/rand/v2"

	"github.com/stretchr/testify/assert"
)

func TestPads(t *testing.T) {
	numPads := uint16(8)
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)
	qpp := NewQPP(seed, numPads)
	matrixBytes := 1 << 8
	t.Log(qpp.pads)
	t.Log(qpp.rpads)
	assert.Equal(t, len(qpp.pads), len(qpp.rpads), "pads not equal")

	for i := 0; i < int(numPads); i++ {
		pad := qpp.pads[i*matrixBytes : (i+1)*matrixBytes]
		rpad := qpp.rpads[i*matrixBytes : (i+1)*matrixBytes]

		for j := 0; j < matrixBytes; j++ {
			assert.Equal(t, rpad[pad[j]], byte(j), "not reservable")
		}
	}
}

func TestEncryption(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 1024)
	receiver := NewQPP(seed, 1024)

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

	qpp := NewQPP(seed, 1024)

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

	qpp := NewQPP(seed, 1024)

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

func BenchmarkQPP(b *testing.B) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 64)
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

type randSource uint64

func (rs randSource) Uint64() uint64 {
	return uint64(rs)
}

func BenchmarkRandV2(b *testing.B) {
	encRand := mathrand2.New(mathrand2.NewPCG(0xff00, 0x00ff))
	//encRand := mathrand2.New(randSource(0xff00ff00))
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		//println(encRand.Uint32())
		_ = encRand.Uint32()
	}
}
