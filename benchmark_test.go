package life

import "testing"

func Benchmark64(b *testing.B)  { benchmark(b, 8) }
func Benchmark1k(b *testing.B)  { benchmark(b, 32) }
func Benchmark1M(b *testing.B)  { benchmark(b, 1024) }
func Benchmark16M(b *testing.B) { benchmark(b, 4096) }
func Benchmark64M(b *testing.B) { benchmark(b, 8192) }

func benchmark(b *testing.B, N int) {
	w := MakeBoard(N, N)
	SetRand(w, 0, 0.5)
	w.Advance(1)
	b.SetBytes(int64(N) * int64(N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Advance(1)
	}
}
