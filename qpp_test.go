package qpp

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"testing"

	mathrand "math/rand"
	mathrand2 "math/rand/v2"

	"github.com/stretchr/testify/assert"
)

const bibleText = `
1:1 In the beginning God created the heaven and the earth.

1:2 And the earth was without form, and void; and darkness was upon
the face of the deep. And the Spirit of God moved upon the face of the
waters.

1:3 And God said, Let there be light: and there was light.

1:4 And God saw the light, that it was good: and God divided the light
from the darkness.

1:5 And God called the light Day, and the darkness he called Night.
And the evening and the morning were the first day.

1:6 And God said, Let there be a firmament in the midst of the waters,
and let it divide the waters from the waters.

1:7 And God made the firmament, and divided the waters which were
under the firmament from the waters which were above the firmament:
and it was so.

1:8 And God called the firmament Heaven. And the evening and the
morning were the second day.

1:9 And God said, Let the waters under the heaven be gathered together
unto one place, and let the dry land appear: and it was so.

1:10 And God called the dry land Earth; and the gathering together of
the waters called he Seas: and God saw that it was good.

1:11 And God said, Let the earth bring forth grass, the herb yielding
seed, and the fruit tree yielding fruit after his kind, whose seed is
in itself, upon the earth: and it was so.

1:12 And the earth brought forth grass, and herb yielding seed after
his kind, and the tree yielding fruit, whose seed was in itself, after
his kind: and God saw that it was good.

1:13 And the evening and the morning were the third day.

1:14 And God said, Let there be lights in the firmament of the heaven
to divide the day from the night; and let them be for signs, and for
seasons, and for days, and years: 1:15 And let them be for lights in
the firmament of the heaven to give light upon the earth: and it was
so.

1:16 And God made two great lights; the greater light to rule the day,
and the lesser light to rule the night: he made the stars also.

1:17 And God set them in the firmament of the heaven to give light
upon the earth, 1:18 And to rule over the day and over the night, and
to divide the light from the darkness: and God saw that it was good.

1:19 And the evening and the morning were the fourth day.

1:20 And God said, Let the waters bring forth abundantly the moving
creature that hath life, and fowl that may fly above the earth in the
open firmament of heaven.

1:21 And God created great whales, and every living creature that
moveth, which the waters brought forth abundantly, after their kind,
and every winged fowl after his kind: and God saw that it was good.

1:22 And God blessed them, saying, Be fruitful, and multiply, and fill
the waters in the seas, and let fowl multiply in the earth.

1:23 And the evening and the morning were the fifth day.

1:24 And God said, Let the earth bring forth the living creature after
his kind, cattle, and creeping thing, and beast of the earth after his
kind: and it was so.

1:25 And God made the beast of the earth after his kind, and cattle
after their kind, and every thing that creepeth upon the earth after
his kind: and God saw that it was good.

1:26 And God said, Let us make man in our image, after our likeness:
and let them have dominion over the fish of the sea, and over the fowl
of the air, and over the cattle, and over all the earth, and over
every creeping thing that creepeth upon the earth.

1:27 So God created man in his own image, in the image of God created
he him; male and female created he them.

1:28 And God blessed them, and God said unto them, Be fruitful, and
multiply, and replenish the earth, and subdue it: and have dominion
over the fish of the sea, and over the fowl of the air, and over every
living thing that moveth upon the earth.

1:29 And God said, Behold, I have given you every herb bearing seed,
which is upon the face of all the earth, and every tree, in the which
is the fruit of a tree yielding seed; to you it shall be for meat.

1:30 And to every beast of the earth, and to every fowl of the air,
and to every thing that creepeth upon the earth, wherein there is
life, I have given every green herb for meat: and it was so.

1:31 And God saw every thing that he had made, and, behold, it was
very good. And the evening and the morning were the sixth day.
`

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

func TestEncryption3(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 1024)
	receiver := NewQPP(seed, 1024)

	original := make([]byte, 12)
	io.ReadFull(rand.Reader, original)
	msg := make([]byte, len(original))
	copy(msg, original)

	// 12 == 3 + 5 + 4
	sender.Encrypt(msg[:3])
	sender.Encrypt(msg[3:8])
	sender.Encrypt(msg[8:])
	//sender.Encrypt(msg)
	assert.NotEqual(t, original, msg, "not encrypted")

	// 12 = 9 + 1 + 2
	receiver.Decrypt(msg[:9])
	receiver.Decrypt(msg[9:10])
	receiver.Decrypt(msg[10:])
	//receiver.Decrypt(msg)
	assert.Equal(t, original, msg, "not equal")
}

// This function test chi-square test for randomness with differnt number of pads by calling testChiSquare
func TestEncryptionChiSquare(t *testing.T) {
	// open a temporary csv file
	f, err := os.Create("chi-square.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// write header to csv
	f.WriteString("pads,chi-square\n")

	// test chi-square for 1 to 255 pads, and output the result to a CSV file
	for i := 1; i < 256; i++ {
		chi := testChiSquare(t, uint16(i))
		f.WriteString(fmt.Sprintf("%d,%f\n", i, chi))
	}
}

func testChiSquare(t *testing.T, pads uint16) float64 {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, pads)
	original := make([]byte, 1024*1024)
	msg := original

	// fill original with bible text
	for len(msg) > 0 {
		n := copy(msg, bibleText)
		msg = msg[n:]
	}

	msg = original
	sender.Encrypt(msg)

	// chi-square test
	// 1. calculate frequency
	freq := make([]int, 256)
	for _, b := range msg {
		freq[b]++
	}
	// 2. calculate expected frequency
	expected := float64(len(msg)) / 256
	// 3. calculate chi-square
	chi := 0.0
	for _, f := range freq {
		chi += (float64(f) - expected) * (float64(f) - expected) / expected
	}
	t.Log("pads:", pads, "chi-square:", chi)
	return chi

}

func TestEncryptionRandLength(t *testing.T) {
	seed := make([]byte, 32)
	io.ReadFull(rand.Reader, seed)

	sender := NewQPP(seed, 1024)
	receiver := NewQPP(seed, 1024)

	original := make([]byte, 1024*1024*1024)
	io.ReadFull(rand.Reader, original)
	msg := make([]byte, len(original))
	copy(msg, original)

	// 12 == 3 + 5 + 4
	start := msg
	for len(start) > 0 {
		l := mathrand.Intn(len(start)+1) % 257
		sender.Encrypt(start[:l])
		start = start[l:]
	}
	//sender.Encrypt(msg)
	assert.NotEqual(t, original, msg, "not encrypted")

	// 12 = 9 + 1 + 2
	start = msg
	for len(start) > 0 {
		l := mathrand.Intn(len(start)+1) % 313
		receiver.Decrypt(start[:l])
		start = start[l:]
	}
	//receiver.Decrypt(msg)
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
