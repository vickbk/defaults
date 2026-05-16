package defaults

import "testing"

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		values       []string
		defaultValue string
		want         string
	}{
		{
			name:         "Empty slice returns default",
			values:       []string{},
			defaultValue: "fallback",
			want:         "fallback",
		},
		{
			name:         "Nil slice returns default",
			values:       nil,
			defaultValue: "fallback",
			want:         "fallback",
		},
		{
			name:         "Returns first element",
			values:       []string{"alpha", "beta"},
			defaultValue: "fallback",
			want:         "alpha",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.values, tt.defaultValue)
			if got != tt.want {
				t.Fatalf("Get(%v, %q) = %q, want %q", tt.values, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	values := []string{"first", "second", "third"}
	defaultValue := "fallback"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Get(values, defaultValue)
	}
}
