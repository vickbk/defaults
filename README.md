# Go-Defaults

**Go-Defaults** is a nil-safe, type-secure utility for handling optional variadic arguments in Go. It bridges the gap between Go’s strict type system and the flexibility needed for optional parameters by providing a clean, generic API.

## ✨ Features

- **Fluent API:** Clean, readable syntax like `defaults.Value(30).SafeCheck(args, 0)`.
- **Nil-Resilience:** Automatically detects and handles both raw `nil` and "typed nil" pointers (e.g., `(*int)(nil)`).
- **Index Safety:** Prevents "index out of range" panics using `SafeCheck` without needing pre-normalization.
- **Bulk Error Reporting:** Collect and join multiple type-mismatch errors using `AggregateErrors`.
- **Zero Dependencies:** Uses only the Go standard library (including Generics and Reflection).

## 📦 Installation

```bash
go get github.com/vickbk/defaults
```

## 🚀 Quick Start

### 1. Simple Case: The `Optional` Helper

If you only need to handle a **single** optional parameter at the end of a function, use the `Optional` function. It is the most performant way to handle a single fallback without struct overhead.

```go
func Search(query string, tags ...string) {
    // Returns tags[0] if exists, otherwise "all"
    targetTag := defaults.Optional(tags, "all")

    fmt.Println("Searching for:", query, "in tag:", targetTag)
}
```

### 2. Multiple Options: Using `SafeCheck` (Recommended)

When dealing with multiple optional parameters of different types, `SafeCheck` is the preferred method.

> **Performance Note:** You do **not** need to call `Normalize` when using `SafeCheck`. It handles index boundaries internally, saving you an unnecessary slice allocation.

```go
func Setup(name string, options ...any) error {
    // SafeCheck handles missing indices automatically
    retries, rStatus := defaults.Value(3).SafeCheck(options, 0, "Retries must be an int")
    timeout, tStatus := defaults.Value(30).SafeCheck(options, 1)

    // Aggregate errors for a clean feedback loop
    if err := defaults.AggregateErrors(rStatus, tStatus); err != nil {
        return err
    }

    fmt.Printf("Configuring %s: %d retries, %ds timeout\n", name, retries, timeout)
    return nil
}
```

### 3. Direct Indexing: Using `Normalize`

Use `Normalize` **only** if you intend to access the slice indices directly (e.g., `args[i]`) or pass the slice to the standard `Check` method.

```go
func CustomLogic(args ...any) {
    // Pad the slice to ensure args[1] exists as nil if not provided
    params := defaults.Normalize(args, 2)

    val, _ := defaults.Value(true).Check(params[1])
}
```

### Distinguishing Defaults

You can check if a user actually provided a value or if the fallback was used.

```go
val, res := defaults.Value("standard").SafeCheck(args, 0)

if res.UsedDefault {
    fmt.Println("User skipped this field, applied fallback.")
}
```

---

## ⚡ Performance & Constraints

- **Interface Boxing (...any):** Using `...any` (variadic interfaces) introduces a performance gap compared to strictly typed parameters. In Go, passing values as `any` causes **interface boxing**, which often leads to heap allocations. This package is designed for ease of use and safety; for extremely high-frequency hot paths where every nanosecond counts, strictly typed structs are always faster.
- **Allocation Efficiency:** Avoid using `Normalize` if you are already using `SafeCheck`. `SafeCheck` performs a simple length check, whereas `Normalize` allocates a new slice.
- **Lazy Evaluation:** Error messages and `fmt.Sprintf` calls are only executed if a failure is detected. The "happy path" (successful type match) is optimized for speed.
- **Reflection Overhead:** The package uses `reflect` only when a standard type assertion fails to detect "typed nils" (e.g., `(*int)(nil)`). This keeps the common success case fast.

---

## 🛠 API Reference

### Core Functions

| Function                       | Description                                                                        |
| :----------------------------- | :--------------------------------------------------------------------------------- |
| `Optional[T](slice, fallback)` | **Highest Performance.** Best for a single trailing option.                        |
| `Value[T](val T)`              | Creates a `Provider` for type `T` with a fallback value.                           |
| `AggregateErrors(...Result)`   | Joins multiple type-mismatch results into a single `error`.                        |
| `Normalize(slice, n)`          | Pads a slice with `nil` to a minimum length `n`. **Use only for direct indexing.** |

### The `Provider[T]` Methods

| Method                  | Description                                                                  |
| :---------------------- | :--------------------------------------------------------------------------- |
| `Check(input)`          | Validates the input type; handles typed nil pointers via reflection.         |
| `SafeCheck(slice, i)`   | **Preferred.** Boundary-safe check that returns default if index is missing. |
| `SafeCheckOrPanic(...)` | Performs a SafeCheck but panics if a type mismatch occurs.                   |

### The `Result` Struct

`Result` satisfies the Go `error` interface, allowing it to be used directly in logging or error joining.

```go
type Result struct {
    Message     string // Detailed mismatch description (e.g., "expected int, got string")
    Ok          bool   // False only if a type mismatch occurred
    UsedDefault bool   // True if the fallback value was utilized
}
```

---

## ❓ Why use this?

Standard Go requires verbose, repetitive boilerplate to safely handle optional variadic arguments while avoiding panics.

**Before:**

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

**With Go-Defaults:**

```go
timeout, res := defaults.Value(30).SafeCheck(args, 0, "timeout must be int")
```

## ⚖️ License

Distributed under the MIT License.
