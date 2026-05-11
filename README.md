# defaults

**defaults** is a lightweight, zero-allocation Go library for safely handling optional variadic arguments and struct configuration. It bridges the gap between Go's strict type system and the flexibility of variadic parameters, providing a clean, verb-based API that handles "typed nil" pointers and index safety out of the box.

---

## ⚡ The defaults Edge

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

---

## 🛠 At a Glance: Choosing the Right Tool

| Category    | Use Case                  | Function                           | Performance                                 |
| :---------- | :------------------------ | :--------------------------------- | :------------------------------------------ |
| **Typed**   | First element             | `Get(slice, default)`              | **Highest.** Zero-alloc                     |
| **Typed**   | Specific index            | `At(slice, i, default)`            | **Highest.** Zero-alloc                     |
| **Typed**   | Batch/Padding             | `Slice(slice, ...defaults)`        | **Optimized.** Zero-alloc if length matches |
| **Dynamic** | Validate type (index 0)   | `Safe(slice, default, msg)`        | **Secure.** Reflection fallback only        |
| **Dynamic** | Validate type (any index) | `SafeAt(slice, i, default, msg)`   | **Secure.** Handles Typed Nil Paradox       |
| **Dynamic** | Critical config           | `Required(slice, i, default, msg)` | **Strict.** Panics on type mismatch         |
| **Structs** | Functional options        | `Apply(target, ...opts)`           | **Composable.** Error aggregation           |

---

## 🚀 Usage Guide

### 1. Typed Access (Get, At, Slice)

Best for strictly typed slices. **Zero allocations, no reflection overhead.**

```go
func Configure(modes ...string) {
    primary := defaults.Get(modes, "debug")          // Index 0
    secondary := defaults.At(modes, 1, "standard")   // Specific index
    strategy := defaults.Slice(intervals, 100, 500, 2000)  // Batch defaults
}
```

### 2. Dynamic Validation (SafeAt, Required)

The preferred method for `...any` variadics. Solves the "Typed Nil Paradox."

```go
func Setup(options ...any) error {
    retries, rStatus := defaults.SafeAt(options, 0, 3, "Retries must be int")
    port := defaults.Required(options, 1, 8080)  // Panics if wrong type

    return defaults.AggregateErrors(rStatus)
}
```

### 3. Struct Configuration (Apply)

Functional options pattern with built-in error aggregation.

```go
cfg, err := defaults.Apply(&Config{Port: 80},
    WithPort(9000),
    ValidateRequired,
)
if err != nil {
    return err  // All validation errors joined
}
```

---

## 🛠 API Reference

### Core Functions

| Function                              | Description                  | Use Case                          |
| :------------------------------------ | :--------------------------- | :-------------------------------- |
| `Get[T](slice, default)`              | Returns index 0 or default   | First element access              |
| `At[T](slice, i, default)`            | Returns index i or default   | Specific index, handles negative  |
| `Slice[T](slice, ...defaults)`        | Pads slice to minimum length | Batch defaults, preserves extras  |
| `Safe[T](slice, default, msg)`        | Validates type at index 0    | Single untyped value              |
| `SafeAt[T](slice, i, default, msg)`   | Validates type at index i    | Untyped variadics with reflection |
| `Required[T](slice, i, default, msg)` | Panics on type mismatch      | Critical internal configs         |
| `Apply[T](target, ...opts)`           | Applies functional options   | Struct initialization             |

### Helper Functions

| Function                     | Purpose                 | Status                                                      |
| :--------------------------- | :---------------------- | :---------------------------------------------------------- |
| `Value[T](val T)`            | Create a typed provider | **Deprecated** — Use direct functions(`Safe()`, `SafeAt()`) |
| `Normalize(slice, n)`        | Pad slice with nils     | Legacy — Use `Safe()` or `SafeAt` to avoid overheads        |
| `AggregateErrors(...Result)` | Join multiple Results   | Supported — Use with batch validation                       |

---

## 📚 Detailed Examples

### Basic Fallbacks: Get & At

```go
func ConfigureServer(ports ...int) {
    primary := defaults.Get(ports, 8080)
    secondary := defaults.At(ports, 1, 8443)
    tertiary := defaults.At(ports, 10, 9000)  // Safe, no panic
    println(primary, secondary, tertiary)
}

// Usage:
ConfigureServer(9000)           // Ports: 9000, 8443, 9000
ConfigureServer(9000, 9443)     // Ports: 9000, 9443, 9000
ConfigureServer()               // Ports: 8080, 8443, 9000
```

### Batch Defaults: Slice

```go
func SetRetryStrategy(intervals ...int) {
    strategy := defaults.Slice(intervals, 100, 500, 2000)
    initial, secondary, tertiary := strategy[0], strategy[1], strategy[2]
}

// Usage:
SetRetryStrategy()                            // [100, 500, 2000]
SetRetryStrategy(50)                          // [50, 500, 2000]
SetRetryStrategy(50, 200, 1000, 5000)         // [50, 200, 1000, 5000] (zero-alloc)
```

### Typed-Nil Protection: SafeAt (prefer `Apply(&Struct)` for type safety)

```go
func ProcessData(options ...any) error {
    timeout, tStatus := defaults.SafeAt(options, 0, 30, "Timeout must be an int")
    if !tStatus.Ok {
        return tStatus
    }

    retries, rStatus := defaults.SafeAt(options, 1, 3, "Retries must be an int")
    if !rStatus.Ok {
        return rStatus
    }

    return nil
}
```

### Struct Configuration: Apply

```go
type Config struct { Port int; Host string }

func WithPort(p int) defaults.Applier[Config] {
    return func(c *Config) error {
        if p <= 0 || p > 65535 {
            return fmt.Errorf("invalid port: %d", p)
        }
        c.Port = p
        return nil
    }
}

func ValidateRequired(c *Config) error {
    if c.Port == 0 {
        return fmt.Errorf("port required")
    }
    return nil
}

// Usage:
cfg, err := defaults.Apply(
    &Config{Host: "localhost", Port: 8080},
    WithPort(9000),
    ValidateRequired,
)
if err != nil {
    return err  // errors.Join aggregates all failures
}
```

### Error Aggregation

```go
func Setup(options ...any) error {
    retries, rStatus := defaults.SafeAt(options, 0, 3, "Retries must be int")
    timeout, tStatus := defaults.SafeAt(options, 1, 30, "Timeout must be int")
    host, hStatus := defaults.SafeAt(options, 2, "localhost", "Host must be string")

    return defaults.AggregateErrors(rStatus, tStatus, hStatus)
}
```

---

## 🔄 Migration Guide (v0.1 → v0.2)

The new verb-based API replaces the provider pattern. **Old functions still work but are deprecated.**

| v0.1 (Deprecated)                 | v0.2 (Recommended)   | Notes                      |
| :-------------------------------- | :------------------- | :------------------------- |
| `Optional(s, d)`                  | `Get(s, d)`          | Direct replacement         |
| `OptionalAt(s, i, d)`             | `At(s, i, d)`        | Zero-alloc maintained      |
| `Optionals(s, ...d)`              | `Slice(s, ...d)`     | Behavior identical         |
| `Value(d).SafeCheck(s, i)`        | `SafeAt(s, i, d)`    | Removed Provider overhead  |
| `Value(d).SafeCheckOrPanic(s, i)` | `Required(s, i, d)`  | Same semantics             |
| (New)                             | `Apply(target, ...)` | Functional options pattern |
| (New)                             | `Safe(s, d)`         | SafeAt at index 0 only     |

**Migration Example:**

Before:

```go
timeout, err := defaults.Value(30).SafeCheck(options, 0)
```

After:

```go
timeout, err := defaults.SafeAt(options, 0, 30)
```

---

## ⚡ Performance & Constraints

- **Zero-Allocation Paths:** `Get`, `At`, and `Slice` (when length matches) provide zero-alloc paths with no reflection.
- **Interface Boxing:** Using `...any` causes boxing, which can allocate. For hot paths, prefer typed slices with `Get`/`At` or struct initializer `Apply`.
- **Lazy Evaluation:** Error strings are only formatted if a type mismatch occurs.
- **Reflection:** Only used as fallback in `Safe` family to detect "typed nils" (e.g., `(*int)(nil)`).
- **Struct Configuration:** `Apply` pre-allocates error slice capacity for linear iteration over options.

---

## Error Handling

Functions that validate return a `Result` type implementing the `error` interface:

```go
type Result struct {
    Message     string  // Error message (if any)
    Ok          bool    // Validation passed?
    UsedDefault bool    // Default value used?
}
```

Combine multiple results with `AggregateErrors`:

```go
err := defaults.AggregateErrors(status1, status2, status3)
// Returns errors.Join if any status.Ok is false
```

---

## Why `defaults`?

Go's type system and variadics create three challenges:

1. **The Index Panic:** Accessing `args[0]` when empty.
2. **The Typed Nil Paradox:** A `nil` pointer in an interface passes `!= nil` checks but panics on access.
3. **The Boilerplate:** 10 lines of type assertions for one optional parameter.

The `defaults` library solves all three idiomatically, with **zero allocations on the fast path**, respecting Go's core values: simplicity, clarity, and performance.

---

## ⚖️ License

Distributed under the MIT License.
