package qpp

import "testing"

func BenchmarkXorShiftStar(b *testing.B) {
	state := uint64(10)
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		state = xorshift64star(state)
	}
}
func BenchmarkXorShift32(b *testing.B) {
	state := uint32(10)
	b.SetBytes(1)
	for i := 0; i < b.N; i++ {
		state = xorshift32(state)
	}
}
