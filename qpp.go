package qpp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/pbkdf2"
)

const PAD_IDENTIFIER = "QPP_%b"
const PM_SELECTOR_IDENTIFIER = "PERMUTATION_MATRIX_SELECTOR"
const SHUFFLE_SALT = "___QUANTUM_PERMUTATION_PAD_SHUFFLE_SALT___"
const PRNG_SALT = "___QUANTUM_PERMUTATION_PAD_PRNG_SALT___"
const NATIVE_BYTE_LENGTH = 8 // bit
const PBKDF2_LOOPS = 128

type QuantumPermutationPad struct {
	pads  [][]byte // encryption
	rpads [][]byte // decryption

	numPads uint16     // number of pads
	qubits  uint8      // number of quantum bits
	encRand *rand.Rand // random source for pattern selection
	decRand *rand.Rand // random source for pattern selection
}

func NewQPP(seed []byte, numPads uint16, qubits uint8) *QuantumPermutationPad {
	qpp := &QuantumPermutationPad{
		numPads: numPads,
		qubits:  qubits,
	}

	qpp.pads = make([][]byte, numPads)
	qpp.rpads = make([][]byte, numPads)

	for i := uint16(0); i < numPads; i++ {
		qpp.pads[i] = make([]byte, 1<<qubits)
		qpp.rpads[i] = make([]byte, 1<<qubits)
		fill(qpp.pads[i])
		shuffle(seed, qpp.pads[i], i)
		reverse(qpp.pads[i], qpp.rpads[i])
	}

	qpp.encRand = qpp.CreatePRNG(seed)
	qpp.decRand = qpp.CreatePRNG(seed)

	return qpp
}

// Encrypt encrypts data using the Quantum Permutation Pad with default PRNG
func (qpp *QuantumPermutationPad) Encrypt(data []byte) {
	switch qpp.qubits {
	case NATIVE_BYTE_LENGTH:
		for i := 0; i < len(data); i++ {
			index := qpp.encRand.Uint32() % uint32(qpp.numPads)
			pad := qpp.pads[index]
			data[i] = pad[data[i]]
		}
	default:
	}
}

// Decrypt decrypts data using the Quantum Permutation Pad with default PRNG
func (qpp *QuantumPermutationPad) Decrypt(data []byte) {
	switch qpp.qubits {
	case NATIVE_BYTE_LENGTH:
		for i := 0; i < len(data); i++ {
			index := qpp.decRand.Uint32() % uint32(qpp.numPads)
			rpad := qpp.rpads[index]
			data[i] = rpad[data[i]]
		}
	default:
	}
}

// CreatePRNG creates a deterministic pseudo-random number generator based on seed provided
func (qpp *QuantumPermutationPad) CreatePRNG(seed []byte) *rand.Rand {
	// condense entropy from seed to 8 bytes
	mac := hmac.New(sha256.New, seed)
	mac.Write([]byte(PM_SELECTOR_IDENTIFIER))
	sum := mac.Sum(nil)
	dk := pbkdf2.Key(sum, []byte(PRNG_SALT), PBKDF2_LOOPS, 8, sha1.New)
	source := rand.NewSource(int64(binary.LittleEndian.Uint64(dk)))
	return rand.New(source)
}

// EncryptWithPRNG encrypts data using the Quantum Permutation Pad with PRNG
//
//	this function shares the same P-matrices
func (qpp *QuantumPermutationPad) EncryptWithPRNG(data []byte, rand *rand.Rand) {
	switch qpp.qubits {
	case NATIVE_BYTE_LENGTH:
		for i := 0; i < len(data); i++ {
			index := rand.Uint32() % uint32(qpp.numPads)
			pad := qpp.pads[index]
			data[i] = pad[data[i]]
		}
	default:
	}
}

// DecryptWithPRNG decrypts data using the Quantum Permutation Pad with PRNG
//
//	this function shares the same P-matrices
func (qpp *QuantumPermutationPad) DecryptWithPRNG(data []byte, rand *rand.Rand) {
	switch qpp.qubits {
	case NATIVE_BYTE_LENGTH:
		for i := 0; i < len(data); i++ {
			index := rand.Uint32() % uint32(qpp.numPads)
			rpad := qpp.rpads[index]
			data[i] = rpad[data[i]]
		}
	default:
	}
}

func fill(pad []byte) {
	for i := 0; i < len(pad); i++ {
		pad[i] = byte(i)
	}
}

func reverse(pad []byte, rpad []byte) {
	for i := 0; i < len(pad); i++ {
		rpad[pad[i]] = byte(i)
	}
}

func shuffle(seed []byte, pad []byte, padID uint16) {
	message := fmt.Sprintf(PAD_IDENTIFIER, padID)
	mac := hmac.New(sha256.New, seed)
	mac.Write([]byte(message))
	sum := mac.Sum(nil)

	// condense entropy to 8 bytes
	dk := pbkdf2.Key(sum, []byte(SHUFFLE_SALT), PBKDF2_LOOPS, 8, sha1.New)
	source := rand.NewSource(int64(binary.LittleEndian.Uint64(dk)))
	rand.New(source).Shuffle(len(pad), func(i, j int) {
		pad[i], pad[j] = pad[j], pad[i]
	})
}
