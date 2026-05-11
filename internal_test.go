package defaults

import "testing"

func TestIsTypedNil(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{
			name:  "literal nil",
			value: nil,
			want:  true,
		},
		{
			name:  "typed nil *int",
			value: (*int)(nil),
			want:  true,
		},
		{
			name:  "typed nil *string",
			value: (*string)(nil),
			want:  true,
		},
		{
			name:  "typed nil *bool",
			value: (*bool)(nil),
			want:  true,
		},
		{
			name:  "typed nil *float64",
			value: (*float64)(nil),
			want:  true,
		},
		{
			name:  "typed nil []any",
			value: ([]any)(nil),
			want:  true,
		},
		{
			name:  "typed nil map[string]any",
			value: (map[string]any)(nil),
			want:  true,
		},
		{
			name:  "typed nil chan any",
			value: (chan any)(nil),
			want:  true,
		},
		{
			name:  "non-nil *int",
			value: func() *int { i := 42; return &i }(),
			want:  false,
		},
		{
			name:  "concrete int",
			value: 42,
			want:  false,
		},
		{
			name:  "concrete string",
			value: "hello",
			want:  false,
		},
		{
			name:  "concrete bool",
			value: true,
			want:  false,
		},
		{
			name:  "non-nil slice",
			value: []int{1, 2, 3},
			want:  false,
		},
		{
			name:  "non-nil map",
			value: map[string]int{"key": 1},
			want:  false,
		},
		{
			name:  "non-nil chan",
			value: make(chan int),
			want:  false,
		},
		{
			name:  "function (non-nil)",
			value: func() {},
			want:  false,
		},
		{
			name:  "typed nil func()",
			value: (func())(nil),
			want:  true,
		},
		{
			name:  "interface with nil",
			value: func() any { var p *int; return p }(),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTypedNil(tt.value)
			if got != tt.want {
				t.Errorf("isTypedNil(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
