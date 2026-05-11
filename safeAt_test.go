package defaults

import (
	"strings"
	"testing"
)

func TestSafeAt(t *testing.T) {
	var typedNil *int

	tests := []struct {
		name                 string
		values               []any
		index                int
		defaultValue         int
		want                 int
		wantOk               bool
		wantUsedDefault      bool
		expectMessageContain string
	}{
		{
			name:            "Nil slice returns default",
			values:          nil,
			index:           0,
			defaultValue:    7,
			want:            7,
			wantOk:          true,
			wantUsedDefault: true,
		},
		{
			name:            "Negative index returns default",
			values:          []any{42},
			index:           -1,
			defaultValue:    7,
			want:            7,
			wantOk:          true,
			wantUsedDefault: true,
		},
		{
			name:            "Out of bounds returns default",
			values:          []any{42},
			index:           5,
			defaultValue:    7,
			want:            7,
			wantOk:          true,
			wantUsedDefault: true,
		},
		{
			name:            "Valid type returns value",
			values:          []any{42, "ignore"},
			index:           0,
			defaultValue:    7,
			want:            42,
			wantOk:          true,
			wantUsedDefault: false,
		},
		{
			name:            "Type mismatch returns default",
			values:          []any{"not int"},
			index:           0,
			defaultValue:    7,
			want:            7,
			wantOk:          false,
			wantUsedDefault: true,
			expectMessageContain: "invalid type",
		},
		{
			name:            "Typed nil pointer returns default",
			values:          []any{42, typedNil},
			index:           1,
			defaultValue:    7,
			want:            7,
			wantOk:          true,
			wantUsedDefault: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, status := SafeAt(tt.values, tt.index, tt.defaultValue)
			if got != tt.want {
				t.Fatalf("SafeAt(%v, %d, %d) = %d, want %d", tt.values, tt.index, tt.defaultValue, got, tt.want)
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

func BenchmarkSafeAtTypeAssertion(b *testing.B) {
	values := []any{42}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = SafeAt(values, 0, 0)
	}
}

func BenchmarkSafeAtTypedNilFallback(b *testing.B) {
	var typedNil *int
	values := []any{typedNil}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = SafeAt(values, 0, 0)
	}
}
