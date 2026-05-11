package defaults

import (
	"fmt"
	"strings"
	"testing"
)

func TestRequired(t *testing.T) {
	var typedNil *int

	tests := []struct {
		name         string
		values       []any
		index        int
		defaultValue int
		want         int
	}{
		{
			name:         "Valid type returns value",
			values:       []any{42},
			index:        0,
			defaultValue: 7,
			want:         42,
		},
		{
			name:         "Nil slice returns default",
			values:       nil,
			index:        0,
			defaultValue: 7,
			want:         7,
		},
		{
			name:         "Negative index returns default",
			values:       []any{42},
			index:        -1,
			defaultValue: 7,
			want:         7,
		},
		{
			name:         "Out of bounds returns default",
			values:       []any{42},
			index:        5,
			defaultValue: 7,
			want:         7,
		},
		{
			name:         "Typed nil pointer returns default",
			values:       []any{typedNil},
			index:        0,
			defaultValue: 7,
			want:         7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("did not expect panic: %v", r)
				}
			}()

			got := Required(tt.values, tt.index, tt.defaultValue)
			if got != tt.want {
				t.Fatalf("Required(%v, %d, %d) = %d, want %d", tt.values, tt.index, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestRequiredPanicsOnTypeMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on type mismatch")
		} else if !strings.Contains(fmt.Sprint(r), "custom mismatch") {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()

	Required([]any{"not int"}, 0, 7, "custom mismatch")
}

func BenchmarkRequired(b *testing.B) {
	values := []any{42}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Required(values, 0, 0)
	}
}
