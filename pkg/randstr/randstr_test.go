package randstr

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkString(b *testing.B) {
	b.Run("new", func(b *testing.B) {
		var s string

		l := 10_000

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s = String(l)
		}

		b.StopTimer()
		b.Log(s)
	})

	b.Run("orig", func(b *testing.B) {
		var s string

		l := 10_000

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			s = _string(l)
		}

		b.StopTimer()
		b.Log(s)
	})
}

var _seed = rand.New(rand.NewSource(time.Now().UnixNano()))

const _charset = "abcdefghijlkmnopqrstuvwxyz"

func _string(length int) string {
	return _stringWithCharset(length, _charset)
}

func _stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[_seed.Intn(len(charset))]
	}

	return string(b)
}
