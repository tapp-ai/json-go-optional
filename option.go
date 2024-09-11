package jsongooptional

import (
	"encoding/json"
)

// Option represents an optional value (Some, None, or Null)
type Option[T any] struct {
	value T
	state State // Tracks the state of the option (Some, None, or Null)
}

// State represents the state of the Option
type State int

const (
	NoneState State = iota // Field is not provided
	NullState              // Field is provided but null
	SomeState              // Field has a value
)

// Some creates an Option with a value (Some)
func Some[T any](v T) Option[T] {
	return Option[T]{
		value: v,
		state: SomeState,
	}
}

// None creates an Option with no value (None)
func None[T any]() Option[T] {
	var zero T
	return Option[T]{
		value: zero,
		state: NoneState,
	}
}

// Null creates an Option with an explicitly null value (Null)
func Null[T any]() Option[T] {
	var zero T
	return Option[T]{
		value: zero,
		state: NullState,
	}
}

// IsSome checks if the Option has a value (Some)
func (o Option[T]) IsSome() bool {
	return o.state == SomeState
}

// IsNone checks if the Option is None (not provided)
func (o Option[T]) IsNone() bool {
	return o.state == NullState
}

// IsNull checks if the Option is Null (provided but null)
func (o Option[T]) IsNull() bool {
	return o.state == NullState
}

// MarshalJSON implements the json.Marshaler interface for custom JSON encoding.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.state == NullState {
		var zeroValue T
		return json.Marshal(zeroValue)
	}
	if o.state == NoneState {
		return json.Marshal(nil)
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface for custom JSON decoding.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*o = Null[T]()
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*o = Some(v)
	return nil
}
