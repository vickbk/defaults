# Go-Defaults

**Go-Defaults** is a nil-safe, type-secure utility for handling optional variadic arguments in Go. It bridges the gap between Go’s strict type system and the flexibility needed for optional parameters by providing a clean, generic API.

## ✨ Features

- **Fluent API:** Clean, readable syntax like `defaults.Value(30).SafeCheck(args, 0)`.
- **Nil-Resilience:** Automatically detects and handles both raw `nil` and "typed nil" pointers (e.g., `(*int)(nil)`).
- **Index Safety:** Prevents "index out of range" panics using `SafeCheck` and `Normalize`.
- **Bulk Error Reporting:** Collect and join multiple type-mismatch errors using `AggregateErrors`.
- **Zero Dependencies:** Uses only the Go standard library (including Generics and Reflection).

## 📦 Installation

```bash
go get github.com/vickbk/defaults
```

## 🚀 Quick Start

### Basic Usage

The most common pattern is extracting values from a variadic `...any` slice.

```go
import "github.com/vickbk/defaults"

func StartServer(args ...any) error {
    // 1. Normalize to ensure the slice has at least 2 slots
    params := defaults.Normalize(args, 2)

    // 2. Define defaults and check types at specific indices
    port, r1 := defaults.Value(8080).SafeCheck(params, 0)
    mode, r2 := defaults.Value("production").SafeCheck(params, 1)

    // 3. Aggregate any type-mismatch errors
    if err := defaults.AggregateErrors(r1, r2); err != nil {
        return err
    }

    fmt.Printf("Starting in %s mode on port %d\n", mode, port)
    return nil
}
```

### Distinguishing Defaults

You can check if a user actually provided a value or if the fallback was used.

```go
val, res := defaults.Value("standard").Check(input)

if res.UsedDefault {
    fmt.Println("User skipped this field, applied fallback.")
}
```

## 🛠 API Reference

### Core Functions

| Function                     | Description                                                 |
| :--------------------------- | :---------------------------------------------------------- |
| `Value[T](val T)`            | Creates a `Provider` for type `T` with a fallback value.    |
| `AggregateErrors(...Result)` | Joins multiple type-mismatch results into a single `error`. |
| `Normalize(slice, n)`        | Pads a slice with `nil` to a minimum length `n`.            |
| `Optional(slice, fallback)`  | Returns the first element of a slice or the fallback.       |

### The `Provider[T]` Methods

| Method                  | Description                                                |
| :---------------------- | :--------------------------------------------------------- |
| `Check(input)`          | Validates the input type; handles typed nil pointers.      |
| `SafeCheck(slice, i)`   | Index-safe check that returns default if index is missing. |
| `SafeCheckOrPanic(...)` | Performs a SafeCheck but panics if a type mismatch occurs. |

### The `Result` Struct

`Result` satisfies the Go `error` interface.

```go
type Result struct {
    Message     string // Detailed mismatch description
    Ok          bool   // False if a type mismatch occurred
    UsedDefault bool   // True if the fallback value was utilized
}
```

## ❓ Why use this?

Standard Go requires verbose boilerplate to safely handle optional variadic arguments while avoiding panics.

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
