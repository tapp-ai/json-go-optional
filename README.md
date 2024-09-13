# `jsonoptional`

The `jsonoptional` package provides a way to handle values that can be present (Some), omitted (None), or explicitly set to null (Null) in a type-safe manner. This is particularly useful when working with JSON data and APIs where fields can be optional or nullable. The package builds on the idea of optional values, allowing for better control and clarity when dealing with nullable or optional fields in Go structures.

## Features

- **Optional Values**: Define values as either Some, None, or Null.
- **Generic Types**: Supports any type, including custom types and structs.
- **JSON Serialization/Deserialization**: Automatically handles marshalling and unmarshalling optional fields in JSON, allowing null values to be explicitly represented.
- **Helper Functions**: Includes utility functions to wrap values, extract values, and handle optionality gracefully.

## Installation

To install this package, use:

```bash
go get github.com/tapp-ai/jsonoptional
```

## Usage

### Defining an Optional Value

To create an optional value, use the `Some`, `None`, or `Null` constructors:

```go
import "github.com/tapp-ai/jsonoptional"

// Some value (non-null)
opt := jsonoptional.Some(42)

// No value (None)
optNone := jsonoptional.None[int]()

// Explicit null value
optNull := jsonoptional.Null[int]()

// Explicity null value when a condition is met
optNullIf := jsonoptional.NullIf(42, true)
```

A particularly useful feature of this package is the `NullIf` function, which allows you to create an optional value that is null if a condition is met:

```go
var t := time.Time{}
optTime := jsonoptional.NullIf(t, t.IsZero())
fmt.Println(optTime.IsNull()) // true
```

### Checking the Value

You can check the state of an optional value using the following methods:

```go
if opt.IsSome() {
    fmt.Println("Value exists:", opt.Unwrap())
}

if optNone.IsNone() {
    fmt.Println("No value present")
}

if optNull.IsNull() {
    fmt.Println("Value is explicitly null")
}
```

### Unwrapping the Value

You can unwrap the value from the `Option`, but if the value is `None` or `Null`, you will receive the default value of the type:

```go
val := opt.Unwrap() // Returns 42
valNone := optNone.Unwrap() // Returns default value (0 for int)
valNull := optNull.Unwrap() // Returns default value (0 for int)
```

### Fallback Values

If you want to provide a fallback in case of `None`, you can use `TakeOr` or `TakeOrElse`:

```go
fallback := optNone.TakeOr(100) // Returns 100 since optNone is None
fallbackElse := optNone.TakeOrElse(func() int { return 200 }) // Returns 200 from the fallback function
```

### JSON Marshalling and Unmarshalling

`Option` values can be automatically marshalled and unmarshalled to/from JSON:

```go
type Example struct {
    Value jsonoptional.Option[int] `json:"value,omitempty"`
}

example := Example{
    Value: jsonoptional.Some(42),
}

jsonData, _ := json.Marshal(example)
fmt.Println(string(jsonData)) // {"value":42}

var exampleNull Example
json.Unmarshal([]byte(`{"value":null}`), &exampleNull)
fmt.Println(exampleNull.Value.IsNull()) // true
```

### Converting to Standard Optional

You can convert the `jsonoptional.Option` to the standard `optional.Option` provided by the `go-optional` package:

```go
standardOpt := opt.ToOptional()
```

## Error Handling

If attempting to retrieve a value from a `None` type, you can handle errors gracefully:

```go
val, err := optNone.Take()
if err != nil {
    fmt.Println("Error:", err) // "none value taken"
}
```

## License

This package is licensed under the MIT License.
