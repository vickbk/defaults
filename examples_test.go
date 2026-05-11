package defaults_test

import (
	"errors"
	"fmt"

	"github.com/vickbk/defaults"
)

// ExampleGet demonstrates basic usage of Get with a slice.
func ExampleGet() {
	ports := []int{8080, 8443}
	port := defaults.Get(ports, 3000)
	fmt.Println(port)

	empty := []int{}
	port = defaults.Get(empty, 3000)
	fmt.Println(port)
	// Output: 8080
	// 3000
}

// ExampleApply demonstrates struct configuration using Apply with functional options.
func ExampleApply() {
	type Config struct {
		Port int
		Host string
	}

	WithPort := func(port int) defaults.Applier[Config] {
		return func(c *Config) error {
			if port <= 0 {
				return fmt.Errorf("invalid port: %d", port)
			}
			c.Port = port
			return nil
		}
	}

	ValidateConfig := func(c *Config) error {
		if c.Host == "" {
			return errors.New("host is required")
		}
		return nil
	}

	cfg, err := defaults.Apply(&Config{Host: "localhost"}, WithPort(9000), ValidateConfig)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Host: %s, Port: %d\n", cfg.Host, cfg.Port)
	// Output: Host: localhost, Port: 9000
}

// ExampleSafeAt_typedNil demonstrates Typed Nil safety where the library correctly
// falls back to a default when a nil pointer is passed in an interface.
func ExampleSafeAt_typedNil() {
	// Create a typed nil pointer
	var nilPtr *int

	// Pass it through SafeAt - should return default due to Typed Nil Paradox
	val, status := defaults.SafeAt([]any{nilPtr}, 0, 42, "must be int")
	fmt.Printf("Value: %d, UsedDefault: %t, Ok: %t\n", val, status.UsedDefault, status.Ok)

	// Compare with direct nil check (which would fail)
	fmt.Println(nilPtr == nil)
	fmt.Println(any(nilPtr) == nil)
	// Output: Value: 42, UsedDefault: true, Ok: true
	// true
	// false
}
