package defaults

import (
	"reflect"
	"testing"
)

func TestSlice(t *testing.T) {
	var nilSlice []int

	tests := []struct {
		name          string
		values        []int
		defaultValues []int
		want          []int
		wantSameSlice bool
	}{
		{
			name:          "Nil slice padded with defaults",
			values:        nilSlice,
			defaultValues: []int{1, 2, 3},
			want:          []int{1, 2, 3},
		},
		{
			name:          "Short slice padded with defaults",
			values:        []int{10},
			defaultValues: []int{5, 20, 30},
			want:          []int{10, 20, 30},
		},
		{
			name:          "Exact length returns original slice",
			values:        []int{5, 6, 7},
			defaultValues: []int{1, 2, 3},
			want:          []int{5, 6, 7},
			wantSameSlice: true,
		},
		{
			name:          "Longer slice returns original slice",
			values:        []int{8, 9, 10, 11},
			defaultValues: []int{1, 2},
			want:          []int{8, 9, 10, 11},
			wantSameSlice: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slice(tt.values, tt.defaultValues...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Slice(%v, %v) = %v, want %v", tt.values, tt.defaultValues, got, tt.want)
			}

			if tt.wantSameSlice && len(tt.values) > 0 {
				if &got[0] != &tt.values[0] {
					t.Fatalf("expected returned slice to reuse input slice header")
				}
			}
		})
	}
}

func BenchmarkSlice(b *testing.B) {
	values := []int{1, 2, 3, 4, 5}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Slice(values, 1, 2, 3)
	}
}
