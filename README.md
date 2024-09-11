# JSON-Go-Optional

The `optional` package provides a generic `Option` type for Go, enabling representation of values that may or may not be present (Some, None, or JsonNull). This package is useful in scenarios where you want to handle optional values without resorting to pointers, and when you need to explicitly differentiate between "no value" and "null."

## Features

- **Generic Option Type**: Create `Option` types for any value, allowing flexibility for optional values.
- **Three States**:
  - **Some**: Represents an existing non-null value.
  - **None**: Represents the absence of a value (field omitted).
  - **JsonNull**: Represents an explicit null value.
- **Database Integration**: Implements `database/sql/driver.Valuer` and `database/sql.Scanner` for database compatibility.
- **JSON Support**: Implements `json.Marshaler` and `json.Unmarshaler` for proper handling of JSON encoding and decoding.

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"optional"
)

func main() {
	// Create an Option with a value
	opt := optional.Some(42)
	if opt.IsSome() {
		fmt.Println("Value:", opt.Unwrap())
	}

	// Create a None Option
	noneOpt := optional.None[int]()
	if noneOpt.IsNone() {
		fmt.Println("No value present.")
	}

	// Create a JsonNull Option
	jsonNullOpt := optional.JsonNull[int]()
	if jsonNullOpt.IsJsonNull() {
		fmt.Println("Explicitly null value.")
	}
}
```

### Common Methods

- `Some(value)`: Wraps a value into an `Option`.
- `None()`: Represents an absent value.
- `JsonNull()`: Represents a `null` value in JSON.
- `Unwrap()`: Extracts the value, returning the zero value of the type if `None` or `JsonNull`.
- `Take()`: Extracts the value and returns an error if `None` or `JsonNull`.
- `TakeOr(fallback)`: Returns the value or a fallback if `None` or `JsonNull`.
- `Filter(predicate)`: Filters the value if it matches a condition.
- `MarshalJSON()` / `UnmarshalJSON()`: Encodes and decodes optional values in JSON.

### JSON Example

```go
package main

import (
	"encoding/json"
	"fmt"
	"optional"
)

func main() {
	opt := optional.Some("hello")
	data, _ := json.Marshal(opt)
	fmt.Println(string(data)) // Output: "hello"

	noneOpt := optional.None[string]()
	data, _ = json.Marshal(noneOpt)
	fmt.Println(string(data)) // Output: nothing (field is omitted)

	jsonNullOpt := optional.JsonNull[string]()
	data, _ = json.Marshal(jsonNullOpt)
	fmt.Println(string(data)) // Output: null
}
```

## Error Handling

- **ErrNoneValueTaken**: Error returned when attempting to extract a value from `None` or `JsonNull`.

## Use Cases

- **Database Fields**: Handling optional database fields with clear semantics for null and omitted fields.
- **APIs**: Managing optional values in API responses without using pointers.
- **Configuration Settings**: Handling optional configuration values in JSON or other structured data formats.

## License

This package is licensed under the MIT License.
