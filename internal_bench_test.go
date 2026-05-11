package defaults

import (
	"reflect"
	"testing"
)

func BenchmarkIsTypedNil(b *testing.B) {
	benchmarks := []struct {
		name  string
		value any
	}{
		{"literal nil", nil},
		{"typed nil *int", (*int)(nil)},
		{"edge-case struct pointer", (*struct{ x int })(nil)},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name+"_isTypedNil", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = isTypedNil(bm.value)
			}
		})

		b.Run(bm.name+"_reflect", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				v := reflect.ValueOf(bm.value)
				_ = v.IsValid() && v.IsNil()
			}
		})
	}
}
