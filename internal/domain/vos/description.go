package vos

import (
	"database/sql/driver"
	"errors"
)

var errInvalidDescriptionSize = errors.New("invalid description")

var maxDescriptionSize = 255

type Description struct {
	value string
}

func NewDescription(value string) (Description, error) {
	if len(value) > maxDescriptionSize {
		return Description{}, errInvalidDescriptionSize
	}
	return Description{value: value}, nil
}

func (d Description) Value() (driver.Value, error) {
	return d.value, nil
}

func (d *Description) Scan(src any) error {
	switch v := src.(type) {
	case string:
		d.value = v
		return nil
	case []byte:
		d.value = string(v)
		return nil
	case nil:
		d.value = ""
		return nil
	}
	return errors.New("unsupported type for Description")
}

func (d Description) String() string {
	return d.value
}