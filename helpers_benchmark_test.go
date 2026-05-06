package defaults

import "testing"

// ============================================================================
// Benchmarks for default functions
// ============================================================================

func BenchmarkValueCheck(b *testing.B) {
	provider := Value(42)
	input := 99

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = provider.Check(input)
	}
}

func BenchmarkValueCheckTypeMismatch(b *testing.B) {
	provider := Value(42)
	input := "string"

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = provider.Check(input)
	}
}

func BenchmarkValueSafeCheck(b *testing.B) {
	provider := Value(10)
	values := []any{42, "test", 3.14}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = provider.SafeCheck(values, 0)
	}
}

func BenchmarkValueSafeCheckOutOfBounds(b *testing.B) {
	provider := Value(10)
	values := []any{42}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = provider.SafeCheck(values, 10)
	}
}

func BenchmarkNormalize(b *testing.B) {
	values := []any{1, "two"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Normalize(values, 5)
	}
}

func BenchmarkAggregateErrors(b *testing.B) {
	results := []Result{
		{Ok: false, Message: "error1"},
		{Ok: true},
		{Ok: false, Message: "error2"},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = AggregateErrors(results...)
	}
}
