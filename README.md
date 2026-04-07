# Go Data Transformations — Test-Writing Exercise

## Problem Statement

You have been given a data transformation library with **no tests**. The library handles string transformations, number formatting, and date parsing.

Your task is to **write comprehensive table-driven tests** that:
1. Verify all functions work correctly with various inputs
2. Cover edge cases (empty strings, nil-like values, boundary numbers)
3. Catch any bugs in the implementation

**There are bugs hidden in this code.** Your tests should find them.

## Package Structure

```
pkg/transform/
  strings.go    — ToCamelCase, ToSnakeCase, Slugify, Abbreviate
  numbers.go    — FormatCurrency, RoundToDecimal, ParsePercentage
  dates.go      — ParseFlexible, FormatISO, RelativeTime
```

## Functions to Test

### strings.go
- `ToCamelCase(s string) string` — "hello_world" -> "helloWorld"
- `ToSnakeCase(s string) string` — "helloWorld" -> "hello_world"
- `Slugify(s string) string` — "Hello World!" -> "hello-world"
- `Abbreviate(s string, maxLen int) string` — truncate with "..." if over maxLen

### numbers.go
- `FormatCurrency(amount float64, currency string) string` — 1234.56, "USD" -> "$1,234.56"
- `RoundToDecimal(value float64, places int) float64` — 3.14159, 2 -> 3.14
- `ParsePercentage(s string) (float64, error)` — "45.5%" -> 0.455

### dates.go
- `ParseFlexible(s string) (time.Time, error)` — Parses multiple date formats
- `FormatISO(t time.Time) string` — Formats as ISO 8601
- `RelativeTime(t time.Time) string` — "2 hours ago", "in 3 days", etc.

## Writing Tests

Create test files alongside the source:
- `pkg/transform/strings_test.go`
- `pkg/transform/numbers_test.go`
- `pkg/transform/dates_test.go`

Use **table-driven tests** (idiomatic Go):

```go
func TestToCamelCase(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"simple underscore", "hello_world", "helloWorld"},
        {"already camel", "helloWorld", "helloWorld"},
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ToCamelCase(tt.input)
            if got != tt.expected {
                t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
            }
        })
    }
}
```

## Running Tests

```bash
go test -v ./pkg/transform/...
```
