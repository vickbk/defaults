package defaults

import "testing"

// ============================================================================
// Provider Check Tests
// ============================================================================

func TestValueCheck(t *testing.T) {
	tests := []struct {
		name        string
		defaultVal  int
		input       any
		expectedVal int
		expectOk    bool
		expectUsed  bool
		description string
	}{
		{
			name:        "Exact type match",
			defaultVal:  10,
			input:       42,
			expectedVal: 42,
			expectOk:    true,
			expectUsed:  false,
			description: "Matching int type",
		},
		{
			name:        "Raw nil",
			defaultVal:  10,
			input:       nil,
			expectedVal: 10,
			expectOk:    true,
			expectUsed:  true,
			description: "Raw nil returns default",
		},
		{
			name:        "Type mismatch string vs int",
			defaultVal:  10,
			input:       "not an int",
			expectedVal: 10,
			expectOk:    false,
			expectUsed:  true,
			description: "String input with int default",
		},
		{
			name:        "Type mismatch float vs int",
			defaultVal:  10,
			input:       3.14,
			expectedVal: 10,
			expectOk:    false,
			expectUsed:  true,
			description: "Float input with int default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Value(tt.defaultVal)
			result, status := provider.Check(tt.input, "Custom error")

			if result != tt.expectedVal {
				t.Errorf("%s: Expected value %d, got %d", tt.description, tt.expectedVal, result)
			}

			if status.Ok != tt.expectOk {
				t.Errorf("%s: Expected Ok=%v, got Ok=%v", tt.description, tt.expectOk, status.Ok)
			}

			if status.UsedDefault != tt.expectUsed {
				t.Errorf("%s: Expected UsedDefault=%v, got UsedDefault=%v", tt.description, tt.expectUsed, status.UsedDefault)
			}
		})
	}
}

func TestValueCheckWithTypedNil(t *testing.T) {
	tests := []struct {
		name        string
		input       any
		expectOk    bool
		expectUsed  bool
		description string
	}{
		{
			name:        "Typed nil pointer to int",
			input:       (*int)(nil),
			expectOk:    true,
			expectUsed:  true,
			description: "Reflects on typed nil and returns default",
		},
		{
			name:        "Typed nil pointer to string",
			input:       (*string)(nil),
			expectOk:    true,
			expectUsed:  true,
			description: "Reflects on different typed nil",
		},
		{
			name:        "Typed nil slice",
			input:       ([]int)(nil),
			expectOk:    true,
			expectUsed:  true,
			description: "Typed nil slice",
		},
		{
			name:        "Valid pointer",
			input:       func() *int { i := 42; return &i }(),
			expectOk:    false,
			expectUsed:  true,
			description: "Non-nil pointer to different type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Value(99)
			_, status := provider.Check(tt.input)

			if status.Ok != tt.expectOk {
				t.Errorf("%s: Expected Ok=%v, got Ok=%v", tt.description, tt.expectOk, status.Ok)
			}

			if status.UsedDefault != tt.expectUsed {
				t.Errorf("%s: Expected UsedDefault=%v, got UsedDefault=%v", tt.description, tt.expectUsed, status.UsedDefault)
			}
		})
	}
}

func TestValueCheckErrorMessage(t *testing.T) {
	t.Run("Default error message on type mismatch", func(t *testing.T) {
		provider := Value(42)
		_, status := provider.Check("string")

		if status.Ok {
			t.Errorf("Expected Ok=false for type mismatch")
		}

		if !contains(status.Message, "invalid type") {
			t.Errorf("Expected message to contain 'invalid type', got %q", status.Message)
		}
	})

	t.Run("Custom error message on type mismatch", func(t *testing.T) {
		provider := Value(10)
		customMsg := "Port must be an integer"
		_, status := provider.Check("invalid", customMsg)

		if status.Ok {
			t.Errorf("Expected Ok=false for type mismatch")
		}

		if status.Message != customMsg {
			t.Errorf("Expected custom message %q, got %q", customMsg, status.Message)
		}
	})

	t.Run("Message preserves actual received type", func(t *testing.T) {
		provider := Value(0)
		_, status := provider.Check("wrong")

		if status.Ok {
			t.Errorf("Expected Ok=false for type mismatch")
		}

		if !contains(status.Message, "string") {
			t.Errorf("Expected message to mention 'string', got %q", status.Message)
		}
	})
}

func TestValueCheckStringType(t *testing.T) {
	t.Run("String type checking", func(t *testing.T) {
		provider := Value("default")
		result, status := provider.Check("input", "Must be string")

		if result != "input" {
			t.Errorf("Expected 'input', got %q", result)
		}

		if !status.Ok {
			t.Errorf("Expected Ok=true, got Ok=false")
		}

		if status.UsedDefault {
			t.Errorf("Expected UsedDefault=false, got true")
		}
	})
}

// ============================================================================
// Provider SafeCheck Tests
// ============================================================================

func TestValueSafeCheck(t *testing.T) {
	tests := []struct {
		name        string
		values      []any
		index       int
		defaultVal  int
		expectedVal int
		expectOk    bool
		expectUsed  bool
		description string
	}{
		{
			name:        "In-bounds valid type",
			values:      []any{42, "ignored"},
			index:       0,
			defaultVal:  10,
			expectedVal: 42,
			expectOk:    true,
			expectUsed:  false,
			description: "Correct type at valid index",
		},
		{
			name:        "Out-of-bounds positive",
			values:      []any{42},
			index:       5,
			defaultVal:  10,
			expectedVal: 10,
			expectOk:    true,
			expectUsed:  true,
			description: "Index beyond slice length returns default",
		},
		{
			name:        "Negative index",
			values:      []any{42, "test"},
			index:       -1,
			defaultVal:  10,
			expectedVal: 10,
			expectOk:    true,
			expectUsed:  true,
			description: "Negative index returns default",
		},
		{
			name:        "Type mismatch in bounds",
			values:      []any{"not an int", 42},
			index:       0,
			defaultVal:  10,
			expectedVal: 10,
			expectOk:    false,
			expectUsed:  true,
			description: "Wrong type at index returns error",
		},
		{
			name:        "Empty slice",
			values:      []any{},
			index:       0,
			defaultVal:  10,
			expectedVal: 10,
			expectOk:    true,
			expectUsed:  true,
			description: "Empty slice returns default",
		},
		{
			name:        "Nil at index",
			values:      []any{nil},
			index:       0,
			defaultVal:  10,
			expectedVal: 10,
			expectOk:    true,
			expectUsed:  true,
			description: "nil value at index returns default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := Value(tt.defaultVal)
			result, status := provider.SafeCheck(tt.values, tt.index)

			if result != tt.expectedVal {
				t.Errorf("%s: Expected value %d, got %d", tt.description, tt.expectedVal, result)
			}

			if status.Ok != tt.expectOk {
				t.Errorf("%s: Expected Ok=%v, got Ok=%v", tt.description, tt.expectOk, status.Ok)
			}

			if status.UsedDefault != tt.expectUsed {
				t.Errorf("%s: Expected UsedDefault=%v, got UsedDefault=%v", tt.description, tt.expectUsed, status.UsedDefault)
			}
		})
	}
}
func TestValueSafeCheckWithTypedNil(t *testing.T) {
	t.Run("Typed nil at index", func(t *testing.T) {
		provider := Value(99)
		_, status := provider.SafeCheck([]any{(*int)(nil)}, 0)

		if !status.Ok {
			t.Errorf("Expected Ok=true for typed nil, got Ok=false")
		}

		if !status.UsedDefault {
			t.Errorf("Expected UsedDefault=true for typed nil, got false")
		}
	})

	t.Run("Typed nil at missing index", func(t *testing.T) {
		provider := Value(50)
		result, status := provider.SafeCheck([]any{}, 0)

		if result != 50 {
			t.Errorf("Expected 50, got %d", result)
		}

		if !status.UsedDefault {
			t.Errorf("Expected UsedDefault=true")
		}
	})
}