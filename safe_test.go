package defaults

import (
	"strings"
	"testing"
)

func TestSafe(t *testing.T) {
	nilInt := (*int)(nil)

	tests := []struct {
		name                 string
		values               []any
		defaultValue         int
		message              []string
		want                 int
		wantOk               bool
		wantUsedDefault      bool
		expectMessageContain string
	}{
		{
			name:            "Nil slice returns default",
			values:          nil,
			defaultValue:    7,
			want:            7,
			wantOk:          true,
			wantUsedDefault: true,
		},
		{
			name:            "Valid type returns value",
			values:          []any{42},
			defaultValue:    7,
			want:            42,
			wantOk:          true,
			wantUsedDefault: false,
		},
		{
			name:            "Type mismatch returns default",
			values:          []any{"not int"},
			defaultValue:    7,
			want:            7,
			wantOk:          false,
			wantUsedDefault: true,
			expectMessageContain: "invalid type",
		},
		{
			name:            "Typed nil pointer returns default",
			values:          []any{nilInt},
			defaultValue:    7,
			want:            7,
			wantOk:          true,
			wantUsedDefault: true,
		},
		{
			name:            "Custom error message preserved",
			values:          []any{"bad"},
			defaultValue:    7,
			message:         []string{"custom mismatch"},
			want:            7,
			wantOk:          false,
			wantUsedDefault: true,
			expectMessageContain: "custom mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, status := Safe(tt.values, tt.defaultValue, tt.message...)
			if got != tt.want {
				t.Fatalf("Safe(%v, %d) = %d, want %d", tt.values, tt.defaultValue, got, tt.want)
			}
			if status.Ok != tt.wantOk {
				t.Fatalf("expected Ok=%v, got %v", tt.wantOk, status.Ok)
			}
			if status.UsedDefault != tt.wantUsedDefault {
				t.Fatalf("expected UsedDefault=%v, got %v", tt.wantUsedDefault, status.UsedDefault)
			}
			if tt.expectMessageContain != "" && !strings.Contains(status.Message, tt.expectMessageContain) {
				t.Fatalf("expected status.Message %q to contain %q", status.Message, tt.expectMessageContain)
			}
		})
	}
}

func BenchmarkSafeTypeAssertion(b *testing.B) {
	values := []any{42}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Safe(values, 0)
	}
}

func BenchmarkSafeTypedNilFallback(b *testing.B) {
	var typedNil *int
	values := []any{typedNil}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Safe(values, 0)
	}
}
