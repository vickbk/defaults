package defaults

import (
	"errors"
	"testing"
)

// ============================================================================
// Normalize Tests
// ============================================================================

func TestNormalize(t *testing.T) {
	tests := []struct {
		name        string
		values      []any
		needed      int
		expectedLen int
		description string
	}{
		{
			name:        "Shorter slice padded",
			values:      []any{1, "two"},
			needed:      5,
			expectedLen: 5,
			description: "Pads with nil to reach needed length",
		},
		{
			name:        "Exact length",
			values:      []any{1, 2, 3},
			needed:      3,
			expectedLen: 3,
			description: "No padding when lengths match",
		},
		{
			name:        "Longer slice kept as is",
			values:      []any{1, 2, 3, 4, 5},
			needed:      3,
			expectedLen: 5,
			description: "Result length is original length, extra elements are kept",
		},
		{
			name:        "Empty slice padded",
			values:      []any{},
			needed:      3,
			expectedLen: 3,
			description: "Empty slice becomes all nil",
		},
		{
			name:        "Normalize to 0",
			values:      []any{1, 2},
			needed:      0,
			expectedLen: 2,
			description: "Normalize to zero length should return the original slice (which is longer than needed) without modification",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.values, tt.needed)
			if len(result) != tt.expectedLen {
				t.Errorf("%s: Expected length %d, got %d", tt.description, tt.expectedLen, len(result))
			}

			// Verify padding is nil
			if tt.needed > len(tt.values) {
				for i := len(tt.values); i < tt.needed; i++ {
					if result[i] != nil {
						t.Errorf("%s: Expected nil at index %d, got %v", tt.description, i, result[i])
					}
				}
			}
		})
	}
}

func TestNormalizePreservesValues(t *testing.T) {
	original := []any{42, "hello", 3.14}
	result := Normalize(original, 5)

	for i := 0; i < len(original); i++ {
		if result[i] != original[i] {
			t.Errorf("Value at index %d changed: expected %v, got %v", i, original[i], result[i])
		}
	}
}

// ============================================================================
// AggregateErrors Tests
// ============================================================================

func TestAggregateErrors(t *testing.T) {
	tests := []struct {
		name        string
		results     []Result
		expectErr   bool
		description string
	}{
		{
			name:        "All Ok results",
			results:     []Result{{Ok: true}, {Ok: true}, {Ok: true}},
			expectErr:   false,
			description: "No errors when all Ok=true",
		},
		{
			name:        "Single error",
			results:     []Result{{Ok: false, Message: "error1"}, {Ok: true}},
			expectErr:   true,
			description: "Error returned when one Result is failed",
		},
		{
			name:        "Multiple errors",
			results:     []Result{{Ok: false, Message: "error1"}, {Ok: false, Message: "error2"}, {Ok: true}},
			expectErr:   true,
			description: "Errors joined when multiple Results fail",
		},
		{
			name:        "Empty results",
			results:     []Result{},
			expectErr:   false,
			description: "No error from empty results",
		},
		{
			name:        "Single Ok result",
			results:     []Result{{Ok: true}},
			expectErr:   false,
			description: "Single passing result returns nil",
		},
		{
			name:        "Single error result",
			results:     []Result{{Ok: false, Message: "fail"}},
			expectErr:   true,
			description: "Single error result returns error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AggregateErrors(tt.results...)

			if tt.expectErr && err == nil {
				t.Errorf("%s: Expected error, got nil", tt.description)
			}

			if !tt.expectErr && err != nil {
				t.Errorf("%s: Expected nil, got error: %v", tt.description, err)
			}
		})
	}
}

func TestAggregateErrorsMessages(t *testing.T) {
	t.Run("Error messages are preserved", func(t *testing.T) {
		result1 := Result{Ok: false, Message: "First error"}
		result2 := Result{Ok: false, Message: "Second error"}

		err := AggregateErrors(result1, result2)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		errStr := err.Error()
		if !contains(errStr, "First error") {
			t.Errorf("Expected 'First error' in output, got: %s", errStr)
		}
		if !contains(errStr, "Second error") {
			t.Errorf("Expected 'Second error' in output, got: %s", errStr)
		}
	})

	t.Run("errors.Join compatibility", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		result1 := Result{Ok: false, Message: err1.Error()}
		result2 := Result{Ok: false, Message: err2.Error()}

		aggregated := AggregateErrors(result1, result2)

		if aggregated == nil {
			t.Fatal("Expected error, got nil")
		}

		// Should be able to unwrap multiple errors
		if !hasError(aggregated, "error 1") {
			t.Errorf("Aggregated error missing 'error 1'")
		}
		if !hasError(aggregated, "error 2") {
			t.Errorf("Aggregated error missing 'error 2'")
		}
	})
}

// ============================================================================
// Helper Functions
// ============================================================================

// Simple implementation to avoid external dependency
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return len(substr) == 0
}

func hasError(err error, target string) bool {
	if err == nil {
		return false
	}
	if err.Error() == target {
		return true
	}
	// Check if it's a wrapped error
	u, ok := err.(interface{ Unwrap() []error })
	if ok {
		for _, e := range u.Unwrap() {
			if hasError(e, target) {
				return true
			}
		}
	}
	return false
}
