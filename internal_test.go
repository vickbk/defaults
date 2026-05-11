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

// TestIsTypedNil_Exhaustive provides 100% statement coverage for isTypedNil.
// It exhaustively tests every type in the fast-path switch and the reflection fallback.
func TestIsTypedNil_Exhaustive(t *testing.T) {
	type customStruct struct {
		value int
	}

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		// Literal nil (base case)
		{"literal nil", nil, true},

		// Pointer types - signed integers
		{"*int nil", (*int)(nil), true},
		{"*int non-nil", func() *int { i := 1; return &i }(), false},
		{"*int8 nil", (*int8)(nil), true},
		{"*int8 non-nil", func() *int8 { i := int8(1); return &i }(), false},
		{"*int16 nil", (*int16)(nil), true},
		{"*int16 non-nil", func() *int16 { i := int16(1); return &i }(), false},
		{"*int32 nil", (*int32)(nil), true},
		{"*int32 non-nil", func() *int32 { i := int32(1); return &i }(), false},
		{"*int64 nil", (*int64)(nil), true},
		{"*int64 non-nil", func() *int64 { i := int64(1); return &i }(), false},

		// Pointer types - unsigned integers
		{"*uint nil", (*uint)(nil), true},
		{"*uint non-nil", func() *uint { i := uint(1); return &i }(), false},
		{"*uint8 nil", (*uint8)(nil), true},
		{"*uint8 non-nil", func() *uint8 { i := uint8(1); return &i }(), false},
		{"*uint16 nil", (*uint16)(nil), true},
		{"*uint16 non-nil", func() *uint16 { i := uint16(1); return &i }(), false},
		{"*uint32 nil", (*uint32)(nil), true},
		{"*uint32 non-nil", func() *uint32 { i := uint32(1); return &i }(), false},
		{"*uint64 nil", (*uint64)(nil), true},
		{"*uint64 non-nil", func() *uint64 { i := uint64(1); return &i }(), false},
		{"*uintptr nil", (*uintptr)(nil), true},
		{"*uintptr non-nil", func() *uintptr { i := uintptr(1); return &i }(), false},

		// Pointer types - floating point
		{"*float32 nil", (*float32)(nil), true},
		{"*float32 non-nil", func() *float32 { f := float32(3.14); return &f }(), false},
		{"*float64 nil", (*float64)(nil), true},
		{"*float64 non-nil", func() *float64 { f := 3.14; return &f }(), false},

		// Pointer types - complex
		{"*complex64 nil", (*complex64)(nil), true},
		{"*complex64 non-nil", func() *complex64 { c := complex64(1 + 2i); return &c }(), false},
		{"*complex128 nil", (*complex128)(nil), true},
		{"*complex128 non-nil", func() *complex128 { c := complex128(1 + 2i); return &c }(), false},

		// Pointer types - bool and string
		{"*bool nil", (*bool)(nil), true},
		{"*bool non-nil", func() *bool { b := true; return &b }(), false},
		{"*string nil", (*string)(nil), true},
		{"*string non-nil", func() *string { s := "test"; return &s }(), false},

		// Slice types
		{"[]any nil", ([]any)(nil), true},
		{"[]any empty", []any{}, false},
		{"[]any with values", []any{1, 2}, false},
		{"[]int nil", ([]int)(nil), true},
		{"[]int empty", []int{}, false},
		{"[]int with values", []int{1, 2}, false},
		{"[]string nil", ([]string)(nil), true},
		{"[]string empty", []string{}, false},
		{"[]string with values", []string{"a"}, false},
		{"[]bool nil", ([]bool)(nil), true},
		{"[]bool empty", []bool{}, false},
		{"[]bool with values", []bool{true}, false},

		// Map types
		{"map[string]any nil", (map[string]any)(nil), true},
		{"map[string]any empty", make(map[string]any), false},
		{"map[string]any with values", map[string]any{"key": "value"}, false},
		{"map[string]int nil", (map[string]int)(nil), true},
		{"map[string]int empty", make(map[string]int), false},
		{"map[string]int with values", map[string]int{"key": 1}, false},
		{"map[int]string nil", (map[int]string)(nil), true},
		{"map[int]string empty", make(map[int]string), false},
		{"map[int]string with values", map[int]string{1: "value"}, false},

		// Channel types
		{"chan any nil", (chan any)(nil), true},
		{"chan any created", make(chan any), false},
		{"chan int nil", (chan int)(nil), true},
		{"chan int created", make(chan int), false},
		{"chan string nil", (chan string)(nil), true},
		{"chan string created", make(chan string), false},

		// Function types
		{"func() nil", (func())(nil), true},
		{"func() non-nil", func() {}, false},
		{"func(int) int nil", (func(int) int)(nil), true},
		{"func(int) int non-nil", func(x int) int { return x }, false},

		// Reflection fallback - custom struct pointer (not in fast-path switch)
		{"*customStruct nil", (*customStruct)(nil), true},
		{"*customStruct non-nil", &customStruct{value: 42}, false},

		// Concrete values (non-nil, non-pointer)
		{"concrete int", 42, false},
		{"concrete float", 3.14, false},
		{"concrete bool", true, false},
		{"concrete string", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTypedNil(tt.value)
			if got != tt.want {
				t.Errorf("isTypedNil(%T(%#v)) = %v, want %v", tt.value, tt.value, got, tt.want)
			}
		})
	}
}
