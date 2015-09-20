package life

import "testing"

func Benchmark64(b *testing.B)   { benchmark(b, 8) }
func Benchmark1k(b *testing.B)   { benchmark(b, 32) }
func Benchmark4k(b *testing.B)   { benchmark(b, 64) }
func Benchmark64k(b *testing.B)  { benchmark(b, 265) }
func Benchmark1M(b *testing.B)   { benchmark(b, 1024) }
func Benchmark16M(b *testing.B)  { benchmark(b, 4096) }
func Benchmark64M(b *testing.B)  { benchmark(b, 8192) }
func Benchmark256M(b *testing.B) { benchmark(b, 16384) }

func benchmark(b *testing.B, N int) {
	w := MakeBoard(N, N)
	//SetRand(w, 0, 0.5)
	b.SetBytes(int64(N) * int64(N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Advance(1)
	}
}
