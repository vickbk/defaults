package defaults

import "testing"

// ============================================================================
// Integration Tests
// ============================================================================

func TestIntegrationMultipleParams(t *testing.T) {
	t.Run("Multiple SafeCheck calls aggregated", func(t *testing.T) {
		options := []any{42, "config"}

		intVal, intStatus := Value(10).SafeCheck(options, 0)
		strVal, strStatus := Value("default").SafeCheck(options, 1)

		if intVal != 42 || strVal != "config" {
			t.Errorf("Expected (42, 'config'), got (%d, %q)", intVal, strVal)
		}

		err := AggregateErrors(intStatus, strStatus)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("Mixed valid and invalid types", func(t *testing.T) {
		options := []any{42, "invalid_int"}

		intVal, intStatus := Value(10).SafeCheck(options, 0)
		intVal2, intStatus2 := Value(20).SafeCheck(options, 1)

		if !intStatus.Ok {
			t.Errorf("First int check should pass")
		}

		if intVal != 42 {
			t.Errorf("Expected first value 42, got %d", intVal)
		}

		if intStatus2.Ok {
			t.Errorf("Second check should fail (string vs int)")
		}

		if intVal2 != 20 {
			t.Errorf("Should use default on type mismatch, got %d", intVal2)
		}
	})
}

func TestIntegrationWithNormalize(t *testing.T) {
	t.Run("Normalize then SafeCheck", func(t *testing.T) {
		raw := []any{5}
		normalized := Normalize(raw, 3)

		val1, _ := Value(1).SafeCheck(normalized, 0)
		val2, _ := Value(2).SafeCheck(normalized, 1)
		val3, _ := Value(3).SafeCheck(normalized, 2)

		if val1 != 5 || val2 != 2 || val3 != 3 {
			t.Errorf("Expected (5, 2, 3), got (%d, %d, %d)", val1, val2, val3)
		}
	})
}
