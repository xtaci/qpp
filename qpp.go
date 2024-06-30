package qpp

import (
	"crypto/sha256"
	"encoding/binary"
)

type QuantumPermutationPad struct {
	pads  [][]byte // encryption
	rpads [][]byte // decryption

	numPads int // number of pads
}

func NewQPP(seed []byte, numPads int) *QuantumPermutationPad {
	qpp := &QuantumPermutationPad{
		numPads: numPads,
	}

	qpp.pads = make([][]byte, numPads)
	qpp.rpads = make([][]byte, numPads)

	for i := 0; i < numPads; i++ {
		qpp.pads[i] = make([]byte, 256)
		qpp.rpads[i] = make([]byte, 256)
		fill(qpp.pads[i])
		shuffle(seed, qpp.pads[i])
		reverse(qpp.pads[i], qpp.rpads[i])
	}

	return qpp
}

func fill(pad []byte) {
	for i := 0; i < 256; i++ {
		pad[i] = byte(i)
	}
}

func reverse(pad []byte, rpad []byte) {
	for i := 0; i < 256; i++ {
		rpad[pad[i]] = byte(i)
	}
}

func shuffle(seed []byte, pad []byte) {
	for i := len(pad) - 1; i > 0; i-- {
		s := sha256.Sum256(seed)
		seed = s[:]

		j := binary.LittleEndian.Uint64(seed) % uint64(i+1)
		pad[i], pad[j] = pad[j], pad[i]
	}
}
