package optional

// NOTE: This file is for backward compatibility with the previous version of this package.
// It is not recommended to use the Scanner and Valuer interfaces for the Option type.

import (
	"database/sql/driver"
	"errors"
	"time"
)

var (
	ErrSQLScannerIncompatibleDataType = errors.New("incompatible data type for SQL scanner on Option[T]")
)

// Scan assigns a value from a database driver.
// This method is required from database/sql.Scanner interface.
func (o *Option[T]) Scan(src any) error {
	if src == nil {
		*o = None[T]()
		return nil
	}

	switch src.(type) {
	case string, []byte, int64, float64, bool, time.Time:
		*o = Some[T](src.(T))
	default:
		return ErrSQLScannerIncompatibleDataType
	}

	return nil
}

// Value returns a driver Value.
// This method is required from database/sql/driver.Valuer interface.
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() || o.IsJsonNull() {
		return nil, nil
	}

	return driver.DefaultParameterConverter.ConvertValue(o.Unwrap())
}
