package defaults

import "testing"

func TestAt(t *testing.T) {
	tests := []struct {
		name         string
		values       []int
		index        int
		defaultValue int
		want         int
	}{
		{
			name:         "Empty slice returns default",
			values:       []int{},
			index:        0,
			defaultValue: 99,
			want:         99,
		},
		{
			name:         "Nil slice returns default",
			values:       nil,
			index:        0,
			defaultValue: 123,
			want:         123,
		},
		{
			name:         "Negative index returns default",
			values:       []int{1, 2, 3},
			index:        -1,
			defaultValue: 42,
			want:         42,
		},
		{
			name:         "Out of bounds returns default",
			values:       []int{1, 2, 3},
			index:        5,
			defaultValue: 13,
			want:         13,
		},
		{
			name:         "Within bounds returns value",
			values:       []int{7, 8, 9},
			index:        1,
			defaultValue: 0,
			want:         8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := At(tt.values, tt.index, tt.defaultValue)
			if got != tt.want {
				t.Fatalf("At(%v, %d, %d) = %d, want %d", tt.values, tt.index, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func BenchmarkAt(b *testing.B) {
	values := []int{10, 20, 30, 40, 50}
	defaultValue := 99
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = At(values, 2, defaultValue)
	}
}
