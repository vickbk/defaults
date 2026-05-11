# defaults

A lightweight, zero-allocation Go library for safely handling variadic arguments and struct configuration defaults. It provides a verb-based API that prioritizes performance and clarity, with built-in support for typed slices and dynamic `[]any` validation using reflection when necessary.

## Overview

The `defaults` library solves a fundamental challenge in Go: providing flexible, safe access to optional variadic parameters without sacrificing type safety or performance. Whether you're extracting values from typed slices, validating mixed-type arguments, or configuring structs with functional options, `defaults` gives you the tools to handle these patterns idiomatically.

### ⚡ The defaults Edge

**Standard Go:**

```go
var timeout int = 30
if len(args) > 0 {
    if v, ok := args[0].(int); ok {
        timeout = v
    } else if args[0] != nil {
        return errors.New("timeout must be int")
    }
}
```

**With defaults:**

```go
timeout := defaults.Get(args, 30)
```

### Key Design Principles

- **Zero Allocations for Core Operations**: `Get`, `At`, and `Slice` operate without memory allocation for typed slices.
- **Pure Functions**: No provider structs required for the typed API—direct, composable functions.
- **Type-Safe Validation**: The "Safe" family uses reflection only when necessary for dynamic `[]any` slices.
- **Graceful Fallbacks**: Built-in support for index boundary checks, nil pointers, and the "Typed Nil Paradox."

---

## Quick Start

```go
// Typed slice access (zero-alloc)
msg := defaults.Get(customMessages, "Default Error")
timeout := defaults.At(retries, 1, 8080)

// Batch defaults
strategy := defaults.Slice(intervals, 100, 500, 2000)

// Dynamic validation with reflection
val, status := defaults.SafeAt(options, 0, 10, "Port must be an int")
if !status.Ok {
    return status.Error() // Aggregated error
}

// Struct configuration with functional options
cfg, err := defaults.Apply(&Config{}, WithPort(9000), WithHost("localhost"))
if err != nil {
    return err
}
```

---

## API Reference

### Typed Access: Zero-Allocation Functions

Use these functions with strongly-typed slices (`[]T`). They have no heap allocations and no reflection overhead.

| Function  | Signature                                        | Use Case                            |
| --------- | ------------------------------------------------ | ----------------------------------- |
| **Get**   | `Get[T](values []T, defaultValue T) T`           | Return first element or default     |
| **At**    | `At[T](values []T, index int, defaultValue T) T` | Return element at index or default  |
| **Slice** | `Slice[T](values []T, defaultValues ...T) []T`   | Ensure minimum length with defaults |

### Dynamic/Safe Access: Reflection-Based Functions

Use these functions with `[]any` to validate type correctness at runtime. The "Typed Nil Paradox" handling ensures nil pointers and nil interfaces are detected correctly.

| Function     | Signature                                                                           | Use Case                                 |
| ------------ | ----------------------------------------------------------------------------------- | ---------------------------------------- |
| **Safe**     | `Safe[T](values []any, defaultValue T, message ...string) (T, Result)`              | Validate type at index 0 with reflection |
| **SafeAt**   | `SafeAt[T](values []any, index int, defaultValue T, message ...string) (T, Result)` | Validate type at specific index          |
| **Required** | `Required[T](values []any, index int, defaultValue T, message ...string) T`         | SafeAt with panic on type mismatch       |

### Struct Configuration

| Function           | Signature                                                 | Use Case                           |
| ------------------ | --------------------------------------------------------- | ---------------------------------- |
| **Apply**          | `Apply[T](target *T, appliers ...Applier[T]) (*T, error)` | Apply functional options to struct |
| **Applier** (type) | `type Applier[T] func(*T) error`                          | Functional option function         |

### Helpers

| Function                                    | Purpose                                                   |
| ------------------------------------------- | --------------------------------------------------------- |
| `Value[T](val T) Provider[T]`               | Create a typed provider (legacy; prefer direct functions) |
| `Normalize(values []any, needed int) []any` | Pad slice with nils to minimum length                     |
| `AggregateErrors(args ...Result) error`     | Join multiple `Result` objects into a single error        |

---

## Usage Examples

### Basic Fallbacks: Get & At

Access elements from typed slices with automatic fallback values:

```go
package main

import "github.com/vickbk/defaults"

func ConfigureServer(ports ...int) {
    // Get the first port, default to 8080
    primary := defaults.Get(ports, 8080)

    // Get the second port (index 1), default to 8443
    secondary := defaults.At(ports, 1, 8443)

    // Out-of-bounds indices return the default
    tertiary := defaults.At(ports, 10, 9000) // Safe, no panic

    println(primary, secondary, tertiary)
}
```

### Index-Safe Access: At with Negative Indices

Safely handle boundary conditions:

```go
func ProcessLogs(entries ...string) {
    // Negative indices are out-of-bounds; return default
    first := defaults.At(entries, -1, "No Entry")

    // This is zero-alloc and won't panic
    println(first)
}
```

### Batch Defaults: Slice

Ensure a minimum number of values while preserving user-provided extras:

```go
func SetRetryStrategy(intervals ...int) {
    // If the user provides fewer than 3 values, pad with defaults
    // If they provide 3 or more, return their slice as-is (zero-alloc)
    strategy := defaults.Slice(intervals, 100, 500, 2000)

    initial, secondary, tertiary := strategy[0], strategy[1], strategy[2]
    println(initial, secondary, tertiary)
}

// Example calls:
// SetRetryStrategy()           → [100, 500, 2000] (all defaults)
// SetRetryStrategy(50)         → [50, 500, 2000] (padded)
// SetRetryStrategy(50, 200)    → [50, 200, 2000] (padded)
// SetRetryStrategy(50, 200, 1000, 5000) → [50, 200, 1000, 5000] (user's slice returned as-is)
```

### Typed-Nil Protection: SafeAt

Handle reflection-based nil checks for interface slices and the "Typed Nil Paradox":

```go
func ProcessData(options ...any) error {
    // SafeAt handles typed nil pointers and interface nil values correctly
    // Unlike simple type assertions, it uses reflection to detect nil pointers
    // that would otherwise compare equal to non-nil interface values

    timeout, tStatus := defaults.SafeAt(options, 0, 30, "Timeout must be an int")
    if !tStatus.Ok {
        return tStatus // Status implements error interface
    }

    retries, rStatus := defaults.SafeAt(options, 1, 3, "Retries must be an int")
    if !rStatus.Ok {
        return rStatus
    }

    println(timeout, retries)
    return nil
}
```

### Struct Configuration: Apply with Functional Options

Use the functional options pattern for composable, validated struct initialization:

```go
type Config struct {
    Port     int
    Host     string
    Timeout  time.Duration
    Retries  int
}

// Define Applier functions
func WithPort(port int) defaults.Applier[Config] {
    return func(c *Config) error {
        if port <= 0 || port > 65535 {
            return fmt.Errorf("invalid port: %d", port)
        }
        c.Port = port
        return nil
    }
}

func WithHost(host string) defaults.Applier[Config] {
    return func(c *Config) error {
        if host == "" {
            return fmt.Errorf("host cannot be empty")
        }
        c.Host = host
        return nil
    }
}

func WithTimeout(d time.Duration) defaults.Applier[Config] {
    return func(c *Config) error {
        if d <= 0 {
            return fmt.Errorf("timeout must be positive")
        }
        c.Timeout = d
        return nil
    }
}

func ValidateRequired(c *Config) error {
    if c.Port == 0 {
        return fmt.Errorf("port is required")
    }
    if c.Host == "" {
        return fmt.Errorf("host is required")
    }
    return nil
}

// Usage
func main() {
    cfg, err := defaults.Apply(
        &Config{
            Port:    8080,
            Host:    "localhost",
            Timeout: 30 * time.Second,
            Retries: 3,
        },
        WithPort(9000),           // Override port
        WithTimeout(60 * time.Second),
        ValidateRequired,         // Validation function
    )

    if err != nil {
        // errors.Join aggregates all Applier errors
        fmt.Printf("Configuration errors: %v\n", err)
        return
    }

    fmt.Printf("Config: %+v\n", cfg)
}
```

### Advanced Validation with Error Aggregation

Combine multiple validations and return all errors at once:

```go
func Setup(options ...any) error {
    retries, rStatus := defaults.SafeAt(options, 0, 3, "Retries must be an int")
    timeout, tStatus := defaults.SafeAt(options, 1, 30, "Timeout must be an int")
    host, hStatus := defaults.SafeAt(options, 2, "localhost", "Host must be a string")

    // Aggregate all validation errors; returns nil if all statuses are Ok
    if err := defaults.AggregateErrors(rStatus, tStatus, hStatus); err != nil {
        return err // All errors joined via errors.Join
    }

    println(retries, timeout, host)
    return nil
}
```

---

## Migration Guide: From v0.1 to v0.2 API

The v0.2 API replaces the provider-based Optional pattern with pure, verb-based functions. Below is a mapping of the old API to the new:

| Old API                                            | New API                              | Notes                                                 |
| -------------------------------------------------- | ------------------------------------ | ----------------------------------------------------- |
| `Optional(slice, default)`                         | `Get(slice, default)`                | Preferred shortcut for first element (index 0)        |
| `OptionalAt(slice, i, default)`                    | `At(slice, i, default)`              | Direct index access, zero-alloc                       |
| `Optionals(slice, ...defaults)`                    | `Slice(slice, ...defaults)`          | Ensures minimum length, preserves extras              |
| `Value(default).SafeCheck(options, i, msg)`        | `SafeAt(options, i, default, msg)`   | Direct function, no provider construction             |
| `Value(default).SafeCheckOrPanic(options, i, msg)` | `Required(options, i, default, msg)` | Panics on type mismatch                               |
| (N/A)                                              | `Apply(cfg, appliers...)`            | **New**: Functional options for struct initialization |
| (N/A)                                              | `Safe(options, default, msg)`        | **New**: SafeAt at index 0 only                       |

### Migration Examples

**Before (v0.1):**

```go
func Configure(options ...any) {
    timeout, err := defaults.Value(30).SafeCheck(options, 0)
    if !err.Ok {
        return err
    }
    println(timeout)
}
```

**After (v0.2):**

```go
func Configure(options ...any) {
    timeout, err := defaults.SafeAt(options, 0, 30)
    if !err.Ok {
        return err
    }
    println(timeout)
}
```

---

## Performance Guarantees

### Zero-Allocation Tier

The core API (`Get`, `At`, `Slice`) maintains **zero allocations** for typical use cases:

- `Get` and `At`: Zero allocations, no reflection. Direct slice access.
- `Slice`: Zero allocations when the input slice has length ≥ the number of defaults. Allocates exactly once when padding is needed.

### Reflection Tier

The "Safe" family (`Safe`, `SafeAt`, `Required`) uses reflection **only when type validation is needed** for `[]any` inputs:

- Type assertion is attempted first (fast path).
- Reflection fallback handles typed nil pointers and interface edge cases.
- No reflection occurs for correctly-typed inputs.

### Struct Configuration Tier

`Apply` is optimized for typical functional options patterns:

- Linear iteration over Appliers.
- Errors are aggregated via `errors.Join`.
- No unnecessary allocations beyond error collection.

---

## Error Handling

All functions that perform validation return a `Result` type that implements the `error` interface:

```go
type Result struct {
    Message     string // Error message (if any)
    Ok          bool   // true if validation passed
    UsedDefault bool   // true if default was used
}

// Implement error interface
func (r Result) Error() string {
    return r.Message
}
```

Use `AggregateErrors` to combine multiple Results into a single error:

```go
err := defaults.AggregateErrors(status1, status2, status3)
if err != nil {
    // err is the result of errors.Join with all failed statuses
}
```

---

## Why `defaults`?

Go's type system and variadic arguments create a tension:

1. **Type Safety**: Go is strongly typed, but variadics are often `...any`.
2. **Index Safety**: Accessing variadic arguments without bounds checking risks panics.
3. **Nil Handling**: Go's interface system makes "typed nil" pointers a subtle bug source.
4. **Performance**: Solutions that solve the above should not require allocations.

The `defaults` library bridges this gap with a clean, idiomatic API that respects Go's values: simplicity, clarity, and performance.
