package main

import (
	"testing"
)

func BenchmarkRun(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doHQ2x("../example/1/demo.png", "../example/1/demo_hq2x.png")
	}
}
