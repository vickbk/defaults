package defaults

import "testing"

// ============================================================================
// Optional[T] Tests
// ============================================================================

func TestOptional(t *testing.T) {
	tests := []struct {
		name       string
		values     []string
		defaultVal string
		expected   string
		expectUsed bool
	}{
		{
			name:       "Empty slice returns default",
			values:     []string{},
			defaultVal: "fallback",
			expected:   "fallback",
			expectUsed: true,
		},
		{
			name:       "Single element slice returns element",
			values:     []string{"value"},
			defaultVal: "fallback",
			expected:   "value",
			expectUsed: false,
		},
		{
			name:       "Multiple elements returns first",
			values:     []string{"first", "second", "third"},
			defaultVal: "fallback",
			expected:   "first",
			expectUsed: false,
		},
		{
			name:       "Zero value element (empty string) is returned",
			values:     []string{""},
			defaultVal: "fallback",
			expected:   "",
			expectUsed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Optional(tt.values, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestOptionalWithTypedNil(t *testing.T) {
	t.Run("Typed nil pointer in generic slice", func(t *testing.T) {
		values := []*int{nil, nil}
		result := Optional(values, (*int)(nil))

		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})

	t.Run("Typed nil with default fallback", func(t *testing.T) {
		values := []*int{}
		defaultVal := new(int)
		result := Optional(values, defaultVal)

		if result != defaultVal {
			t.Errorf("Expected default fallback, got %v", result)
		}
	})
}

func TestOptionalWithIntegers(t *testing.T) {
	tests := []struct {
		name       string
		values     []int
		defaultVal int
		expected   int
	}{
		{
			name:       "Empty int slice",
			values:     []int{},
			defaultVal: 42,
			expected:   42,
		},
		{
			name:       "Zero value is returned",
			values:     []int{0},
			defaultVal: 42,
			expected:   0,
		},
		{
			name:       "Negative value is returned",
			values:     []int{-100},
			defaultVal: 42,
			expected:   -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Optional(tt.values, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// ============================================================================
// OptionalAt[T] Tests
// ============================================================================

func TestOptionalAt(t *testing.T) {
	values := []string{"a", "b", "c", "d"}
	defaultVal := "default"

	tests := []struct {
		name       string
		index      int
		expected   string
		desc       string
	}{
		{
			name:     "Index 0 in-bounds",
			index:    0,
			expected: "a",
			desc:     "Returns first element",
		},
		{
			name:     "Index 2 in-bounds",
			index:    2,
			expected: "c",
			desc:     "Returns element at middle",
		},
		{
			name:     "Index at length boundary",
			index:    3,
			expected: "d",
			desc:     "Returns last element",
		},
		{
			name:     "Index out of bounds positive",
			index:    4,
			expected: "default",
			desc:     "Returns default when index >= len",
		},
		{
			name:     "Large positive index",
			index:    1000,
			expected: "default",
			desc:     "Returns default for large index",
		},
		{
			name:     "Negative index -1",
			index:    -1,
			expected: "default",
			desc:     "Returns default for negative index",
		},
		{
			name:     "Negative index -100",
			index:    -100,
			expected: "default",
			desc:     "Returns default for large negative index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OptionalAt(values, tt.index, defaultVal)
			if result != tt.expected {
				t.Errorf("%s: Expected %q, got %q", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestOptionalAtEmpty(t *testing.T) {
	t.Run("Empty slice at index 0", func(t *testing.T) {
		result := OptionalAt([]string{}, 0, "fallback")
		if result != "fallback" {
			t.Errorf("Expected fallback, got %q", result)
		}
	})

	t.Run("Empty slice at negative index", func(t *testing.T) {
		result := OptionalAt([]int{}, -1, 99)
		if result != 99 {
			t.Errorf("Expected 99, got %d", result)
		}
	})
}

// ============================================================================
// Optionals[T] Tests
// ============================================================================
