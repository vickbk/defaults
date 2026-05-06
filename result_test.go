package defaults

import "testing"

// ============================================================================
// Result Tests
// ============================================================================

func TestResultError(t *testing.T) {
	t.Run("Result implements error interface", func(t *testing.T) {
		result := Result{
			Message:     "test error",
			Ok:          false,
			UsedDefault: false,
		}

		var _ error = result
		if result.Error() != "test error" {
			t.Errorf("Expected 'test error', got %q", result.Error())
		}
	})

	t.Run("Empty message", func(t *testing.T) {
		result := Result{Message: "", Ok: true}
		if result.Error() != "" {
			t.Errorf("Expected empty string, got %q", result.Error())
		}
	})
}
