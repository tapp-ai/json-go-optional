package jsonoptional

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrNoneValueTaken represents the error that is raised when None value is taken.
	ErrNoneValueTaken = errors.New("none value taken")
	// NullBytes is a byte slice representation of the string "null"
	NullBytes = []byte("null")
)

// Option is a data type that must be Some (i.e. having a value),None (i.e. doesn't have a value), or Null (i.e. has a value but it's null).
type Option[T any] map[bool]T

// Some is a function to make an Option type value with the actual value.
func Some[T any](v T) Option[T] {
	return Option[T]{
		true: v,
	}
}

// None is a function to make an Option type value that doesn't have a value.
func None[T any]() Option[T] {
	return map[bool]T{}
}

// Null is a function to make an Option type value that has an explicit null value.
func Null[T any]() Option[T] {
	var defaultVal T
	return Option[T]{
		false: defaultVal,
	}
}

// NullIf is a function to make an Option type value that has an explicit null value if the condition is true.
// Otherwise, it returns Some value.
func NullIf[T any](v T, cond bool) Option[T] {
	if cond {
		return Null[T]()
	}
	return Some(v)
}

// FromNillable is a function to make an Option type value with the nillable value with value de-referencing.
// If the given value is not nil, this returns Some[T] value.
// On the other hand, if the value is nil, this returns None[T].
// This function does "dereference" for the value on packing that into Option value.
// If this value is not preferable, please consider using PtrFromNillable() instead.
func FromNillable[T any](v *T) Option[T] {
	if v == nil {
		return None[T]()
	}
	return Some(*v)
}

// PtrFromNillable is a function to make an Option type value with the nillable value without value de-referencing.
// If the given value is not nil, this returns Some[*T] value. On the other hand, if the value is nil, this returns None[*T].
// This function doesn't "dereference" the value on packing that into the Option value; in other words, this puts the as-is pointer value into the Option envelope.
// This behavior contrasts with the FromNillable() function's one.
func PtrFromNillable[T any](v *T) Option[*T] {
	if v == nil {
		return None[*T]()
	}
	return Some(v)
}

// IsSome returns whether the Option has a value or not and is not null.
func (o Option[T]) IsSome() bool {
	if len(o) == 0 {
		return false
	}
	_, ok := o[true]
	return ok
}

// IsNone returns whether the Option doesn't have a value or not.
func (o Option[T]) IsNone() bool {
	return len(o) == 0
}

// IsNull returns whether the Option has an explicit null value or not.
func (o Option[T]) IsNull() bool {
	if len(o) == 0 {
		return false
	}
	_, ok := o[false]
	return ok
}

// Unwrap returns the value regardless of Some/None/Null status.
// If the Option value is Some, this method returns the actual value.
// On the other hand, if the Option value is None or Null, this method returns the *default* value according to the type.
func (o Option[T]) Unwrap() T {
	if o.IsNone() || o.IsNull() {
		var defaultValue T
		return defaultValue
	}

	return o[true]
}

// UnwrapAsPtr returns the contained value in receiver Option as a pointer.
// This is similar to `Unwrap()` method but the difference is this method returns a pointer value instead of the actual value.
// If the receiver Option value is None or Null, this method returns nil.
func (o Option[T]) UnwrapAsPtr() *T {
	if o.IsNone() || o.IsNull() {
		return nil
	}

	var v = o[true]
	return &v
}

// Take takes the contained value in Option.
// If Option value is Some, this returns the value.
// If Option value is Null, this returns the default value.
// If Option value is None, this returns an ErrNoneValueTaken as the second return value.
func (o Option[T]) Take() (T, error) {
	if o.IsNone() {
		var defaultValue T
		return defaultValue, ErrNoneValueTaken
	}

	return o.Unwrap(), nil
}

// TakeOr returns the actual value if the Option has a value (Some or Null).
// Otherwise, it returns the provided fallback value.
func (o Option[T]) TakeOr(fallbackValue T) T {
	if o.IsNone() {
		return fallbackValue
	}

	return o.Unwrap()
}

// TakeOrElse returns the actual value if the Option has a value (Some or Null).
// Otherwise, it executes the fallback function and returns the result.
func (o Option[T]) TakeOrElse(fallbackFunc func() T) T {
	if o.IsNone() {
		return fallbackFunc()
	}

	return o.Unwrap()
}

// Or returns the current Option if it has a value (Some or Null).
// If the current Option is None, it returns the fallback Option.
func (o Option[T]) Or(fallbackOptionValue Option[T]) Option[T] {
	if o.IsNone() {
		return fallbackOptionValue
	}

	return o
}

// Filter returns the current Option if it has a value and the value matches the predicate.
// If the Option is Null, the default value is used for the predicate.
// Otherwise, it returns None.
func (o Option[T]) Filter(predicate func(v T) bool) Option[T] {
	if o.IsNone() || !predicate(o.Unwrap()) {
		return None[T]()
	}

	return o
}

// IfSome calls the provided function with the value of Option if it is Some.
func (o Option[T]) IfSome(f func(v T)) {
	if o.IsNone() || o.IsNull() {
		return
	}

	f(o.Unwrap())
}

// IfSomeWithError calls the provided function with the value of Option if it is Some.
// This propagates the error from the provided function.
func (o Option[T]) IfSomeWithError(f func(v T) error) error {
	if o.IsNone() || o.IsNull() {
		return nil
	}

	return f(o.Unwrap())
}

// IfNone calls the provided function if the Option is None.
func (o Option[T]) IfNone(f func()) {
	if !o.IsNone() {
		return
	}

	f()
}

// IfNoneWithError calls the provided function if the Option is None.
// This propagates the error from the provided function.
func (o Option[T]) IfNoneWithError(f func() error) error {
	if !o.IsNone() {
		return nil
	}

	return f()
}

// IfNull calls the provided function if the Option is Null.
func (o Option[T]) IfNull(f func()) {
	if !o.IsNull() {
		return
	}

	f()
}

// IfNullWithError calls the provided function if the Option is Null.
// This propagates the error from the provided function.
func (o Option[T]) IfNullWithError(f func() error) error {
	if !o.IsNull() {
		return nil
	}

	return f()
}

// String returns a string representation of the Option.
// It includes the unwrapped value for both Some and Null, and if the value implements fmt.Stringer, it uses its custom string representation.
func (o Option[T]) String() string {
	if o.IsNone() {
		return "None[]"
	}

	// Unwrap the value for both Some and Null
	v := o.Unwrap()

	// Check if the value implements fmt.Stringer for custom string formatting
	if stringer, ok := interface{}(v).(fmt.Stringer); ok {
		if o.IsNull() {
			return fmt.Sprintf("Null[%s]", stringer.String())
		}
		return fmt.Sprintf("Some[%s]", stringer.String())
	}

	// Default formatting when fmt.Stringer is not implemented
	if o.IsNull() {
		return fmt.Sprintf("Null[%v]", v)
	}
	return fmt.Sprintf("Some[%v]", v)
}

// MarshalJSON implements the json.Marshaler interface for Option.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	// if field was specified, and `null`, marshal it
	if o.IsNull() {
		return NullBytes, nil
	}

	// if field was unspecified, and `omitempty` is set on the field's tags, `json.Marshal` will omit this field

	// otherwise: we have a value, so marshal it
	return json.Marshal(o[true])
}

// UnmarshalJSON implements the json.Unmarshaler interface for Option.
func (o *Option[T]) UnmarshalJSON(data []byte) error {
	// if field is unspecified, UnmarshalJSON won't be called

	// if field is specified, and `null`
	if bytes.Equal(data, NullBytes) {
		*o = Null[T]()
		return nil
	}
	// otherwise, we have an actual value, so parse it
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}
