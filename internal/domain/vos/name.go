package vos

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNameEmpty   = errors.New("name cannot be empty")
	ErrNameTooLong = errors.New("name cannot exceed 38 characters")
	ErrInvalidName = errors.New("invalid name type from database")
)

type Name struct {
	value string
}

// NewName creates a new validated Name value object.
func NewName(v string) (Name, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return Name{}, ErrNameEmpty
	}
	// Assuming max 38 runes for proper unicode support, but len(v) works for bytes.
	// We'll use []rune(v) to accurately count characters instead of bytes.
	if len([]rune(v)) > 38 {
		return Name{}, ErrNameTooLong
	}
	return Name{value: v}, nil
}

// String returns the underlying string value.
func (n Name) String() string {
	return n.value
}

// Value implements the database/sql/driver Valuer interface.
func (n Name) Value() (driver.Value, error) {
	return n.value, nil
}

// Scan implements the database/sql Scanner interface.
func (n *Name) Scan(value interface{}) error {
	if value == nil {
		return ErrNameEmpty
	}

	var parsedName string
	switch v := value.(type) {
	case string:
		parsedName = v
	case []byte:
		parsedName = string(v)
	default:
		return fmt.Errorf("%w: %T", ErrInvalidName, value)
	}

	name, err := NewName(parsedName)
	if err != nil {
		return err
	}

	*n = name
	return nil
}
