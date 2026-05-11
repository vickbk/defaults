// Package defaults provides zero-allocation utilities for handling optional variadic arguments
// and struct configuration in Go. It bridges the gap between Go's strict type system and
// the flexibility of variadic parameters, providing a clean, verb-based API that handles
// "typed nil" pointers and index safety out of the box.
//
// # Why defaults?
//
// Go's variadic functions (func(...any)) are powerful but require manual nil checks and
// type assertions, leading to verbose and error-prone code. The defaults package solves
// this by providing type-safe, performant alternatives that handle edge cases automatically.
//
// # The Typed Nil Paradox
//
// Go's interface{} values can contain nil pointers that still satisfy != nil:
//
//	var p *int
//	fmt.Println(p == nil)        // true
//	fmt.Println(any(p) == nil)   // false - paradox!
//
// This "Typed Nil Paradox" causes bugs when checking any((*int)(nil)) == nil.
// The defaults package's Safe and Required families use reflection to detect these
// cases, ensuring robust type validation for dynamic arguments.
//
// # Performance
//
// Core functions (Get, At, Slice) execute in sub-nanosecond time with zero heap allocations.
// The Safe family adds minimal overhead (~10-20ns) only when reflection is needed for
// Typed Nil detection. All functions are designed for high-throughput applications.
//
// # Usage Patterns
//
// For strictly typed slices, use Get/At/Slice for maximum performance:
//
//	timeout := defaults.Get(timeouts, 30*time.Second)
//
// For dynamic any variadics, use SafeAt/Required for type safety:
//
//	port, status := defaults.SafeAt(args, 0, 8080, "port must be int")
//	if !status.Ok {
//	    log.Fatal(status.Message)
//	}
//
// For struct configuration, use Apply with functional options:
//
//	cfg, err := defaults.Apply(&Config{}, WithPort(9000), ValidateConfig)
package defaults
