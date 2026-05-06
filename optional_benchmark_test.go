package defaults

import "testing"

// ============================================================================
// AggregateErrors Tests
// ============================================================================

func BenchmarkOptional(b *testing.B) {
	values := []int{1, 2, 3, 4, 5}
	defaultVal := 0

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Optional(values, defaultVal)
	}
}

func BenchmarkOptionalEmpty(b *testing.B) {
	values := []int{}
	defaultVal := 99

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Optional(values, defaultVal)
	}
}

func BenchmarkOptionalAt(b *testing.B) {
	values := []string{"a", "b", "c", "d", "e"}
	defaultVal := "default"

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = OptionalAt(values, 2, defaultVal)
	}
}

func BenchmarkOptionalAtOutOfBounds(b *testing.B) {
	values := []string{"a", "b", "c"}
	defaultVal := "default"

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = OptionalAt(values, 10, defaultVal)
	}
}

func BenchmarkOptionalsNoAlloc(b *testing.B) {
	values := []int{10, 20, 30}
	defaults := []int{1, 2, 3}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Optionals(values, defaults...)
	}
}

func BenchmarkOptionalsWithAlloc(b *testing.B) {
	values := []int{10}
	defaults := []int{1, 2, 3}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Optionals(values, defaults...)
	}
}
