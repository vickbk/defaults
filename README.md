# Go-Defaults

**Go-Defaults** is a nil-safe, type-secure utility for handling optional variadic arguments in Go. It bridges the gap between Go’s strict type system and the flexibility of variadic parameters, providing a clean, generic API that handles "typed nil" pointers and index safety out of the box.

## ⚡ The Go-Defaults Edge

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

**With Go-Defaults:**

```go
timeout := defaults.Optional(args, 30)
```

---

## 🛠 At a Glance: Choosing the Right Tool

| Use Case                     | Recommended Function             | Performance Profile                                  |
| :--------------------------- | :------------------------------- | :--------------------------------------------------- |
| **First** trailing option    | `Optional(slice, default)`       | **Highest.** Zero-alloc                              |
| **Specific** index           | `OptionalAt(slice, i, default)`  | **Highest** Zero-alloc                               |
| **Batch** same-type options  | `Optionals(slice, ...defaults)`  | **Optimized.** Zero-alloc if length matches.         |
| **Mixed** types & validation | `Value(default).SafeCheck(s, i)` | **Secure.** Handles boxing and reflection fallbacks. |

---

## 🚀 Usage Guide

### 1. Targeted Extraction: `OptionalAt` & `Optional`

Best for strictly typed slices. These provide zero-allocation access to specific indices without the need for manual boundary checks.

```go
func Configure(modes ...string) {
    // Grab a specific index; handles out-of-bounds and negative indices safely
    secondary := defaults.OptionalAt(modes, 1, "standard")

    // Shortcut for the first element
    primary := defaults.Optional(modes, "debug")
}
```

### 2. Batch Synchronized Defaults: `Optionals`

Ensures a minimum set of values are present while **preserving any extra values** provided by the user.

```go
func SetRetryStrategy(intervals ...int) {
    // Ensures at least 3 tiers; pads with defaults if missing.
    // Returns original slice (Zero-Alloc) if user provided >= 3 values.
    strategy := defaults.Optionals(intervals, 100, 500, 2000)

    initial, secondary := strategy[0], strategy[1]
}
```

### 3. Type-Safe Validation: `SafeCheck` (Recommended)

The preferred method for multiple optional parameters of different types (`...any`).

> **Important:** `SafeCheck` handles index boundaries internally. You **do not** need to call `Normalize` when using this method, saving an unnecessary slice allocation.

```go
func Setup(options ...any) error {
    // Automatically handles missing indices and type validation
    retries, rStatus := defaults.Value(3).SafeCheck(options, 0, "Retries must be int")
    timeout, tStatus := defaults.Value(30).SafeCheck(options, 1)

    if err := defaults.AggregateErrors(rStatus, tStatus); err != nil {
        return err
    }
    return nil
}
```

### 4. Direct Indexing: `Normalize`

Use `Normalize` only if you require direct index access (`args[i]`) and want to use the standard `Check` method manually.

```go
func CustomLogic(args ...any) {
    params := defaults.Normalize(args, 5) // Pads with nil up to index 4
    val, _ := defaults.Value("data").Check(params[3])
}
```

---

## ⚡ Performance & Constraints

- **Zero-Allocation Paths:** `Optional`, `OptionalAt` and `Optionals` provide zero-alloc paths when the input slice already meets the required length.
- **Interface Boxing:** Using `...any` causes **interface boxing**, which can lead to heap allocations. For ultra-high-frequency hot paths, prefer `Optional` functions with concrete types.
- **Lazy Evaluation:** Error strings and formatting are only computed if a type mismatch actually occurs.
- **Reflection:** `reflect` is only used as a fallback to detect "typed nils" (e.g., `(*int)(nil)`) when standard type assertion fails.

---

## 🛠 API Reference

### Core Functions

| Function                            | Description                                                                 |
| :---------------------------------- | :-------------------------------------------------------------------------- |
| `Optional[T](slice, fallback)`      | Returns index 0 or the fallback value.                                      |
| `OptionalAt[T](slice, i, fallback)` | Returns index `i` or fallback. Negative indexes are considered out of band. |
| `Optionals[T](slice, ...defaults)`  | Pads a slice to a minimum length with specified defaults.                   |
| `Value[T](val T)`                   | Entry point for the generic `Provider` logic.                               |
| `AggregateErrors(...Result)`        | Joins multiple `Result` errors into a single `error`.                       |
| `Normalize(slice, n)`               | Pads a slice to length `n`. **Use only for direct indexing**.               |

### The `Result` Struct

`Result` satisfies the Go `error` interface for direct use in logging or error joining.

- `Message string`: Detailed mismatch description.
- `Ok bool`: False if a type mismatch occurred.
- `UsedDefault bool`: True if the fallback value was utilized.

## ⚖️ License

Distributed under the MIT License.
