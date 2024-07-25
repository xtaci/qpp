// # Copyright (c) 2024 xtaci
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
